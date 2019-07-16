package main

import (
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openw"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/blocktree/vns-adapter/entity"
	"github.com/pkg/errors"
	"strings"
)

//初始化钱包
func InitWalletManager() *openw.WalletManager {

	log.SetLogFuncCall(true)
	tc := openw.NewConfig()

	tc.ConfigDir = configFilePath
	tc.EnableBlockScan = false
	tc.SupportAssets = []string{
		"VNS",
	}
	return openw.NewWalletManager(tc)
}

//创建钱包
func CreateWallet(tm *openw.WalletManager, walletAlias, password string) {

	w := &openwallet.Wallet{Alias: walletAlias, IsTrust: true, Password: password}
	nw, key, err := tm.CreateWallet(app, w)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("wallet:", nw)
	log.Info("key:", key)

}

//获取钱包信息
func GetWalletInfo(tm *openw.WalletManager, walletId string) {

	wallet, err := tm.GetWalletInfo(app, walletId)
	if err != nil {
		log.Error("unexpected error:", err)
		return
	}
	log.Info("wallet:", wallet)
}

//创建账号
//accountAlias 账户别名
//currencySymbol 币总
func CreateAssetsAccount(tm *openw.WalletManager, walletID, accountAlias, currencySymbol, password string) (*openwallet.AssetsAccount, *openwallet.Address, error) {

	account := &openwallet.AssetsAccount{Alias: accountAlias, WalletID: walletID, Required: 1, Symbol: currencySymbol, IsTrust: true}
	account, address, err := tm.CreateAssetsAccount(app, walletID, password, account, nil)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	log.Info("account:", account)
	log.Info("address:", address)

	return account, address, nil

}

//检查资产账户是否属于某个币种
func CheckAccountMatchSymbol(tm *openw.WalletManager, walletID, accountID, symbol string) (bool, error) {

	account, err := tm.GetAssetsAccountInfo(app, walletID, accountID)
	if err != nil {
		log.Error(err)
		return false, err
	}
	if strings.ToUpper(account.Symbol) != strings.ToUpper(symbol) {
		return false, nil
	}
	return true, nil

}

//创建充值地址
func CreateAddress(tm *openw.WalletManager, walletID string, accountID string, count uint64) ([]*entity.DepositAddressResponse, error) {
	address, err := tm.CreateAddress(app, walletID, accountID, count)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//log.Info("address:", address)
	respondAddressList := make([]*entity.DepositAddressResponse, 0)
	for i, w := range address {
		log.Info("address[", i, "] :", w.Address)
		respondAddress := &entity.DepositAddressResponse{AccountID: w.AccountID, WalletID: walletID, Address: w.Address, Symbol: strings.ToLower(w.Symbol)}
		respondAddressList = append(respondAddressList, respondAddress)
	}

	return respondAddressList, nil
}

//获取地址列表
func GetAddressList(tm *openw.WalletManager, walletID string, accountID string) ([]*openwallet.Address, error) {

	list, err := tm.GetAddressList(app, walletID, accountID, 0, -1, false)
	if err != nil {
		log.Error("unexpected error:", err)
		return nil, err
	}
	for i, w := range list {
		log.Info("address[", i, "] :", w.Address)
	}
	log.Info("address count:", len(list))
	return list, nil
}

func SignTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	_, err := tm.SignTransaction(app, rawTx.Account.WalletID, rawTx.Account.AccountID, "12345678", rawTx)
	if err != nil {
		log.Error("SignTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Infof("rawTx: %+v", rawTx)
	return rawTx, nil
}

func VerifyTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	//log.Info("rawTx.Signatures:", rawTx.Signatures)

	_, err := tm.VerifyTransaction(app, rawTx.Account.WalletID, rawTx.Account.AccountID, rawTx)
	if err != nil {
		log.Error("VerifyTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Infof("rawTx: %+v", rawTx)
	return rawTx, nil

}

func SubmitTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	tx, err := tm.SubmitTransaction(app, rawTx.Account.WalletID, rawTx.Account.AccountID, rawTx)
	if err != nil {
		log.Error("SubmitTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Std.Info("tx: %+v", tx)
	log.Info("wxID:", tx.WxID)
	log.Info("txID:", rawTx.TxID)

	return rawTx, nil
}

func CreateTransactionStep(tm *openw.WalletManager, walletID, accountID, to, amount, feeRate string, contract *openwallet.SmartContract) (*openwallet.RawTransaction, error) {

	//err := tm.RefreshAssetsAccountBalance(app, accountID)
	//if err != nil {
	//	log.Error("RefreshAssetsAccountBalance failed, unexpected error:", err)
	//	return nil, err
	//}

	rawTx, err := tm.CreateTransaction(app, walletID, accountID, amount, to, feeRate, "", contract)

	if err != nil {
		log.Error("CreateTransaction failed, unexpected error:", err)
		return nil, err
	}

	return rawTx, nil
}

//转币
func TransferCoin(tm *openw.WalletManager, walletID, accountID, to string, transferAmount string, contract *openwallet.SmartContract) (string, error) {
	var smartContract *openwallet.SmartContract
	if contract == nil || contract.Symbol == "" || contract.Address == "" {
		smartContract = nil
	} else {
		smartContract = contract
	}
	rawTx, err := CreateTransactionStep(tm, walletID, accountID, to, transferAmount, "", smartContract)
	if err != nil {
		return "", err
	}

	log.Std.Info("rawTx: %+v", rawTx)

	_, err = SignTransactionStep(tm, rawTx)
	if err != nil {
		return "", err
	}

	_, err = VerifyTransactionStep(tm, rawTx)
	if err != nil {
		return "", err
	}
	_, err = SubmitTransactionStep(tm, rawTx)
	if err != nil {
		return "", err
	}
	return rawTx.TxID, nil
}

//获取账号eth余额
func GetAssetsAccountBalance(tm *openw.WalletManager, walletID, accountID string) *openwallet.Balance {
	balance, err := tm.GetAssetsAccountBalance(app, walletID, accountID)
	if err != nil {
		log.Error("GetAssetsAccountBalance failed, unexpected error:", err)
		return nil
	}
	log.Info("balance:", balance)
	return balance
}

//获取账号erc20余额
func GetAssetsAccountTokenBalance(tm *openw.WalletManager, walletID, accountID string, contract openwallet.SmartContract) *openwallet.TokenBalance {
	balance, err := tm.GetAssetsAccountTokenBalance(app, walletID, accountID, contract)
	if err != nil {
		log.Error("GetAssetsAccountTokenBalance failed, unexpected error:", err)
		return nil
	}
	log.Info("token balance:", balance.Balance)
	return balance
}

//创建汇总交易
func CreateSummaryTransactionStep(
	tm *openw.WalletManager,
	walletID, accountID, summaryAddress, minTransfer, retainedBalance, feeRate string,
	start, limit int,
	contract *openwallet.SmartContract,
	feeSupportAccount *openwallet.FeesSupportAccount) ([]*openwallet.RawTransactionWithError, error) {

	rawTxArray, err := tm.CreateSummaryRawTransactionWithError(app, walletID, accountID, summaryAddress, minTransfer,
		retainedBalance, feeRate, start, limit, contract, feeSupportAccount)

	if err != nil {
		log.Error("CreateSummaryTransaction failed, unexpected error:", err)
		return nil, err
	}

	return rawTxArray, nil
}

//汇总充值钱包的币
func SummaryCoin(tm *openw.WalletManager, walletID, accountID, feesSupportAccountID, addr string, contract openwallet.SmartContract) error {

	summaryAddress := addr

	GetAssetsAccountBalance(tm, walletID, accountID)
	if contract.Symbol != "" {
		log.Infof("Summary token,info:", contract)
	}
	var summarySmartContract *openwallet.SmartContract = nil
	if contract.Symbol == "" {
		summarySmartContract = nil
	} else {
		summarySmartContract = &contract
	}
	var feesSupport *openwallet.FeesSupportAccount = nil

	//如果合约为空说明是主币汇总
	if summarySmartContract != nil {
		//检查汇总目标账户和费率支持账户是否相等
		if feesSupportAccountID == "" {
			return errors.New("must have feeSupportAccount")
		}

		if accountID == feesSupportAccountID {
			return errors.New(" feeSupportAccountId must not same as summary accountID")
		}

		feesSupport = &openwallet.FeesSupportAccount{
			AccountID: feesSupportAccountID,
			//FixSupportAmount: "0.01",
			FeesSupportScale: "1.3",
		}
	}
	rawTxArray, err := CreateSummaryTransactionStep(tm, walletID, accountID,
		summaryAddress, "", "", "",
		0, 1000, summarySmartContract, feesSupport)
	if err != nil {
		log.Errorf("CreateSummaryTransaction failed, unexpected error: %v", err)
		return err
	}

	//执行汇总交易
	for _, rawTxWithErr := range rawTxArray {

		if rawTxWithErr.Error != nil {
			log.Error(rawTxWithErr.Error.Error())
			continue
		}

		_, err = SignTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return err
		}

		_, err = VerifyTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return err
		}
		//rawTxWithErr.RawTx.IsCompleted = true
		_, err = SubmitTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			log.Errorf("txId=%s err=%s", rawTxWithErr.RawTx.TxID, err.Error())
			return err
		}
	}
	return nil
}
