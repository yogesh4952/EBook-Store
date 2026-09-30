package model

type PaymentCallback struct {
	TransactionCode  string `json:"transaction_code"`
	Status           string `json:"status"`
	TotalAmount      string `json:"total_amount"`
	TransactionUUID  string `json:"transaction_uuid"`
	ProductCode      string `json:"product_code"`
	SignedFieldNames string `json:"signed_field_names"`
	Signature        string `json:"signature"`
}
