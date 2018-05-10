package model

import (
	"time"
)

type TransactionLog struct {
	Id int64	`json:"id" xorm:"pk autoincr unique"`
	TaskId int64 `json:"task_id"`
	TransactionId string `json:"transaction_id"`
	Uid int64 `json:"uid"`
	Type string `json:"type"`
	//
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	UpdatedAt time.Time	`json:"updated_at" xorm:"updated"`
	DeletedAt time.Time	`json:"deleted_at" xorm:"deleted"`
}

func (t *TransactionLog)TableName() string {
	return "b_transaction_log"
}

func NewTransactionLog(taskId, uid int64, transactionId, t_type string) error {
	trs := TransactionLog{
		TaskId: taskId,
		Uid: uid,
		TransactionId: transactionId,
		Type: t_type,
	}

	_, err := GetEngine().Insert(trs)
	return err
}