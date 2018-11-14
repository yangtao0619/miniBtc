package core

import "fmt"

//发起转账给某人
func (cli *CLI) SendToSomeOne(from string, to, miner string, amount float64, data []byte) {
	fmt.Println("address ", from, "转账给 address ", to, "金额是", amount, "数据是", string(data))
	blockChain := GetBlockChainObject()
	//defer blockChain.db.Close()
	var txs []*Transaction
	//创建CoinBase
	coinBase := CreateCoinBase(string(data), miner)
	txs = append(txs, coinBase)
	//创建普通交易
	transaction := NewTransaction(from, to, amount, blockChain)
	if transaction != nil {
		txs = append(txs, transaction)
	} else {
		fmt.Println("余额不足,创建交易失败")
	}
	blockChain.AddBlock(txs)
}

//创建一个钱包
func (cli *CLI) CreateWallet() {
	walletCollection := NewWalletCollection()
	newWalletAddress := walletCollection.CreateWallet()
	if newWalletAddress == "" {
		fmt.Println("地址创建失败")
		return
	}
	fmt.Println("您新的钱包地址是:", newWalletAddress)

}
func (cli *CLI) ListAllAddress() {
	fmt.Println("print address...")
	walletCollection := NewWalletCollection()
	allAddress := walletCollection.GetAllAddress()
	for index, address := range allAddress {
		fmt.Printf("%d :%s\n", index, address)
	}
}
func (cli *CLI) PrintBlockChain() {
	fmt.Println("打印区块链")
	bc := GetBlockChainObject()
	bc.PrintBlockChain()
}
func (cli *CLI) GetBalance(address string) {
	bc := GetBlockChainObject()
	balance := bc.GetBalance(address)
	fmt.Println("current balance is", balance)
}
