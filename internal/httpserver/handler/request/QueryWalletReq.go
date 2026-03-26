package request

type QueryWalletReq struct {
	Address string `json:"address" validate:"required"`
}
