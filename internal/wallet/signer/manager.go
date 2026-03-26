package signer

import (
	"sync"
	modlee "tron-scan/models"
)

type Manager struct {
	store     Store
	wallets   map[string]*modlee.Wallet
	mu        sync.RWMutex
	masterKey []byte // AES 密钥
}

func NewManager(store Store, masterKey []byte) *Manager {
	return &Manager{
		store:     store,
		wallets:   make(map[string]*modlee.Wallet),
		masterKey: masterKey,
	}
}

// 启动加载
func (m *Manager) LoadWallet(walletID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	wallet, err := m.store.LoadWallet(walletID)
	if err != nil {
		return err
	}

	// 解密 xprv
	xprv, err := DecryptAESGCM(m.masterKey, wallet.XPRVEncrypted)
	if err != nil {
		return err
	}

	// 内存加载
	wallet.Addresses = []*modlee.ChildAddress{}
	InitAddressPool(wallet, 100) // 默认 100 个地址
	m.wallets[walletID] = wallet

	// 内存擦除 xprv 临时变量
	for i := range xprv {
		xprv[i] = 0
	}
	return nil
}

func (m *Manager) GetWallet(walletID string) *modlee.Wallet {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.wallets[walletID]
}
