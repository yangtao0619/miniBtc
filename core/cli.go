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
	fmt.Println("args:", args)
	if len(args) < 2 {
		printUsage("参数格式错误")

	}
	cmd := args[1]
	switch cmd {
	//创建区块链的操作和添加区块的操作需要分开
	case "generateBc":
		checkArgs(4)
		//获取创建该区块链的地址
		var address string
		if "--address" == args[2] && args[3] != "" {
			address = args[3]
		} else {
			printUsage("地址格式错误")
		}
		CreateBlockChain(address)
	case "printBc":
		checkArgs(2)
		cli.PrintBlockChain()
	case "getBalance":
		//查看某个账户的余额
		checkArgs(4)
		//获取创建该区块链的地址
		var address string
		if "--address" == args[2] && args[3] != "" {
			address = args[3]
		} else {
			printUsage("地址格式错误")
		}
		cli.GetBalance(address)
	case "send":
		//查看某个账户的余额
		//block send from ADDRESS1 to ADDRESS2 miner ADDRESS3 amount Money data DATA
		checkArgs(12)
		//获取创建该区块链的地址
		from := args[3]
		to := args[5]
		miner := args[7]
		amount, _ := strconv.ParseFloat(args[9], 64)
		data := []byte(args[11])
		cli.SendToSomeOne(from, to, miner, amount, data)
	case "createWallet":
		checkArgs(2)
		cli.CreateWallet()
	case "listAllAddress":
		checkArgs(2)
		cli.ListAllAddress()
	default:
		fmt.Println("参数错误")
	}
}
func checkArgs(need int) {
	if len(os.Args) != need {
		fmt.Println("参数无效")
		os.Exit(-1)
	}
}

func printUsage(errInfo string) {
	fmt.Printf("%s,%s", errInfo, Usage)
	os.Exit(-1)
}
