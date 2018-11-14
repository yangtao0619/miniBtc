package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/ripemd160"
	"crypto/sha256"
	"github.com/base58"
	"bytes"
)

/*
比特币的钱包,需要保存私钥和公钥哈希
钱包执行的功能:
创建一个新的钱包
生成公钥的hash
生成地址
*/

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PubKey     []byte //公钥用一对x,y字节切片拼装组成
}

func NewWallet() *Wallet {
	//创建钱包

	//首先生成私钥
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("generate key err:", err)
		panic(err)
	}

	//使用私钥生成公钥
	publicKey := privateKey.PublicKey

	var pubkey []byte
	pubkey = append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)

	return &Wallet{PrivateKey: privateKey, PubKey: pubkey}
}



//生成钱包的地址,地址是由公钥进行一系列运算之后得到的
func (wallet *Wallet) GetAddress() string {
	//首先对公钥进行RIPEMD160计算对pubkey的sha256进行运算
	pubkeyHash := HashPubkey(wallet.PubKey)
	//将公钥的hash和version拼接得到21字节的byte切片b
	version := byte(00)
	b := append([]byte{version}, pubkeyHash...)
	//将b拷贝一份进行两次sha256运算并取出前四个字节组成切片c
	c := checkSum(b)
	//将b和c拼接在一起进行base58运算得到地址
	d := append(b, c...)
	address := base58.Encode(d)
	return address
}

//校验地址是否正确
func IsAddressValid(address string) bool{
	//1. 解码base58
	decodeInfo := base58.Decode(address)

	//2. 截取前21字节和 后四个字节
	payload := decodeInfo[0: len(decodeInfo)-4]
	checksum1 := decodeInfo[len(decodeInfo)-4:]

	//3. 对前21字节进行checksum计算
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])
	checksum2 := secondHash[0:4]

	//4. 比较生成cheksum1和截取cheksum2
	return bytes.Equal(checksum1, checksum2)
}

//将b进行两次sha256运算
func checkSum(bytes []byte) []byte {
	firstHash := sha256.Sum256(bytes)
	secondHash := sha256.Sum256(firstHash[:])
	checkSum := secondHash[0:4]
	return checkSum
}

//对公钥进行hash
func HashPubkey(pubkey []byte) []byte {
	sum256Hash := sha256.Sum256(pubkey)

	ripHasher := ripemd160.New()
	_, err := ripHasher.Write(sum256Hash[:])
	if err != nil {
		panic(err)
	}
	sum := ripHasher.Sum(nil)
	return sum
}
