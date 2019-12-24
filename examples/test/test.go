package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/seerx/chain/pkg/context"

	"github.com/seerx/chain/pkg/intf"

	"github.com/seerx/chain"
)

type ApiTest struct {
	AN string    `json:"an"`
	V  int       `json:"v"`
	D  float32   `json:"d"`
	T  time.Time `json:"t"`
}

func (a ApiTest) Group() *intf.Group {
	return &GA
}

type N struct {
	II bool
}

type U struct {
	ID   int    `json:"id" c:"desc:ID,require"`
	Name string `json:"name" c:"desc:名称"`

	US []*U `json:"us" c:"desc:yes"`
	Ni N    `json:"ni"`
}

type Response struct {
	Val   string   `json:"val" c:"desc:123,deprecated"`
	Key   string   `json:"key" c:"desc:键,require"`
	Items []string `json:"items" c:"desc:数组"`
	UAry  []*U     `json:"uAry" c:"desc:U数组"`
}

type ReqID struct {
	ID int `json:"id" c:"desc="ID 值"`
}

type Req struct {
	A    string   `json:"a,omitempty" c:"desc:测试A,require,limit:10<$v"`
	B    *string  `json:"b" c:"desc:测试B ptr"`
	Req  ReqID    `json:"req" c:"desc:测试结构"`
	Reqs []*ReqID `json:"reqs" c:"desc:啊哈"`
}

func (a *ApiTest) Test1(aa *Req) ([]*Response, error) {
	return []*Response{&Response{Val: "123"}}, nil
}

func (a ApiTest) Test2(abb string) (*Response, error) {
	return &Response{
		Val:   "123",
		Key:   abb,
		Items: nil,
		UAry:  nil,
	}, nil
}

func InjectFn(a map[string]interface{}) (io.Closer, error) {
	return nil, nil
}

func InjectFn1(a map[string]interface{}) (io.ReadCloser, error) {
	return nil, nil
}

func main() {
	ch := chain.New()

	ch.Register(&ApiTest{})
	err := ch.Explain()
	if err != nil {
		panic(err)
	}

	B := "aaaaaaa"

	//req := &Req{
	//	A:    "123",
	//	B:    &B,
	//	Req:  ReqID{ID: 100},
	//	Reqs: []*ReqID{{ID: 11}, {ID: 12}},
	//}

	reqs := chain.Requests{}
	reqs = append(reqs, &chain.Request{
		Service: "test.Test2",
		Alias:   "ABC",
		Args:    B,
	})
	//reqs = append(reqs, &chain.Request{
	//	Service: "test.Test1",
	//	Args:    req,
	//})

	data, _ := json.Marshal(reqs)
	str := string(data)

	rsp, err := ch.Execute(&context.Context{}, str)
	if err != nil {
		panic(err)
	}

	data, _ = json.Marshal(rsp)
	str = string(data)
	fmt.Println(str)
}
