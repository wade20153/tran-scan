package request

type TransferTRXReq struct {
	FromAddr string `json:"fromAddr" validate:"required"`
	ToAddr   string `json:"toAddr" validate:"required"`
	Amount   int64  `json:"amount" validate:"required"`
}
