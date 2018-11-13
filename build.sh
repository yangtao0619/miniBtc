#! /bin/bash
rm block.exe
rm BlockChain.db
rm BlockChain.db.lock
rm wallets.dat
go build
block.exe generateBc --address "yangtao"