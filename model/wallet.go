package model

import "time"

/**
	逻辑上的wallet
	包含用户支付验证的验证字段
 */
type Wallet struct {
	Id int64	`json:"id" xorm:"pk autoincr unique" `
	Uid int64 `json:"uid"`
	Name string `json:"name"`
	Status int64 `json:"status"` // 0 locked; 1 unlocked
	EncodeSecret string `json:"encode_secret"`
	// ... Other information
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	UpdatedAt time.Time	`json:"updated_at" xorm:"updated"`
	DeletedAt time.Time	`json:"deleted_at" xorm:"deleted"`
}

func (w *Wallet)TableName() string {
	return "b_wallet"
}
