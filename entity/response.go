package entity

//地址返回实体
type DepositAddressResponse struct {
	WalletID  string `json:"walletId" storm:"index"` //钱包ID
	AccountID string `json:"accountId"`              //账号ID
	Address   string `json:"address" storm:"id"`     //地址字符串
	Symbol    string `json:"currencySymbol"`
	TxId      string
	//币种类别
}

//充值通知实体
type DepositNotify struct {
	WxId            string `json:"wxId"`            //钱包生成的唯一ID
	AccountID       string `json:"accountId"`       //账号ID
	Address         string `json:"depositAddress"`  //地址字符串
	TxId            string `json:"txId"`            //地址字符串
	BlockNumber     uint64 `json:"blockNumber"`     //区块数字
	BlockHash       string `json:"blockHash"`       //区块HASH
	DepositAmount   string `json:"depositAmount"`   //地址字符串
	Symbol          string `json:"currencySymbol"`  //币种类别
	ConfirmTimes    int64  `json:"confirmTimes"`    //确认次数
	IsContract      bool   `json:"isContract"`      //是否合约转账
	ContractAddress string `json:"contractAddress"` //合约地址

}
