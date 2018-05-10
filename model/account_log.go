package model

import "time"

type AccountLog struct {
	Id int64	`json:"id" xorm:"pk autoincr unique" `
	FromAccount string `json:"from_account"`
	Action string `json:"action"`
	ToAccount string `json:"to_account"`
	Amount float64 `json:"amount"`
	FreezeAmount float64 `json:"freeze_amount"`
	Balance float64 `json:"balance"`
	Status int `json:"status" xorm:"default 0"`
	Memo string `json:"memo"`
	TransactionId string `json:"transaction_id"`
	// ... Other information
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	UpdatedAt time.Time	`json:"updated_at" xorm:"updated"`
}


func (t *AccountLog)TableName() string {
	return "b_account_log"
}

func ActionName(action string) string {
	switch action {
	case "syncBalance":
		return "账户同步"
	case "transferOut":
		return "转出"
	case "transferIn":
		return "转入"
	default:
		return "未知操作"
	}
}