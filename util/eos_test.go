package util_test

import(
	//"testing"
	//"fmt"
	//"github.com/tidwall/gjson"
	eos "github.com/eoscanada/eos-go"
	"net/url"
	//"github.com/eoscanada/eos-go/ecc"
	//"github.com/daoone/hammer/util"
	//"strconv"
	//"math/rand"
	//"time"
)

 /**
 POST /v1/chain/get_required_keys HTTP/1.0\r

{"transaction":{
 "ref_block_num":1861,
 "ref_block_prefix":4023696733,
 "expiration":"2018-03-20T23:15:59",
 "scope":["eos","inita"],"read_scope":[],
 "messages":[{"code":"eos",
 "type":"newaccount","authorization":
 [{"account":"inita","permission":"active"}],
 "data":"000000000093dd740000003499ab2afd0100000001033e6263844877f0235259e1a0feb944ea23c48eef4992a3b6bb85a58208a277f201000001000000010254b362d60ff4817a744ed0443c9ef71808012d069df21c5891b6977def4cb070010000010000000001000000000093dd7400000000a8ed32320100010000000000000004454f5300000000"}],
 "signatures":[]},"available_keys":["EOS5XnvKE4f5RRZHr713VAEbmeUkWfQ4f8aHj4ck2sEqxo6WkDbWm","EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV","EOS7Ji4678S1JSk774akqxvDxbQkGynQknzmCxYCYx5d94wdqoxb6"]}
  */

func getEos() *eos.API  {
	url, _ := url.Parse("http://127.0.0.1:8888/")
	chainId := make([]byte, 32)
	return eos.New(url, chainId)
}

//func TestDeployContract(t *testing.T) {
//	contractName, err := createAccount("zhangsan")
//	if err != nil {
//		fmt.Println(err)
//		t.Fail()
//	}
//	//contractDir := "../../contract/bounty"
//	abiFile := "../../contract/bounty/bounty.abi"
//	wastFile := "../../contract/bounty/bounty.wast"
//	e := getEos()
//	//newProjectName := "project1"
//	//res, err := esys.NewSetCodeTx("a1", wastFile, abiFile)
//	e.SetSigner(eos.NewWalletSigner(e, "eosio"))
//	res, err := e.SetCode(contractName, wastFile, abiFile)
//	if err != nil {
//		fmt.Println(err)
//		t.Fail()
//	}
//
//	fmt.Println(res)
//	t.Log("success")
//}

//func TestCreate(t *testing.T) {
//	e := util.GetEos()
//	res,err := e.CreateWallet("daooner1241")
//	if err != nil {
//		fmt.Println("error::", err)
//	}
//
//	fmt.Println("normal::", res)
//}

//func TestKey(t *testing.T) {
//	res,err := ecc.NewRandomPrivateKey()
//	if err !=nil {fmt.Println(err)}
//	fmt.Println("Private Key", res.String())
//
//	fmt.Println(res.PublicKey().String())
//}

//func TestCreateAccount(t *testing.T) {
//	e := util.GetEos()
//	rand.Seed(time.Now().Unix())
//	accountName := "user" + strconv.Itoa(rand.Int())[:4]
//	//accountName = "fucker1"
//	fmt.Println(accountName)
//	fmt.Println(e.CreateAccount(accountName))
//}

//func TestCreateTask(t *testing.T) {
//	e := util.GetEos()
//	fmt.Println(e.CreateTask("a1", "daoone", "daoone", "dddd", 100))
//}

//func TestGetBalance(t *testing.T) {
//	e := util.GetEos()
//	fmt.Println(e.GetCurrencyBalance("daoyqupzcrz2", "DOT", "daoone.token"))
//}
