package core

import (
	"encoding/gob"
	"bytes"
	"crypto/elliptic"
	"io/ioutil"
	"fmt"
	"os"
)

/*
保存所有钱包信息的集合,根据地址来保存
提供的功能有:
1.创建钱包
2.获取所有的钱包地址
3.加载所有的钱包信息到内存中
4.将钱包信息从内存写入文件
 */

const walletCollectionName = "wallets.dat"

type WalletCollection struct {
	WalletClt map[string]*Wallet
}

func NewWalletCollection() *WalletCollection {
	//从本地加载所有的钱包集合
	var walletCtl WalletCollection
	//注意这里要给map初始化
	walletCtl.WalletClt = make(map[string]*Wallet)
	walletCtl.loadFromFile()

	return &walletCtl
}


//创建新的钱包,返回新钱包的地址
func (walletClt *WalletCollection) CreateWallet() string {
	wallet := NewWallet()
	address := wallet.GetAddress()

	//存入map集合中
	walletClt.WalletClt[address] = wallet

	//存入到本地
	if !walletClt.saveToFile() {
		fmt.Println("保存钱包信息失败")
		return ""
	}

	return address
}

//获得所有的钱包地址
func (walletClt *WalletCollection) GetAllAddress() []string {
	var allAddress []string
	//遍历所有的钱包信息,取出地址信息
	for address := range walletClt.WalletClt {
		allAddress = append(allAddress, address)
	}
	return allAddress
}

//保存所有钱包信息到本地文件中
func (walletClt *WalletCollection) saveToFile() bool {
	fmt.Println("save ws to file")
	//使用gob编码
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)

	gob.Register(elliptic.P256())

	encodeErr := encoder.Encode(walletClt)
	if encodeErr != nil {
		fmt.Println("Encode err:", encodeErr)
		return false
	}
	//使用ioutil写入文件
	writeErr := ioutil.WriteFile(walletCollectionName, buffer.Bytes(), 0600)
	if writeErr != nil {
		fmt.Println("write err:", writeErr)
		return false
	}
	return true
}

func (walletClt *WalletCollection) loadFromFile() bool {
	//判断钱包集合文件是否存在
	if !isFileExist(walletCollectionName) {
		fmt.Println("钱包集合文件不存在")
		return false
	}
	//从文件中读取gob编码的数据
	fileBuffer, err := ioutil.ReadFile(walletCollectionName)
	if err != nil {
		fmt.Println("ioutil read file err:", err)
		return false
	}
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileBuffer))

	var tempWs WalletCollection
	decodeErr := decoder.Decode(&tempWs)
	if decodeErr != nil {
		fmt.Println("decodeErr:", decodeErr)
		return false
	}
	//转成WalletCollection
	walletClt.WalletClt = tempWs.WalletClt
	//赋值给调用该方法对象
	return true
}
func isFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
