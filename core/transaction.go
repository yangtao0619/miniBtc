package core

import (
	"encoding/gob"
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/base58"
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
	input := TXInput{nil, -1, nil,[]byte(data)}
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
		fmt.Println("这是一笔挖矿交易", string(tx.Id))
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
	if wallet != nil {
		fmt.Printf("本地没有%s钱包,无法创建新的交易", from)
		return nil
	}
	pubKey := wallet.PubKey

	pubkeyHash := HashPubkey(pubKey)
	//找出交易发起人所有可支配的交易输出
	suitableUtxos, clac := bc.GetSuitableUtxos(pubkeyHash, amount)
	//得到适合的交易输出中的信息,组合成交易的输入和输出信息
	var inputs []TXInput
	var outputs []TXOutput
	for txId, txIndexArr := range suitableUtxos {
		//组合成交易输入
		for _, index := range txIndexArr {
			txInput := TXInput{[]byte(txId), index,nil, pubKey}
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
	fmt.Println("交易创建完成,inputs are ", inputs, "outputs are ", outputs)
	return &transaction
}
