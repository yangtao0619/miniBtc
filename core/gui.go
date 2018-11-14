package core

import (
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/walk"
	"fmt"
	"strconv"
)

func StartGui() {
	//var inTE, outTE *walk.TextEdit
	var teAllAddress *walk.TextEdit
	var teSearchAddress *walk.TextEdit
	var labelSearchResult *walk.Label
	var teTxFrom *walk.TextEdit
	var teTxTo *walk.TextEdit
	var teTxMiner *walk.TextEdit
	var teTxAmount *walk.TextEdit
	var teNewAddress *walk.TextEdit
	var teTxData *walk.TextEdit
	var teNewChainState *walk.TextEdit
	var labelTxState *walk.Label
	var labelCreateChainState *walk.Label

	MainWindow{
		Title:   "SCREAMO",
		MinSize: Size{Width: 800, Height: 1000},
		Layout:  VBox{},
		Children: []Widget{
			PushButton{
				Text: "创建新的地址",
				OnClicked: func() {
					walletCollection := NewWalletCollection()
					newWalletAddress := walletCollection.CreateWallet()
					if newWalletAddress == "" {
						fmt.Println("地址创建失败")
						return
					}
					fmt.Println("您新的钱包地址是:", newWalletAddress)
					teNewAddress.SetText(newWalletAddress)
				},
			},
			TextEdit{
				Text:     "新的钱包地址",
				AssignTo: &teNewAddress,
				MaxSize:  Size{Width: 800, Height: 80},
			},
			TextEdit{
				Text:     "输入创建区块链的地址",
				AssignTo: &teNewChainState,
				MaxSize:  Size{Width: 800, Height: 80},
			},
			PushButton{
				Text: "创建区块链",
				OnClicked: func() {
					address := teNewChainState.Text()
					if address == "" {
						teNewChainState.SetText("地址不正确")
						return
					}
					fmt.Printf("使用地址%s创建区块链\n", address)
					blockChain := CreateBlockChain(address)
					if blockChain != nil {
						labelCreateChainState.SetText("建链成功")
					} else {
						labelCreateChainState.SetText("建链失败")
					}
				},
			},
			Label{
				AssignTo: &labelCreateChainState,
				Text:     "新建区块链状态为:",
			},
			PushButton{
				Text: "查询所有地址",
				OnClicked: func() {
					teAllAddress.SetText("")
					//获取teSearchAddress的数据,查询
					walletCollection := NewWalletCollection()
					allAddress := walletCollection.GetAllAddress()
					for index, address := range allAddress {
						fmt.Printf("%d :%s\n", index, address)
						teAllAddress.AppendText(address + "\r\n")
					}
				},
			},
			TextEdit{
				Text:     "所有的地址",
				Row:      20,
				VScroll:  true,
				AssignTo: &teAllAddress,
				MaxSize:  Size{Width: 800, Height: 140},
			},
			TextEdit{
				Text:     "输入查询地址",
				AssignTo: &teSearchAddress,
				MaxSize:  Size{Width: 800, Height: 30},
			},
			PushButton{
				Text: "查询余额",
				OnClicked: func() {
					//获取teSearchAddress的数据,查询
					address := teSearchAddress.Text()
					if address == "" {
						teNewChainState.SetText("地址不正确")
						return
					}
					//查询出结果之后给
					go func() {
						balance := queryValue(address)
						fmt.Printf("查询%s余额\n", address)
						tBalance := strconv.FormatFloat(balance, 'f', -1, 64)
						labelSearchResult.SetText(tBalance)
					}()

				},
			},
			Label{
				AssignTo: &labelSearchResult,
				Text:     "余额为:",
			},
			TextEdit{
				Text:     "交易发起人",
				AssignTo: &teTxFrom,
				MaxSize:  Size{Width: 800, Height: 50},
			},
			TextEdit{
				Text:     "交易接收人",
				AssignTo: &teTxTo,
				MaxSize:  Size{Width: 800, Height: 50},
			},
			TextEdit{
				Text:     "矿工",
				AssignTo: &teTxMiner,
				MaxSize:  Size{Width: 800, Height: 50},
			},
			TextEdit{
				Text:     "金额",
				AssignTo: &teTxAmount,
				MaxSize:  Size{Width: 800, Height: 50},
			},
			TextEdit{
				Text:     "数据",
				AssignTo: &teTxData,
				MaxSize:  Size{Width: 800, Height: 50},
			},
			PushButton{
				Text: "开始交易",
				OnClicked: func() {
					fmt.Println("开始交易")
					from := teTxFrom.Text()
					to := teTxTo.Text()
					miner := teTxMiner.Text()
					amount := teTxAmount.Text()
					data := teTxData.Text()
					fmt.Printf("from %s,to %s,miner %s,amount %s,data %s\n", from, to, miner, amount, data)
					go func() {
						state := startTransfer(from, to, miner, amount, data)
						if state {
							labelTxState.SetText("交易成功")
						} else {
							labelTxState.SetText("交易失败")
						}
					}()
				},
			},
			Label{
				AssignTo: &labelTxState,
				Text:     "交易状态",
			},
		},
	}.Run()
}

//发起一笔转账交易
func startTransfer(from string, to string, miner string, amount string, data string) bool {
	fAmount, err := strconv.ParseFloat(amount, 64)
	state := true
	if err != nil {
		fmt.Println("解析金额失败,err:", err)
		state = false
	}
	fmt.Println("address ", from, "转账给 address ", to, "金额是", amount, "数据是", string(data))
	blockChain := GetBlockChainObject()
	defer blockChain.db.Close()
	var txs []*Transaction
	//创建CoinBase
	coinBase := CreateCoinBase(string(data), miner)
	txs = append(txs, coinBase)
	//创建普通交易
	transaction := NewTransaction(from, to, fAmount, blockChain)
	if transaction != nil {
		txs = append(txs, transaction)
	} else {
		fmt.Println("余额不足,创建交易失败")
		state = false
	}
	blockChain.AppendBlockToChain(txs)
	return state
}

//查询地址账户余额
func queryValue(address string) float64 {
	bc := GetBlockChainObject()
	balance := bc.GetBalance(address)
	fmt.Println("current balance is", balance)
	return balance
}
