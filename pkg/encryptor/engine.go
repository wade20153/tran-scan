package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

// deriveKey 使用 PBKDF2 从 password + salt 派生出固定长度的密钥
//
// 参数说明：
// - password: 原始秘密（用户密码 / 助记词 / master key）
// - salt: 随机盐（必须与密文一起存储，用于防止彩虹表攻击）
// - iterations: 迭代次数（越大越安全，但越慢，常用 100_000）
// - keyLen: 输出密钥长度（AES-256 为 32 字节）
//
// 返回值：
// - 派生后的密钥 []byte，可直接用于 AES
func deriveKey(password string, salt []byte, iterations int, keyLen int) []byte {
	return pbkdf2.Key([]byte(password), salt, iterations, keyLen, sha256.New)
}

// Encrypt 使用 AES-GCM 加密字符串
// 返回的密文格式：
//
//	v1:salt:nonce:ciphertext
//
// 各字段说明：
// - v1        : 版本号（方便将来算法升级）
// - salt      : PBKDF2 使用的随机盐（base64）
// - nonce     : GCM 使用的随机 nonce（base64）
// - ciphertext: 加密后的数据（base64）
func Encrypt(plainText, password string, iterations int) (string, error) {
	// 生成随机 salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	// 派生 key
	key := deriveKey(password, salt, iterations, 32)
	// 3️⃣ 创建 AES block（AES 本身不负责模式）
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	// 4️⃣ 使用 GCM 模式（同时提供加密 + 完整性校验）
	// GCM 是 AEAD（Authenticated Encryption）
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	// 5️⃣ 生成随机 nonce（长度由 GCM 决定，通常是 12 字节）
	// nonce 必须唯一，但可以公开
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	// 6️⃣ 执行加密
	// Seal 会同时完成：
	// - 加密
	// - 生成认证标签（防止篡改）
	cipherText := gcm.Seal(
		nil,               // dst（nil 表示新建 slice）
		nonce,             // 随机 nonce
		[]byte(plainText), // 明文
		nil,               // additional data（AAD，可选）
	)
	// 7️⃣ 拼接最终输出
	// 所有二进制数据统一 base64，便于存数据库 / JSON
	return fmt.Sprintf(
		"v1:%s:%s:%s",
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(nonce),
		base64.StdEncoding.EncodeToString(cipherText),
	), nil

}

// Decrypt 解密 Encrypt 生成的 v1:salt:nonce:ciphertext 格式密文
func Decrypt(cipherText, password string, iterations int) (string, error) {

	// 1️⃣ 按分隔符拆分密文
	parts := strings.Split(cipherText, ":")
	if len(parts) != 4 {
		return "", fmt.Errorf("非法密文格式")
	}

	// 2️⃣ 解析 salt
	salt, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	// 3️⃣ 解析 nonce
	nonce, err := base64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return "", err
	}

	// 4️⃣ 解析密文数据
	data, err := base64.StdEncoding.DecodeString(parts[3])
	if err != nil {
		return "", err
	}

	// 5️⃣ 使用相同 password + salt + iterations 派生密钥
	key := deriveKey(password, salt, iterations, 32)

	// 6️⃣ 创建 AES block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 7️⃣ 创建 GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 8️⃣ 解密并校验完整性
	// 如果 password 错误 / 数据被篡改，这里会直接返回 error
	plain, err := gcm.Open(
		nil,
		nonce,
		data,
		nil,
	)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}
