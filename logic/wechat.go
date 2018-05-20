package logic

import (
	"net/url"
	"encoding/json"
	"errors"
	"strconv"
	"fmt"

	"github.com/gin-gonic/gin"
	. "github.com/daoone/hammer/util"
	"github.com/daoone/hammer/model"
	"github.com/eoscanada/eos-go"
	"encoding/hex"
)

const (
	errWxUserData = iota + 4000
	userPrefix = "dao"
)
// 注册送测试dot
var registerAward = true

type WxUserInfo struct {
	Province string `json:"province"`
	City string `json:"city"`
	Country string `json:"country"`
	UnionId string `json:"unionId"`
	Gender int	`json:"gender"`
	NickName string `json:"nickName"`
	AvatarUrl string `json:"avatarUrl"`
	OpenId string `json:"openId"`
}

var wxRedisPrefix = "wx"
/**
	微信用户登录session，使用status字段控制session可用性
 */
func WxCheckLogin(code, data, iv string) (string, error) {

	// 解码用户信息
	sessionKey, _ := WxGetSessionKey(code)
	token := Md5(sessionKey)
	redisClient := GetRedis()
	rs := redisClient.WithPrefix(wxRedisPrefix)
	defer redisClient.Close()

	// session已经登录了
	if val, _ := rs.HGET(token, "status");val == "1" {
		rs.EXPIRE(token, 600)
		return token, nil
	}

	// 没有已登录状态，重新解码用户信息
	realEcrypt, _ := url.QueryUnescape(data)
	wxe := WxBizDataCrypt{
		AppID: GetConfig("wechat", "appId"),
		SessionKey: sessionKey,
	}
	jsstr, err := wxe.Decrypt(realEcrypt, iv, true)
	if err != nil {
		DoLog(err.Error(), "wx_error")
		return "", errors.New("data decrypt error")
	}

	var wu WxUserInfo
	json.Unmarshal(jsstr.([]byte), &wu)
	// 查看是否已经注册
	err = wu.WxAutoRegister()
	if err != nil {
		return "", err
	}

	// 设置redis
	rs.HSET(token, "status", "1")
	rs.HSET(token, "union_id", wu.UnionId)

	u := &model.User{UnionId: wu.UnionId}
	model.GetEngine().Get(u)
	if u.Uid != 0 {
		rs.HSET(token, "uid", strconv.FormatInt(u.Uid,10))
	}

	rs.EXPIRE(token, 600)
	return token, nil
}

// 自动注册
func (wu *WxUserInfo)WxAutoRegister() (error) {
	wechatUser := &model.WechatUser{}
	_, err := model.GetEngine().Where("union_id=?", wu.UnionId).Get(wechatUser)
	if err != nil {
		return  err
	}

	// 事务注册用户
	if wechatUser.Id == 0 {
		session := model.GetEngine().NewSession()
		defer session.Close()
		err = session.Begin()
		newWechatUser := model.WechatUser{
			UnionId: wu.UnionId,
			NickName: FilterEmoji(wu.NickName),
			City: wu.City,
			Province: wu.Province,
			Country: wu.Country,
			OpenId: wu.OpenId,
			Gender: wu.Gender,
			Avatar: wu.AvatarUrl,
		}

		_, err := session.Insert(newWechatUser)
		if err != nil {
			session.Rollback()
			return err
		}
		newUser := model.User{}
		newUser.GenRandomUserName()
		newUser.GenMd5Passwd()
		newUser.UnionId = newWechatUser.UnionId
		if _, err = session.Insert(newUser); err != nil {
			session.Rollback()
			return err
		}
		err = session.Commit()
	}

	return nil
}

func WxGetUserDotAndRp(unionId string) gin.H {
	account := &model.Account{}
	_, err := model.GetEngine().Cols("b_account.*").
		Join("LEFT", "b_user","b_account.uid=b_user.uid").
		Where("b_user.union_id=?", unionId).Get(account)
	if err != nil {
		DoLog(err.Error(),"logic_wechat")
		return gin.H{"error": true}
	}

	if account.Id == 0 {
		return gin.H{"success":true, "dot":0, "rp": 0}
	}

	return gin.H{"success":true, "dot": account.Balance, "rp": account.Reputation}
}

func WxGetUserDot(unionId string) gin.H {
	account := &model.Account{}
	_, err := model.GetEngine().Cols("b_account.*").
		Join("LEFT", "b_user","b_account.uid=b_user.uid").
			Where("b_user.union_id=?", unionId).Get(account)
	if err != nil {
		DoLog(err.Error(),"logic_wechat")
		return gin.H{"error": true}
	}

	if account.Id == 0 {
		return gin.H{"success":true, "dot":0}
	}

	e := GetEos()
	asset, err := e.GetCurrencyBalance(eos.AN(account.AccountName), TOKEN_SYM, eos.AN(DOT_CONTRACT))
	fmt.Println(asset)
	if len(asset) <= 0 {
		return gin.H{"success":true, "dot": 0}
	}
	if err != nil {
		DoLog("Query balance error|" + err.Error(), "logic_wechat")
		return gin.H{"success":true, "dot": account.Balance}
	}

	if account.Balance != DotUnitToFloat(asset[0].Amount) {
		fmt.Println(account)
		account.Balance = 	DotUnitToFloat(asset[0].Amount)
		session := model.GetEngine().NewSession()
		defer session.Close()
		session.Begin()
		_, err1 := session.ID(account.Id).
			Cols("balance", "version", "updated_at").Update(account)
		if err1 != nil {
			DoLog(err1.Error(), "logic_wechat")
			session.Rollback()
		}
		_, err2 := session.Insert(model.AccountLog{
			FromAccount: account.AccountName,
			Action: "syncBalance",
			ToAccount: account.AccountName,
			Amount: account.Balance,
			FreezeAmount: 0,
			Balance: account.Balance,
			Status: 1,
			Memo: "Synchronize from node",
		})
		if err2 != nil {
			DoLog(err2.Error(), "logic_wechat")
			session.Rollback()
		}

		session.Commit()
	}

	return gin.H{"success":true, "dot": account.Balance}
}

func WxGetUserRp(unionId string) gin.H {
	account := &model.Account{}
	_, err := model.GetEngine().Cols("b_account.*").
		Join("LEFT", "b_user","b_account.uid=b_user.uid").
		Where("b_user.union_id=?", unionId).Get(account)
	if err != nil {
		DoLog(err.Error(),"logic_wechat")
		return gin.H{"error": true}
	}

	if account.Id == 0 {
		return gin.H{"success":true, "rp":0}
	}

	return gin.H{"success":true, "rp": account.Reputation}
}

func WxGetUserWallet(unionId string) gin.H {
	wallet := &model.Wallet{}
	_, err := model.GetEngine().Cols("b_wallet.*").
		Join("LEFT", "b_user", "b_wallet.uid=b_user.uid").
		Where("b_user.union_id=?", unionId).Get(wallet)
	if err != nil {
		DoLog(err.Error(),"logic_wechat")
		return gin.H{"error": true}
	}

	if wallet.Id == 0 {
		return gin.H{"success":true, "wallet": ""}
	}

	return gin.H{"success":true, "wallet": wallet.Name, "status": wallet.Status}
}

func WxChangeWalletStatus(unionId string, status int) gin.H {
	return gin.H{ "success": true }
}

// 第一次创建钱包，附带一个默认账户
func WxCreateUserWallet(unionId, passcode string) gin.H {
	errorH := WxError("创建钱包失败")
	// 检测用户是否已经有钱包
	if WxCheckHasWallet(unionId) {
		return WxError("钱包已存在，请直接创建账户")
	}

	// 创建钱包，返回钱包密码
	name := autoName()
	e := GetEos()
	walletKey, err := e.CreateWallet(name)
	if err != nil {
		return errorH
	}
	// 关联用户
	user := &model.User{UnionId: unionId}
	_, err = model.GetEngine().Get(user)
	if err != nil || user.Uid == 0 {
		DoLog("wallet create but not related, need to be deleted:" + name + "|" + unionId, "logic_wechat")
		return errorH
	}
	// 密文aes密码 start
	pcode := []byte(passcode + passcode)
	ecode, err := Encrypt([]byte(walletKey), pcode)
	// 密文aes密码 over
	if err != nil {
		DoLog(err.Error(), "logic_wechat")
		return errorH
	}

	newWallet := model.Wallet{
		Uid: user.Uid,
		Name: name,
		Status: 1,
		EncodeSecret: fmt.Sprintf("%x", ecode), // 解锁需要密码但是不能明文存储
	}

	_, err = model.GetEngine().Insert(newWallet)
	model.NewTransactionLog(0, user.Uid, "", "create_wallet")

	if err != nil {
		DoLog(err.Error(), "logic_wechat")
		return errorH
	}
	// 创建默认账户
	prk, err := WxCreateUserAccount(name, name, passcode, user.Uid)
	//prk, err := "5Jb83EU35GpKYjG7rVJxReoVxqTwanaN2qFZa6RNdt4cwztdAPU", nil
	if err != nil {
		DoLog(err.Error(), "logic_wechat")
		return WxError("Fail to create account.")
		//return gin.H{"error": 2, "message":}
	}
	// 添加测试金额
	if registerAward {
		e.IssueDot(name, "first time reward", 0.1)
	}

	return gin.H{"success":true, "accountKey": prk, "name": name}
}

// @todo 支持单独创建账户行为
func WxPreCreateAccount(unionId, accountName string) {

}

func WxCreateUserAccount(accountName, walletName, passcode string, uid int64) (string, error) {
	// 创建账户并导入到钱包
	e := GetEos()
	prk, txId, err := e.CreateAccount(accountName)
	if err != nil {return "", err}
	err = e.WalletImportKey(walletName, prk)
	if err != nil {return "", err}
	// from wallet
	wallet := &model.Wallet{Name: walletName}
	_, err = model.GetEngine().Get(wallet)
	if err != nil { return "",err }
	if wallet.Id == 0 { return "", errors.New("no wallet found") }

	// 账户的支付口令摘要
	hash := Md5(passcode + strconv.FormatInt(uid,10))
	// 关联account
	ac := model.Account{
		Uid: uid,
		WalletId: wallet.Id,
		Balance: 0,
		AccountName: accountName,
		Reputation: 5,
		Status: 1,
		VerifyHash: hash,
	}
	_, err = model.GetEngine().Insert(ac)
	if err != nil {
		return "", err
	}
	// 记录transaction
	model.NewTransactionLog(0, uid, txId, "create_account")

	return prk, nil
}

// @todo 支持用户导入账户密钥
func WxImportUserAccount() {

}

func WxDoTransfer(unionId, to, amount, memo, passcode string) gin.H{
	engine := model.GetEngine()
	// 验证口令和用户
	user := &model.User{UnionId: unionId}
	_, err := engine.Get(user)
	if err != nil || user.Uid == 0 {
		return WxError("User does not exist.")
	}

	hash := Md5(passcode + strconv.FormatInt(user.Uid,10))

	account := &model.Account{Uid: user.Uid}
	_, err = engine.Get(account)
	if err != nil || account.Id == 0 {
		return WxError("账户不存在")
	}
	if account.VerifyHash != hash {
		return WxError("口令不正确")
	}

	wallet := &model.Wallet{Id: account.WalletId}
	_, err = engine.Get(wallet)
	if err != nil || account.Id == 0 {
		return WxError("钱包不存在")
	}
	// 口令正确解锁钱包
	e := GetEos()
	byteEncode, _ := hex.DecodeString(wallet.EncodeSecret)
	realPass, err := Decrypt(byteEncode, []byte(passcode + passcode))
	if err != nil || e.WalletUnlock(wallet.Name, string(realPass)) != nil {
		//fmt.Println(e.WalletUnlock(wallet.Name, string(realPass)))
		return WxError("解锁钱包出错")
	}

	var flo float64
	flo, err = strconv.ParseFloat(amount, 64)
	if err != nil || account.Balance < flo{
		// 金额出错或与余额不足
		return WxError("DOT余额不足")
	}

	toAccount := &model.Account{AccountName: to}
	_, err = engine.Get(toAccount)
	if err != nil || toAccount.Id == 0 {
		return WxError("对方账户不存在")
	}


	txId, err := e.TransferDot(account.AccountName, to, memo, wallet.Name, flo)
	if err != nil {
		fmt.Println(err)
		return WxError("请稍后重试")
	}

	model.NewTransactionLog(0, user.Uid, txId, "transfer_dot")
	session := model.GetEngine().NewSession()
	defer session.Close()
	err = session.Begin()
	// 扣款 加款
	account.Balance = account.Balance - flo
	toAccount.Balance = toAccount.Balance + flo
	_, err1 := session.ID(account.Id).Cols("balance","version","updated_at").Update(account)
	_, err2 := session.ID(toAccount.Id).Cols("balance","version","updated_at").Update(toAccount)
	if err1 != nil || err2 != nil {
		session.Rollback()
		DoLog("update balance error", "logic_wechat")
		return WxError("请稍后重试")
	}
	//fmt.Println("--------", to, account.AccountName)
	// 用户记录转出
	_, err1 = session.Insert(model.AccountLog{
		FromAccount: account.AccountName,
		Action: "transferOut",
		ToAccount: to,
		Amount: flo,
		FreezeAmount: 0,
		Balance: (account.Balance - flo),
		Status: 1,
		Memo: memo,
		TransactionId: txId,
	})
	// 用户记录转入
	_, err2 = session.Insert(model.AccountLog{
		FromAccount: to,
		Action: "transferIn",
		ToAccount: account.AccountName,
		Amount: flo,
		FreezeAmount: 0,
		Balance: (toAccount.Balance + flo),
		Status: 1,
		Memo: memo,
		TransactionId: txId,
	})

	if err1 != nil || err2 != nil {
		session.Rollback()
		DoLog("update accountlog error", "logic_wechat")
		return WxError("请稍后重试")
	}

	if session.Commit() != nil {
		return WxError("请稍后重试")
	}

	e.WalletLock(wallet.Name)

	return gin.H{"success": true}
}

func WxGetUserAccountLog(unionId, page string) gin.H {
	account := &model.Account{}
	_, err := model.GetEngine().Cols("b_account.*").
		Join("LEFT", "b_user", "b_account.uid=b_user.uid").
		Where("b_user.union_id=?", unionId).Get(account)
	if err != nil || account.Id == 0{
		return WxError("账户不存在")
	}
	logs := make([]model.AccountLog, 0)
	limit, _ := strconv.Atoi(page)
	err = model.GetEngine().
		Where("from_account=?", account.AccountName).
		Limit(10, limit*10).
		OrderBy("created_at DESC").
		Find(&logs)
	if err != nil {
		return WxError("日志不存在")
	}
	res := make([]map[string]string, 0)
	for _, v := range logs {
		a := make(map[string]string)
		a["id"] = strconv.FormatInt(v.Id, 10)
		a["action"] = model.ActionName(v.Action)
		a["amount"] = fmt.Sprintf("%.4f", v.Amount)
		a["to_user"] = v.ToAccount
		a["created_at"] = v.CreatedAt.Format("2006/01/02 15:04:05")
		res = append(res, a)
	}

	return gin.H{"success":true, "logs": res}
}

func WxGetUserAccountLogDetail(unionId, tid string) gin.H {
	account := &model.Account{}
	_, err := model.GetEngine().Cols("b_account.*").
		Join("LEFT", "b_user", "b_account.uid=b_user.uid").
		Where("b_user.union_id=?", unionId).Get(account)
	if err != nil || account.Id == 0{
		return WxError("账户不存在")
	}
	id, _ := strconv.ParseInt(tid, 10, 64)
	log := &model.AccountLog{Id: id}
	_, err = model.GetEngine().Where("from_account=?", account.AccountName).Get(log)
	if err != nil {
		return WxError("账户不存在")
	}

	res := make(map[string]interface{})
	// 非链上交易
	if log.TransactionId == "" {
		res["transFrom"] = log.FromAccount
		res["transTo"] = log.ToAccount
		res["dot"] = log.Amount
		res["transTime"] = log.CreatedAt
		res["transId"] = "-非链上交易-"
		res["contract"] = "-非链上交易-"
		res["memo"] = log.Memo

		return gin.H{"success": true, "detail": res}
	}

	txid := log.TransactionId
	transaction, err := GetEos().GetAccountTransaction(txid)
	actionData := transaction.Data.(map[string]interface{})

	if err != nil {
		return WxError("获取记录详细失败")
	}

	res["transFrom"] = actionData["from"]
	res["transTo"] = actionData["to"]
	res["dot"] = actionData["quantity"]
	res["transTime"] = log.CreatedAt
	res["transId"] = txid
	res["contract"] = transaction.Account
	res["memo"] = actionData["memo"]

	return gin.H{"success": true, "detail": res}
}


func WxCheckHasWallet(unionId string) bool {
	wallet := &model.Wallet{}
	_, err := model.GetEngine().Cols("b_wallet.*").
		Join("LEFT", "b_user", "b_wallet.uid=b_user.uid").
		Where("b_user.union_id=?", unionId).Get(wallet)
	if err != nil {
		DoLog(err.Error(),"logic_wechat")
		return false
	}

	if wallet.Id == 0 {
		return false
	}

	return true
}

func autoName() string {
	return userPrefix + MakeRandomStrLower(9)
}
