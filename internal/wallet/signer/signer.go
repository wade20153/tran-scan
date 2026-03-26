package signer

import (
	"fmt"
	"tron-scan/models"
)

func Sign(wallet *models.Wallet, message string) string {
	// 示例签名逻辑（真实用 BIP32 / ECDSA 等库）
	return fmt.Sprintf("Signed[%s]-by-%s", message, wallet.WalletID)
}
