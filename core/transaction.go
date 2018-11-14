package core

import (
	"encoding/gob"
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/base58"
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"
	"crypto/elliptic"
)

//需要定义输入和输出

type TXInput struct {
	//引用的交易id
	TXId []byte

	//索引值
	TxIndex int64

	//签名
	Sig []byte

	//解锁脚本
	PubKey []byte
}

type TXOutput struct {
	//金额
	Value float64

	//锁定脚本
	PublicKeyHash []byte
}

type Transaction struct {
	Id        []byte
	TxInputs  []TXInput
	TxOutputs []TXOutput
}

func (tx *Transaction) String() {
	//todo 重写String方法
}

func (transaction *Transaction) SetTxId() {
	//设置交易当前的hash,需要将tx编码
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encodeErr := encoder.Encode(transaction)
	if encodeErr != nil {
		panic(encodeErr)
	}
	sumHash := sha256.Sum256(buffer.Bytes())
	transaction.Id = sumHash[:]
}

const reward = 12.5

func NewOutput(value float64, address string) *TXOutput {
	txOutput := TXOutput{Value: value,}
	//使用地址进行锁定
	decode := base58.Decode(address)
	hash := decode[1 : len(decode)-4]
	txOutput.PublicKeyHash = hash
	return &txOutput
}

//这个方法创建新区块的第一笔交易,没有输入,只有输出
func CreateCoinBase(data, miner string) *Transaction {
	fmt.Println("创建挖矿交易", miner)
	//这里的输入应该为空
	input := TXInput{nil, -1, nil, []byte(data)}
	//输出到矿工的地址
	output := NewOutput(reward, miner)
	//组合成交易
	tx := Transaction{nil, []TXInput{input}, []TXOutput{*output}}
	tx.SetTxId()
	//返回交易
	return &tx
}

//判断是否是第一笔交易
func (tx *Transaction) isCoinBase() bool {
	if len(tx.TxInputs) == 1 && tx.TxInputs[0].TXId == nil && tx.TxInputs[0].TxIndex == -1 {
		fmt.Println("这是一笔挖矿交易")
		return true
	}
	return false
}

//创建一笔交易的函数
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	fmt.Println("正在创建交易,from ", from, " to ", to, " amount ", amount)
	//from转化成publichash
	walletCollection := NewWalletCollection()
	wallet := walletCollection.WalletClt[from]
	if wallet == nil {
		fmt.Printf("本地没有%s钱包,无法创建新的交易\n", from)
		return nil
	}
	pubKey := wallet.PubKey
	privateKey := wallet.PrivateKey

	fmt.Println("pubkey is ", pubKey, " privateKey is ", privateKey, " real pubkey ",
		append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...))

	pubkeyHash := HashPubkey(pubKey)
	//找出交易发起人所有可支配的交易输出
	suitableUtxos, clac := bc.GetSuitableUtxos(pubkeyHash, amount)
	//得到适合的交易输出中的信息,组合成交易的输入和输出信息
	var inputs []TXInput
	var outputs []TXOutput
	for txId, txIndexArr := range suitableUtxos {
		//组合成交易输入
		for _, index := range txIndexArr {
			txInput := TXInput{[]byte(txId), index, nil, pubKey}
			inputs = append(inputs, txInput)
		}
	}

	//判断余额是否充足
	if clac < amount {
		fmt.Println("余额不足,交易创建失败")
	}

	//如果充足的话,需要创建一个交易输出
	output := NewOutput(amount, to)
	//如果还有剩余的话,要再创建一个输出给自己
	outputs = append(outputs, *output)
	if clac-amount > 0 {
		fmt.Println("创建一笔交易给自己,from ", from)
		outputToSelf := NewOutput(clac-amount, from)
		outputs = append(outputs, *outputToSelf)
	}
	transaction := Transaction{nil, inputs, outputs}
	//给交易设置交易id
	transaction.SetTxId()
	if !bc.Sign(&transaction, privateKey) {
		fmt.Println("签名失败,交易创建失败")
		return nil
	}
	fmt.Println("交易创建成功")
	return &transaction
}

//实际进行交易签名的方法
func (transaction *Transaction) Sign(private *ecdsa.PrivateKey, transactions map[string]*Transaction) bool {

	//如果当前交易是coinbase交易的话不需要进行签名,因为它没有输入,怎么遍历呢?
	if transaction.isCoinBase() {
		return true
	}

	//首先要对当前的交易数据进行拷贝,创建一个副本
	copyTx := transaction.CopyTransaction()

	//遍历副本中的inputs集合,找到其引用的交易,对每一笔交易都要进行签名
	for i, input := range copyTx.TxInputs {
		txId := input.TXId
		tx := transactions[string(txId)]

		//对tx交易进行签名,然后更改原数据

		//给pubkey赋值
		copyTx.TxInputs[i].PubKey = tx.TxOutputs[input.TxIndex].PublicKeyHash

		//给sign赋值,首先要生成签名的数据
		copyTx.SetTxId()

		//设置完数据后将copy中的pubkey置为空,防止验证的时候数据不一致
		copyTx.TxInputs[i].PubKey = nil
		hash := copyTx.Id
		fmt.Println("i is ", i, " hash is ", hash)
		r, s, err := ecdsa.Sign(rand.Reader, private, hash[:])
		if err != nil {
			fmt.Println("签名失败,err:", err)
			return false
		}
		var sign []byte
		//签名之后将签名的值赋给原始交易的sig字段
		sign = append(r.Bytes(), s.Bytes()...)
		transaction.TxInputs[i].Sig = sign
		fmt.Println("sign is ",sign)
	}
	//所有的交易数据签名完成之后,返回true
	fmt.Println("签名成功")
	return true
}
func (transaction *Transaction) CopyTransaction() *Transaction {

	var inputs []TXInput

	for _, input := range transaction.TxInputs {
		inputs = append(inputs, TXInput{input.TXId, input.TxIndex, nil, nil})
	}
	return &Transaction{transaction.Id, inputs, transaction.TxOutputs}
}

//验证所有的交易是否合法
func (transaction *Transaction) Verify(transactions map[string]*Transaction) bool {
	//校验传来的所有的交易是否是合法有效的
	copyTx := transaction.CopyTransaction()
	for i, input := range transaction.TxInputs {
		//将output中的公钥hash赋值给input的公钥
		prevTx := transactions[string(input.TXId)]
		//得到要签名的数据
		copyTx.TxInputs[i].PubKey = prevTx.TxOutputs[input.TxIndex].PublicKeyHash

		//得到要签名的hash
		copyTx.SetTxId()
		hash := copyTx.Id
		copyTx.TxInputs[i].PubKey = nil

		pubKey := input.PubKey
		sign := input.Sig

		//将sign一分为2
		r := big.Int{}
		s := big.Int{}

		rData := sign[:len(sign)/2]
		sData := sign[len(sign)/2:]

		r.SetBytes(rData)
		s.SetBytes(sData)

		x1 := big.Int{}
		y1 := big.Int{}

		x1Data := pubKey[:len(pubKey)/2]
		y1Data := pubKey[len(pubKey)/2:]

		x1.SetBytes(x1Data)
		y1.SetBytes(y1Data)

		curve := elliptic.P256()
		pubKeyOrigin := ecdsa.PublicKey{Curve: curve, X: &x1, Y: &y1}

		//这里需要的数据分别是公钥,要校验的hash值,r和s的值

		if !ecdsa.Verify(&pubKeyOrigin, hash, &r, &s) {
			fmt.Println("校验失败")
			fmt.Println("i is ", i, " hash ", hash)
			fmt.Println("sign is ",sign)
			return false
		}
	}
	return true
}
