package util

import(
	//"bytes"
	//"io/ioutil"
	//"encoding/json"
	//"net/http"
	"errors"
	"strconv"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/token"
	"github.com/eoscanada/eos-go/system"
	"fmt"
)

const(
	WALLET_NAME = "daoone"
	CREATOR = "daoone"
	BOUNTY_CONTRACT = "bounty"
	DOT_CONTRACT = "daoone.token"
	//TOKEN_SYM = "DOT"
	//DOT_CONTRACT = "eosio.token"
	TOKEN_SYM = "DOT"
)

var eosInstance EosInstance

type EosInstance struct {
	*eos.API
	httpEndPoint string
	rpcVersion string
}

type NewTask struct {
	Owner eos.AccountName `json:"owner"`
	Approver eos.AccountName `json:"approver"`
	Contributor eos.AccountName `json:"contributor"`
	Desc string `json:"desc"`
	Bounty int64 `json:"bounty"`
}


func GetEos() *EosInstance {
	if eosInstance.API == nil {
		url := GetConfig("eos","httpEndpoint")
		chainId := make([]byte, 32)
		eosInstance = EosInstance{
			eos.New(url, chainId),
			GetConfig("eos","httpEndpoint"),
			GetConfig("eos","rpcVersion"),
		}
		return &eosInstance
	}

	return &eosInstance
}

// 创建用户
func (e *EosInstance)CreateAccount(accountName string) (string, string, error) {
	privateKey, err := ecc.NewRandomPrivateKey()
	if err != nil {
		DoLog(err.Error(), "eos_api")
		return "", "", err
	}
	publicKey := privateKey.PublicKey()
	e.SetSigner(eos.NewWalletSigner(e.API, WALLET_NAME))
	action := system.NewNewAccount(CREATOR, eos.AccountName(accountName), publicKey)

	tx, err := e.SignPushActions(action)
	if err != nil {
		DoLog(err.Error(), "eos_api")
		return "", "", err
	}

	return privateKey.String(), tx.TransactionID, nil
}

// 创建唯一的钱包
// 每个用户只能创建一个钱包，钱包名字关联用户表
// 不存储钱包密钥，创建钱包无法从外部调用
func (e *EosInstance)CreateWallet(walletName string) (string, error) {
	res, err := DoPost(fmt.Sprintf("%s/%s/%s", e.httpEndPoint, e.rpcVersion, "wallet/create"), walletName)
	if err != nil {
		DoLog(err.Error(), "eos_api")
		return "", err
	}
	str, _ := strconv.Unquote(string(res))

	return str, nil
}

func (e *EosInstance)WalletUnlock(name, secret string) error {
	data := make([]string,0)
	data = append(data, name)
	data = append(data, secret)
	fmt.Println(data)
	_, err := DoPost(fmt.Sprintf("%s/%s/%s", e.httpEndPoint, e.rpcVersion, "wallet/unlock"), data)
	if err != nil {
		return err
	}

	return nil
}

func (e *EosInstance)WalletLock(name string) error {
	_, err := DoPost(fmt.Sprintf("%s/%s/%s", e.httpEndPoint, e.rpcVersion, "wallet/lock"), name)
	if err != nil {
		return err
	}

	return nil
}

func (e *EosInstance)GetAccountTransaction(id string) (eos.Action, error) {
	res, err := e.API.GetTransaction(id)
	if err != nil {
		return eos.Action{}, err
	}

	if len(res.Transaction.Transaction.Actions) > 0 {
		return *res.Transaction.Transaction.Actions[0], nil
	}

	return eos.Action{}, errors.New("No such Transaction")
}

func (e *EosInstance)GetAccountTransactionsRaw(accountName string, startNum, count int) (interface{},error) {
	payload := map[string]interface{}{
		"account_name": accountName,
		"skip_seq": startNum,
		"num_seq": count,
	}
	resp, err := DoPost(fmt.Sprintf("%s/%s/%s", e.httpEndPoint, e.rpcVersion, "account_history/get_transactions"), payload)
	if err != nil {
		return "",err
	}
	return resp, nil
}

func (e *EosInstance)WalletImport(privKey string) (error) {
	return e.WalletImportKey(WALLET_NAME, privKey)
}

// 调用合约生成bounty
func (e *EosInstance)CreateTask(owner, approver, contributor, desc string, bounty int64) (string, error) {
	a := &eos.Action{
		Account: eos.AN(BOUNTY_CONTRACT),
		Name: eos.ActN("create"),
		Authorization: []eos.PermissionLevel{
			{eos.AN(owner), eos.PN("active")},
		},
		ActionData: eos.NewActionData(NewTask{
			Owner: eos.AN(owner),
			Approver: eos.AN(approver),
			Contributor: eos.AN(contributor),
			Desc: desc,
			Bounty: bounty,
		}),
	}

	e.SetSigner(eos.NewWalletSigner(e.API, WALLET_NAME))
	tx, err := e.SignPushActions(a)
	if err != nil {
		DoLog(err.Error(), "eos_api")
		return "", err
	}

	return tx.TransactionID, nil
}


func (e *EosInstance)TransferDot(from, to, memo, walletName string, quantity float64) (string, error){
	s := fmt.Sprintf("%.4f", quantity)
	asset, err := eos.NewAsset(s + " " + TOKEN_SYM)
	if err != nil { return "", err }
	a := &eos.Action{
		Account: eos.AN(DOT_CONTRACT),
		Name:    eos.ActN("transfer"),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AN(from), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(token.Transfer{
			From:     eos.AN(from),
			To:       eos.AN(to),
			Quantity: asset,
			Memo:     memo,
		}),
	}

	e.SetSigner(eos.NewWalletSigner(e.API, walletName))
	tx, err := e.SignPushActions(a)
	if err != nil {
		DoLog(err.Error(), "eos_api")
		return "", err
	}
	return tx.TransactionID, nil
}

func (e *EosInstance)IssueDot(to, memo string, quantity float64) (string, error){
	s := fmt.Sprintf("%.4f", quantity)
	asset, err := eos.NewAsset(s + " " + TOKEN_SYM)
	a := &eos.Action{
		Account: eos.AN(DOT_CONTRACT),
		Name:    eos.ActN("issue"),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AN(CREATOR), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(token.Issue{
			To:       eos.AN(to),
			Quantity: asset,
			Memo:     memo,
		}),
	}

	e.SetSigner(eos.NewWalletSigner(e.API, WALLET_NAME))
	tx, err := e.SignPushActions(a)
	if err != nil {
		DoLog(err.Error(), "eos_api")
		return "", err
	}
	return tx.TransactionID, nil
}

