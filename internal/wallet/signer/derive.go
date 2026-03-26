package signer

import "fmt"
import wall "tron-scan/models"

func DeriveAddress(wallet *wall.Wallet, index uint32) string {
	// 示例派生逻辑
	return fmt.Sprintf("%s-%d", wallet.WalletID, index)
}

func InitAddressPool(wallet *wall.Wallet, count uint32) {
	wallet.Addresses = make([]*wall.ChildAddress, count)
	for i := uint32(0); i < count; i++ {
		wallet.Addresses[i] = &wall.ChildAddress{
			Index:    i,
			Address:  DeriveAddress(wallet, i),
			WalletID: wallet.WalletID,
			Coin:     wallet.Coin,
		}
	}
}
