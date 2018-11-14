package core

import (
	"github.com/boltdb/bolt"
	"os"
	"log"
	"fmt"
	"time"
	"errors"
	"bytes"
	"github.com/base58"
	"crypto/ecdsa"
)

type BlockChain struct {
	//Blocks []*Block
	//区块链中应该存储一个数据库操作的句柄
	db *bolt.DB
	//还有区块链最尾部区块的hash,便于添加区块
	tailHash []byte
}

const (
	dbName      = "BlockChain.db"
	bucketName  = "BlockChainBucket"
	lastHashKey = "lastHashKey"
)

func CreateBlockChain(address string) *BlockChain {
	//创建区块链之前需要检测一下区块链数据库文件是否存在
	if isBlockChainExists() {
		fmt.Println("BlockChain already exist,please addBlock!")
		//os.Exit(1)
		return nil
	}

	//返回一个区块链对象
	//首先需要检查一下数据库中是否有区块链的数据,没有的话需要新建,有的话直接给区块链对象赋值即可
	blockChain := new(BlockChain)
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		panic(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil { //如果为空的话就创建这个bucket
			bucket, err = tx.CreateBucket([]byte(bucketName))
			if err != nil {
				return err
			}
			//创建创世区块
			coinBase := CreateCoinBase("hello coin base", address)
			genesisBlock := CreateGenesisBlock([]*Transaction{coinBase}, []byte{})
			//将创世区块写入区块链数据库中,同时更新lastHsh
			bucket.Put(genesisBlock.Hash, genesisBlock.toBytes())
			bucket.Put([]byte(lastHashKey), genesisBlock.Hash)
			blockChain.tailHash = genesisBlock.Hash
		}
		return nil
	})
	//将数据库的句柄赋值给db
	blockChain.db = db
	return blockChain
}
func isBlockChainExists() bool {
	_, err := os.Stat(dbName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func GetBlockChainObject() *BlockChain {
	//创建区块链之前需要检测一下区块链数据库文件是否存在
	if !isBlockChainExists() {
		fmt.Println("BlockChain not exist,please create first!")
		//os.Exit(1)
		return nil

	}
	//返回一个区块链对象
	//首先需要检查一下数据库中是否有区块链的数据,没有的话需要新建,有的话直接给区块链对象赋值即可
	blockChain := new(BlockChain)
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		panic(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		blockChain.tailHash = bucket.Get([]byte(lastHashKey))
		return nil
	})
	//将数据库的句柄赋值给db
	blockChain.db = db
	return blockChain
}

func (blockChain *BlockChain) AddBlock(transactions []*Transaction) {
	//添加区块的时候要对区块中的交易数据进行校验
	validTxs := make([]*Transaction, 0)

	for _, tx := range transactions {
		if blockChain.Verify(tx) {
			validTxs = append(validTxs, tx)
		} else {
			fmt.Println("交易校验失败")
		}
	}

	//向数据库中添加新的区块
	blockChain.db.Update(func(tx *bolt.Tx) error {
		//打开水桶
		bucket := tx.Bucket([]byte(bucketName))
		//判断水桶是否为空
		if bucket == nil {
			//状态错误,直接退出
			log.Fatal("bucket can not be nil")
			//os.Exit(0)
			return errors.New("bucket can not be nil")
		} else {
			//不为空,向里面插入数据
			//创建一个新的block对象
			newBlock := NewBlock(validTxs, blockChain.tailHash)
			bucket.Put(newBlock.Hash, newBlock.toBytes())
			//切记要将最后一个hash put进去
			bucket.Put([]byte(lastHashKey), newBlock.Hash)
			blockChain.tailHash = newBlock.Hash
			return nil
		}
		//更新内存区块链对象的数据
		return nil
	})
}

//区块链需要有一个迭代器用于返回一个当前指向的区块
type BlockChainIterator struct {
	db               *bolt.DB
	currentBlockHash []byte
}

//创建迭代器的方法
func (blockChain *BlockChain) CreateIterator() *BlockChainIterator {
	return &BlockChainIterator{blockChain.db, blockChain.tailHash}
}

func (bcIterator *BlockChainIterator) GetBlock() *Block {
	//首先需要打开数据库
	var currentBlock *Block
	handler := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			//如果bucket为nil的时候,panic
			//os.Exit(-1)
			return errors.New("bucket is nil")
		} else {
			//读取数据
			blockBytes := bucket.Get(bcIterator.currentBlockHash)
			currentBlock = ToBlock(blockBytes)
			//左移一位
			bcIterator.currentBlockHash = currentBlock.PrevHash
		}
		return nil
	}
	bcIterator.db.View(handler)
	return currentBlock
}

func (blockChain *BlockChain) PrintBlockChain() {
	//创建迭代器
	iterator := blockChain.CreateIterator()
	//遍历
	for {
		block := iterator.GetBlock()
		fmt.Printf("===============================\n")
		fmt.Printf("Version :%d\n", block.Version)
		fmt.Printf("PrevBlockHash :%x\n", block.PrevHash)
		fmt.Printf("MerkeRoot :%x\n", block.MerkelRoot)
		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("TimeStamp: %s\n", timeFormat)
		//fmt.Printf("TimeStamp :%d\n", block.TimeStamp)
		fmt.Printf("Difficulty :%d\n", block.Difficulty)
		fmt.Printf("Nonce :%d\n", block.Nonce)
		fmt.Printf("Hash :%x\n", block.Hash)
		//fmt.Printf("Transactions :%s\n", block.Transactions[0].TxInputs[0].PubKey)
		pow := NewProofOfWork(block)
		fmt.Printf("IsValid : %v\n\n", pow.IsValid())
		if len(block.PrevHash) == 0 {
			break
		}
	}
}

//这是一个获得所有可使用的UTXO集合的方法
func (blockChain *BlockChain) GetUtxos(pubkeyHash []byte) []UtxoInfo {
	//遍历所有的区块,得到所有的交易数据
	//获得一个遍历器
	fmt.Println("正在查询所有该地址的utxo")
	iterator := blockChain.CreateIterator()
	var utxoInfos []UtxoInfo
	spentUtxo := make(map[string][]int64)
	for {
		block := iterator.GetBlock()
		fmt.Println("遍历区块", block.TimeStamp)
		//遍历该区块所有的交易数据
		for _, tx := range block.Transactions {
		ScanTransaction:
		//找出自己能解锁的输出
			for opIndex, txOutPut := range tx.TxOutputs {
				//判断当前的index和已经花费的输入index是否一致,一致的话,继续下一次循环

				//为什么后遍历input集合,简单想就是在进行下一笔交易的输出之前先检查之前所有的输入是否被消费
				if bytes.Equal(txOutPut.PublicKeyHash, pubkeyHash) {
					if len(spentUtxo[string(tx.Id)]) != 0 {
						fmt.Println("当前交易有消耗过的output,address = ", string(pubkeyHash))
						inputIndexArr := spentUtxo[string(tx.Id)]
						for _, index := range inputIndexArr {
							if opIndex == int(index) {
								fmt.Println("找到被花费的输出,txid is ", string(tx.Id), " index is ", index)
								continue ScanTransaction
							}
						}
					}

					//当找到一笔自己能解锁的交易的时候,就把这条输出放在集合里面
					utxoInfo := UtxoInfo{tx.Id, int64(opIndex), txOutPut}
					utxoInfos = append(utxoInfos, utxoInfo)
				}
			}
			//在遍历输入之前需要先判断当前的交易是否是CoinBase交易,是的话就没有必要遍历了
			if !tx.isCoinBase() {
				//需要找出这个地址已经消耗掉的交易输出,所以要遍历这个地址的输入,查看是否引用到之前的交易输出
				for _, txInput := range tx.TxInputs {
					if bytes.Equal(HashPubkey(txInput.PubKey), pubkeyHash) {
						//如果存在能解开的输入,将这个输入的索引记录下来
						fmt.Println("这笔交易已经被花费,记录,txid is ", string(txInput.TXId), "index is ", txInput.TxIndex)
						spentUtxo[string(txInput.TXId)] = append(spentUtxo[string(txInput.TXId)], txInput.TxIndex)
					}
				}
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	fmt.Println("所有的utxo查询完毕")
	return utxoInfos
}

//这是一个获得指定账户所有余额的方法
func (blockChain *BlockChain) GetBalance(address string) float64 {
	//对传进来的地址进行解码
	decode := base58.Decode(address)
	pubkeyHash := decode[1 : len(decode)-4]

	//首先需要得到所有可以使用的交易输出的集合
	utxoinfos := blockChain.GetUtxos(pubkeyHash)
	var count float64
	for _, utxoInfo := range utxoinfos {
		count += utxoInfo.Utxo.Value
	}
	return count
}

type UtxoInfo struct {
	TxId  []byte
	Index int64
	Utxo  TXOutput
}

func (blockChain *BlockChain) GetSuitableUtxos(pubkeyHash []byte, amount float64) (map[string][]int64, float64) {
	//查看自己能够使用的utxo,满足amount的时候就直接退出
	//遍历所有的区块
	utxoInfos := blockChain.GetUtxos(pubkeyHash)
	//如果余额满足就退出遍历
	calc := 0.0
	needUtxos := make(map[string][]int64)
	for _, utxoInfo := range utxoInfos {
		key := string(utxoInfo.TxId)
		needUtxos[key] = append(needUtxos[key], utxoInfo.Index)
		calc += utxoInfo.Utxo.Value
		if calc >= amount {
			fmt.Println("已经满足,calc is", calc)
			return needUtxos, calc
		}
	}
	fmt.Println("已经满足,calc is", calc)
	return needUtxos, calc
}

//这个方法提供当前交易inputs应用的所有交易对象
func (blockChain *BlockChain) Sign(transaction *Transaction, privateKey *ecdsa.PrivateKey) bool {
	fmt.Println("对交易的数据进行签名")
	prevTxs := make(map[string]*Transaction)
	for _, input := range transaction.TxInputs {
		txId := string(input.TXId)
		tx := blockChain.FindTxById(txId)
		if tx == nil {
			return false
		}
		prevTxs[txId] = tx
	}

	return transaction.Sign(privateKey, prevTxs)
}

//根据交易id,找到对应的交易
func (blockChain *BlockChain) FindTxById(txId string) *Transaction {
	//遍历所有的区块
	iterator := blockChain.CreateIterator()
	for {
		block := iterator.GetBlock()

		transactions := block.Transactions
		for _, tx := range transactions {
			if bytes.Compare(tx.Id, []byte(txId)) == 0 {
				return tx
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return nil
}

//对区块中的交易进行验证
func (blockChain *BlockChain) Verify(transaction *Transaction) bool {
	//根据交易id获得之前所有相关的交易
	prevTxs := make(map[string]*Transaction)
	for _, input := range transaction.TxInputs {
		//根据id找交易
		txId := string(input.TXId)
		prevTx := blockChain.FindTxById(txId)
		if prevTx == nil {
			return false
		}
		prevTxs[txId] = prevTx
	}
	return transaction.Verify(prevTxs)
}
func (blockChain *BlockChain) GetChainData() string {
	//创建迭代器
	iterator := blockChain.CreateIterator()
	//遍历
	data := ""
	for {
		block := iterator.GetBlock()
		data += fmt.Sprintf("Version :%d\n", block.Version)
		data += fmt.Sprintf("PrevBlockHash :%x\n", block.PrevHash)
		data += fmt.Sprintf("MerkeRoot :%x\n", block.MerkelRoot)
		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		data += fmt.Sprintf("TimeStamp: %s\n", timeFormat)
		//fmt.Printf("TimeStamp :%d\n", block.TimeStamp)
		data += fmt.Sprintf("Difficulty :%d\n", block.Difficulty)
		data += fmt.Sprintf("Nonce :%d\n", block.Nonce)
		data += fmt.Sprintf("Hash :%x\n", block.Hash)
		/*for _,tx := range block.Transactions{
			data += fmt.Sprintf("Transactions :%s\n", tx.String)
		}*/
		pow := NewProofOfWork(block)
		data += fmt.Sprintf("IsValid : %v\n\n", pow.IsValid())
		if len(block.PrevHash) == 0 {
			break
		}
		data += "---------------------------------------next block----------------------------------------------------\n"
	}
	return data
}
