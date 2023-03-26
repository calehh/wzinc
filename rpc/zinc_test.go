package rpc

import (
	"fmt"
	"testing"
	client "github.com/zinclabs/sdk-go-zincsearch"
)

const zincUrl = ""
const port = "1234"
const username = ""
const password = ""

func init() {
	InitRpcService(zincUrl, port, username, password)
}

func TestInput(t *testing.T) {
	res, err := RpcServer.zincQuery(QueryReq{
		SearchType: "",
		Query:      Query{
			Term: "",
		},
		From:       0,
		MaxResult:  0,
	}, "test")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(res))
}
