# miniBtc
实现一个小型的比特币客户端,功能包括

+ 1.block generateBc --address Address 生成一条新的区块链,Address表示创建人的地址
+ 2.block printBc 打印出当前的区块链数据
+ 3.block getBalance --address Address 得到Address地址的账户余额
+ 4.block send from ADDRESS1 to ADDRESS2 miner ADDRESS3 amount MoneyValue data DATA
+ 5.block createWallet 创建一个钱包
+ 6.block listAllAddress 列出所有的地钱包地址

### 使用方式 ###

1.本客户端有命令行和图形界面两种运行方式,可以在main.go中按照注释使用

```go
	//命令行模式
	/*	cli := core.CLI{}
		cli.Run()*/

	//gui模式,默认模式
	core.StartGui()
```

2.需要先生成若干地址,使用其中一个地址创建一个区块链,之后就可以使用该地址进行转账,目前挖矿奖励是12.5

3.成功运行的截图

![avatar](https://t1.aixinxi.net/o_1csag6sgh1jfj1ep3g0rqqh1ctia.png-w.jpg)