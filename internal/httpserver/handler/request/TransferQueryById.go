package request

type TransferQueryById struct {
	TxId string `json:"txId" validate:"required"`
}
