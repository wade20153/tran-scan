package request

type DeriveSubWalletReq struct {
	Index uint32 `json:"index" validate:"required"`
}
