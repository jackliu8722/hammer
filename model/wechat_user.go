package model

import "time"

type WechatUser struct {
	Id int64	`json:"id" xorm:"pk autoincr unique" `
	UnionId string `json:"union_id" xorm:"varchar(128) notnull unique"`
	NickName string `json:"nick_name" xorm:"varchar(128)"`
	Avatar string `json:"avatar" xorm:"varchar(128)"`
	City string `json:"city"`
	Province string `json:"province"`
	Country string `json:"country"`
	OpenId string `json:"open_id"`
	Gender int `json:"gender"`
	// ... Other information
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	UpdatedAt time.Time	`json:"updated_at" xorm:"updated"`
	DeletedAt time.Time	`json:"deleted_at" xorm:"deleted"`
}

func (w *WechatUser) TableName() string {
	return "b_wechat_user"
}


