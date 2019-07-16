package main

import (
	"github.com/astaxie/beego/config"
	"github.com/blocktree/openwallet/openw"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/blocktree/vns-adapter/vns"
	"path/filepath"
)

var (
	app            = "vns-adapter"
	configFilePath = filepath.Join("conf")
)

func main() {
	absFile := filepath.Join(configFilePath, "VNS.ini")
	c, err := config.NewConfig("ini", absFile)
	if err != nil {
		panic(err.Error())
	}

	openw.RegAssets(vns.Symbol, vns.NewWalletManager(), c)
	tm := InitWalletManager()
	walletId := "WKcxAcyQyY2zUpd9EVWGJ2UM5rwgM9Bqup"
	accountId := "82uzz91iDN98rhQ2szNEZfxhSb412GJxL1BFnkJZzerK"
	//CreateWallet(tm,"vns-main-wallet","12345678")

	//CreateAssetsAccount(tm,walletId,"xxx",vns.Symbol,"12345678")
	//CreateAddress(tm, walletId, accountId, 10)
	//GetAddressList(tm,walletId,accountId)
	//GetAssetsAccountBalance(tm,walletId,accountId)
	//tw := vns.Client{
	//	BaseURL: "http://192.168.31.142:2000",
	//	Debug:   true,
	//}
	//r, _:= tw.GetAddrBalance("0x142b097b802b5224f9d4bfc93db189b4f4621df2")
	// log.Error(r.String())
	//TransferCoin(tm,walletId,accountId,"0xb6758e52b9a75bb345892f06fcb98d65aee2585c","10",nil)
	var contract openwallet.SmartContract

	SummaryCoin(tm, walletId, accountId, "", "0x696f725B176905e4300f45CABC91b8afB8db644c", contract)
}
