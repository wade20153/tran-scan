package models

type Coin string

const (
	BTC Coin = "BTC"
	ETH Coin = "ETH"
)

type Wallet struct {
	WalletID      string
	Coin          Coin
	XPRVEncrypted []byte // 数据库存储的密文
	DeriveIndex   uint32
	Addresses     []*ChildAddress
}

type ChildAddress struct {
	Index    uint32
	Address  string
	WalletID string
	Coin     Coin
}
