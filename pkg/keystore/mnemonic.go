package keystore

//系统级加密助记词（phase.yml 系统 key + salt）
//用户级二次加密（user key + salt）
//双层加密设计，安全性工业级
import encryptor "tron-scan/pkg/encryptor"

// EncryptMnemonic 系统级加密助记词
func EncryptMnemonic(mnemonic, sysKey, sysSalt string, iter int) (string, error) {
	return encryptor.Encrypt(mnemonic, sysKey+":"+sysSalt, iter)
}

// DecryptMnemonic 系统级解密助记词
func DecryptMnemonic(cipher, sysKey, sysSalt string, iter int) (string, error) {
	return encryptor.Decrypt(cipher, sysKey+":"+sysSalt, iter)
}
