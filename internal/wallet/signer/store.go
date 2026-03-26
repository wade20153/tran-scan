package signer

import wall "tron-scan/models"

type Store interface {
	SaveWallet(wallet *wall.Wallet) error
	LoadWallet(walletID string) (*wall.Wallet, error)
}

// 简单内存实现示例
type MemoryStore struct {
	wallets map[string]*wall.Wallet
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{wallets: make(map[string]*wall.Wallet)}
}

func (s *MemoryStore) SaveWallet(wallet *wall.Wallet) error {
	s.wallets[wallet.WalletID] = wallet
	return nil
}

func (s *MemoryStore) LoadWallet(walletID string) (*wall.Wallet, error) {
	w, ok := s.wallets[walletID]
	if !ok {
		return nil, nil
	}
	return w, nil
}
