package core

import (
	"time"
	"bytes"
	"encoding/gob"
	"fmt"
	"crypto/sha256"
)

//定义每个区块的结构体
type Block struct {
	//版本号
	Version uint64
	//前一个区块的hash
	PrevHash []byte
	//默克尔根
	MerkelRoot []byte
	//时间戳
	TimeStamp uint64
	//当前区块的hash
	Hash []byte
	//挖矿的难度
	Difficulty uint64
	//随机数
	Nonce uint64
	//区块需要存储的hash,存储交易数据
	Transactions []*Transaction
}

func CreateGenesisBlock(data []*Transaction, prevHash []byte) *Block {
	return NewBlock(data, prevHash)
}

//创建新的区块
func NewBlock(data []*Transaction, prevHash []byte) *Block {
	newBlock := &Block{
		Version:      00,
		PrevHash:     prevHash,
		TimeStamp:    uint64(time.Now().UnixNano()),
		Difficulty:   difficulty,
		Transactions: data,
	}
	merkelRoot := newBlock.GetMerkelRoot()
	newBlock.MerkelRoot = merkelRoot
	proofOfWork := NewProofOfWork(newBlock)
	//newBlock.SetHash()
	hash, nonce := proofOfWork.Run()
	newBlock.Hash = hash
	newBlock.Nonce = nonce
	fmt.Printf("generate new block,hash is %x\n nonce is %d", hash, nonce)
	return newBlock
}

//设置组合的梅克尔根
func (block *Block) GetMerkelRoot() []byte{
	//模拟组合的哈希
	var info []byte
	for _, tx := range block.Transactions {
		info = append(info, tx.Id...)
	}
	sumHash := sha256.Sum256(info)
	return sumHash[:]
}

/*//设置当前区块的hash
func (block *Block) SetHash() {
	//Join入参,前者为一个二维数组,后者为进行拼接的一维数组
	bytesInfo := bytes.Join([][]byte{uint64ToByte(block.Version), uint64ToByte(block.Nonce), uint64ToByte(block.Difficulty),
		uint64ToByte(block.Nonce), uint64ToByte(block.TimeStamp), block.MerkelRoot, block.Transactions, block.Hash, block.PrevHash}, []byte{})
	//对得到的bytes数组进行hash运算
	hash := sha256.Sum256(bytesInfo)
	block.Hash = hash[:]
}*/

//将一个区块对象转换成字节流
func (block *Block) toBytes() []byte {
	//使用gob的encoder
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encodeErr := encoder.Encode(block)
	if encodeErr != nil {
		panic(encodeErr)
	}
	return buffer.Bytes()
}

//将字节流转换成一个block对象
func ToBlock(data []byte) *Block {
	var block *Block
	var buffer bytes.Buffer
	_, err := buffer.Write(data)
	if err != nil {
		panic(err)
	}
	decoder := gob.NewDecoder(&buffer)
	decodeErr := decoder.Decode(&block)
	if decodeErr != nil {
		panic(decodeErr)
	}
	return block
}
