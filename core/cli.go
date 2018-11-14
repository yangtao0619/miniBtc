package core

import (
	"os"
	"fmt"
	"strconv"
)

//命令行工具

type CLI struct {
}

//接收命令行的参数
const Usage = `
block generateBc --address Address 生成一条新的区块链,Address表示创建人的地址
block printBc 打印出当前的区块链数据
block getBalance --address Address 得到Address地址的账户余额
block send from ADDRESS1 to ADDRESS2 miner ADDRESS3 amount MoneyValue data DATA
block createWallet 创建一个钱包
block listAllAddress 列出所有的地钱包地址
`

func (cli *CLI) Run() {
	//先检查一下命令行参数的个数是否合法
	args := os.Args
	if len(args) < 2 {
		printUsage("参数格式错误")
		return
	}
	cmd := args[1]
	switch cmd {
	//创建区块链的操作和添加区块的操作需要分开
	case "generateBc":
		if !checkArgs(4) {
			return
		}
		//获取创建该区块链的地址
		var address string
		if "--address" == args[2] && args[3] != "" {
			address = args[3]
		} else {
			printUsage("地址格式错误")
			return
		}
		if !IsAddressValid(address) {
			return
		}
		CreateBlockChain(address)
	case "printBc":
		if !checkArgs(2) {
			return
		}
		cli.PrintBlockChain()
	case "getBalance":
		//查看某个账户的余额
		if !checkArgs(4) {
			return
		}
		//获取创建该区块链的地址
		var address string
		if "--address" == args[2] && args[3] != "" {
			address = args[3]
		} else {
			printUsage("地址格式错误")
			return
		}
		if !IsAddressValid(address) {
			return
		}
		cli.GetBalance(address)
	case "send":
		//查看某个账户的余额
		//block send from ADDRESS1 to ADDRESS2 miner ADDRESS3 amount Money data DATA
		if !checkArgs(12) {
			return
		}
		//获取创建该区块链的地址
		from := args[3]
		to := args[5]
		miner := args[7]
		amount, _ := strconv.ParseFloat(args[9], 64)
		data := []byte(args[11])
		if !IsAddressValid(from) || !IsAddressValid(to) || !IsAddressValid(miner) {
			return
		}
		cli.SendToSomeOne(from, to, miner, amount, data)
	case "createWallet":
		if !checkArgs(2) {
			return
		}
		cli.CreateWallet()
	case "listAllAddress":
		if !checkArgs(2) {
			return
		}
		cli.ListAllAddress()
	default:
		fmt.Println("参数错误")
	}
}
func checkArgs(need int) bool {
	if len(os.Args) != need {
		fmt.Println("参数无效")
		return false
	}
	return true
}

func printUsage(errInfo string) {
	fmt.Printf("%s,%s", errInfo, Usage)
}
