#!/bin/bash
rm BlockChain.db
rm BlockChain.db.lock
rm wallets.dat
rm block.exe
go build
#./block