package model

import(
	"math/rand"
	"errors"
	"time"
	"fmt"
	"github.com/daoone/hammer/util"
)

type User struct {
	Uid int64 `json:"uid" xorm:"pk"`
	Username  string    `json:"username" xorm:"varchar(25) notnull unique 'username'"`
	Passwd    string    `json:"passwd"`
	Email     string    `json:"email"`
	Phone    string    `json:"phone" xorm:"varchar(25) unique"`
	Status	  int 		`json:"status" xorm:"default 1"` // 用户状态
	LoginTime time.Time `json:"login_time" xorm:"<-"` // 最后登录时间
	Passcode string		`json:"passcode"` // 生成密码的随机数
	UnionId 	string 	`json:"union_id"`
	// ... Other information
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	UpdatedAt time.Time	`json:"updated_at" xorm:"updated"`
	DeletedAt time.Time	`json:"deleted_at" xorm:"deleted"`
}

func (u *User) TableName() string {
	return "b_user"
}

type UserInfo struct {
	Id int64	`json:"id" xorm:"pk"`
	Uid int64	`json:"uid"`
	RealName string `json:"realname"`
}

func (u *UserInfo) TableName() string {
	return "b_user_info"
}

func (u *User) GenMd5Passwd() error {
	rand.Seed(time.Now().UnixNano())
	if u.Passwd == "" {
		return errors.New("password is required")
	}
	u.Passcode = fmt.Sprintf("%x", rand.Int31())
	u.Passwd = util.Md5(u.Passwd + u.Passcode)
	return nil
}

func (u *User) GenRandomUserName() {
	rand.Seed(time.Now().UnixNano())
	u.Username = "daooner_" + util.MakeRandomStr(6)
	u.Passwd = util.MakeRandomStr(8)
	fmt.Println(u.Passwd)
}