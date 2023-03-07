package mymodels

import "gopkg.in/guregu/null.v4"

type DBCryptoRate struct {
	ID          int         `json:"id"`
	DateCreated null.Time   `json:"date_created"`
	Data        null.String `json:"data"`
}

func (DBCryptoRate) TableName() string {
	return "crypto_rate"
}
