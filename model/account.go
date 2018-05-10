package model

import "time"

type Account struct {
	Id int64	`json:"id" xorm:"pk autoincr unique" `
	Uid int64 `json:"uid"`
	WalletId int64 `json:"wallet_id"`
	Balance float64 `json:"balance"`
	AccountName string `json:"account_name"`
	Reputation int64 `json:"reputation"`
	Status int `json:"status" xorm:"default 1"`
	VerifyHash string `json:"verify_hash"`
	Version int `xorm:"version"`
	// ... Other information
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	UpdatedAt time.Time	`json:"updated_at" xorm:"updated"`
	DeletedAt time.Time	`json:"deleted_at" xorm:"deleted"`
}

func (a *Account) TableName() string {
	return "b_account"
}
