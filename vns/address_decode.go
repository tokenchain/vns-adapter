package vns

import (
	"github.com/blocktree/go-owcdrivers/addressEncoder"
	owcrypt "github.com/blocktree/go-owcrypt"
)

//AddressDecoder 地址解析器
type AddressDecoder struct{}

//PrivateKeyToWIF 私钥转WIF
func (decoder *AddressDecoder) PrivateKeyToWIF(priv []byte, isTestnet bool) (string, error) {
	return "", nil

}

//PublicKeyToAddress 公钥转地址
func (decoder *AddressDecoder) PublicKeyToAddress(pub []byte, isTestnet bool) (string, error) {

	cfg := addressEncoder.ETH_mainnetPublicAddress
	if isTestnet {
		cfg = addressEncoder.ETH_mainnetPublicAddress
	}

	//pkHash := btcutil.Hash160(pub)
	//address, err :=  btcutil.NewAddressPubKeyHash(pkHash, &cfg)
	//if err != nil {
	//	return "", err
	//}

	//log.Debug("public key:", common.ToHex(pub))
	publickKey := owcrypt.PointDecompress(pub, owcrypt.ECC_CURVE_SECP256K1)
	//log.Debug("after encode public key:", common.ToHex(publickKey))
	pkHash := owcrypt.Hash(publickKey[1:len(publickKey)], 0, owcrypt.HASH_ALG_KECCAK256)

	//地址添加0x前缀
	address := appendOxToAddress(addressEncoder.AddressEncode(pkHash, cfg))

	return address, nil

}

//RedeemScriptToAddress 多重签名赎回脚本转地址
func (decoder *AddressDecoder) RedeemScriptToAddress(pubs [][]byte, required uint64, isTestnet bool) (string, error) {
	return "", nil
}

//WIFToPrivateKey WIF转私钥
func (decoder *AddressDecoder) WIFToPrivateKey(wif string, isTestnet bool) ([]byte, error) {
	return nil, nil

}
