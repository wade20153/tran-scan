package security

// GetSystemKeyFromKMS 从 KMS/HSM 获取系统密钥
func GetSystemKeyFromKMS() (key, salt string, err error) {
	// 🔹 TODO: AWS KMS / HashiCorp Vault / HSM 代替本地 phase.yml
	// 业务调用无需修改，保持 encryptor/keystore 接口一致
	return "", "", nil
}
