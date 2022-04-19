package example

import (
	"das_register_server/http_server/handle"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/sign"
	"github.com/scorpiotzh/toolib"
	"testing"
)

func TestAddressDeposit(t *testing.T) {
	var req handle.ReqAddressDeposit
	req.AlgorithmId = 3
	req.Address = "0x15a33588908cf8edb27d1abe3852bf287abd3891"

	url := TestUrl + "/address/deposit"
	var data handle.RespAddressDeposit
	if err := doReq(url, req, &data); err != nil {
		t.Fatal(err)
	}
	fmt.Println(toolib.JsonString(data))
}

func TestBalanceDeposit(t *testing.T) {
	var req handle.ReqBalanceDeposit
	req.FromCkbAddress = "ckt1qyqdpc0e57dtjdscy89u8vclenhqjn904ngqrq78am"
	req.ToCkbAddress = "ckt1qqexmutxu0c2jq9q4msy8cc6fh4q7q02xvr7dc347zw3ks3qka0m6qg9zk3ntzys3nuwmvnar2lrs54l9pat6wy3q526xdvgjzx03mdj05dtuwzjhu5840fcjy2c9u8d"
	req.Amount = 116 * common.OneCkb

	url := TestUrl + "/balance/deposit"

	var data handle.RespBalanceDeposit
	if err := doReq(url, req, &data); err != nil {
		t.Fatal(err)
	}
	fmt.Println(toolib.JsonString(data))

	private := ""
	for i, v := range data.SignList {
		signData, err := sign.EthSignature(common.Hex2Bytes(v.SignMsg), private)
		if err != nil {
			t.Fatal(err)
		}
		data.SignList[i].SignMsg = common.Bytes2Hex(signData)
	}
	fmt.Println(toolib.JsonString(data))
}
