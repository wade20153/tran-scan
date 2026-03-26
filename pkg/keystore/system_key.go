package keystore

import (
	"fmt"
	"strings"

	"tron-scan/pkg/encryptor"
)

// EncryptSystemKeyForUser 用户二次加密系统 key
func EncryptSystemKeyForUser(sysKey, sysSalt, userKey, userSalt string, iter int) (string, error) {
	payload := sysKey + ":" + sysSalt
	return encryptor.Encrypt(payload, userKey+":"+userSalt, iter)
}

// DecryptSystemKeyForUser 用户二次解密系统 key
func DecryptSystemKeyForUser(cipher, userKey, userSalt string, iter int) (string, string, error) {
	plain, err := encryptor.Decrypt(cipher, userKey+":"+userSalt, iter)
	if err != nil {
		return "", "", err
	}
	parts := strings.SplitN(plain, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("数据损坏")
	}
	return parts[0], parts[1], nil
}
