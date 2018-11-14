package core

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"fmt"
)

//定义工作量证明的结构
type ProofOfWork struct {
	//区块作为数据的来源
	block Block

	//需要有一个固定的难度值
	target big.Int
}

const difficulty = 16

//创建一个新的工作量证明对象
func NewProofOfWork(block *Block) *ProofOfWork {
	//比特币中的难度值是计算的
	pow := ProofOfWork{
		block: *block,
	}

	//开始计算,直接写成hash
	//targetStr := "0000100000000000000000000000000000000000000000000000000000000000"
	targetInt := big.NewInt(1)
	targetInt.Lsh(targetInt, 256-difficulty)
	pow.target = *targetInt
	return &pow
}

//运算hash,满足小于难度值
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	//步骤一.拿到block的数据
	fmt.Println("pow running")
	var nonce uint64 = 0
	//步骤二.对block数据进行hash运算
	//Join入参,前者为一个二维数组,后者为进行拼接的一维数组
	for {
		//fmt.Println("pow circle......")
		bytesInfo := pow.prepareData(nonce)
		//对得到的bytes数组进行hash运算
		sum256Hash := sha256.Sum256(bytesInfo)
		//步骤三.比较hash值
		//先将hash转成bitInt便于和target做比较
		var tempInt big.Int
		tempInt.SetBytes(sum256Hash[:])
		//小于目标值
		/*
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		//
		func (x *Int) Cmp(y *Int) (r int) {}
		 */
		if tempInt.Cmp(&pow.target) == -1 {
			//挖矿成功
			return sum256Hash[:], nonce
		} else { //大于目标值
			//log.Println("大于目标值,hash is", sum256Hash[:], "nonce is", nonce)
			nonce++
		}

	}
	return []byte{}, uint64(0)
}

func (pow *ProofOfWork) prepareData(nonce uint64) []byte {
	block := pow.block
	merkelRoot := pow.block.GetMerkelRoot()
	bytesInfo := bytes.Join([][]byte{uint64ToByte(block.Version),
		uint64ToByte(nonce),
		uint64ToByte(block.Difficulty),
		uint64ToByte(block.TimeStamp),
		block.PrevHash},
		merkelRoot)
	return bytesInfo
}

//验证hash
func (pow *ProofOfWork) IsValid() bool {
	//先计算pow中的block的hash值
	block := pow.block
	prepareData := pow.prepareData(block.Nonce)
	sum256Hash := sha256.Sum256(prepareData)
	//转成bitInt
	var tempInt big.Int
	tempInt.SetBytes(sum256Hash[:])
	//和当前的target做比较
	return tempInt.Cmp(&pow.target) == -1
}

//将uint64转换成byte数组
func uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	//将num以二进制大尾端的形式写入缓冲区
	writeErr := binary.Write(&buffer, binary.BigEndian, num)
	if writeErr != nil {
		log.Fatal("uint64 to num err:", writeErr)
		return nil
	}
	return buffer.Bytes()
}
