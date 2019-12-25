package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/seerx/runjson"

	"github.com/seerx/runjson/pkg/context"

	"github.com/seerx/runjson/pkg/intf"
)

type ApiTest struct {
	AN string    `json:"an"`
	V  int       `json:"v"`
	D  float32   `json:"d"`
	T  time.Time `json:"t"`

	R  intf.Require
	I1 Cls
	I2 *Cls
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

	US []*U `json:"us" c:"desc:yes,deprecated"`
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

func (a *ApiTest) Test1Info() string {
	return `测试函数 1`
}

func (a *ApiTest) Test1(aa Req, cls Cls) ([]*Response, error) {
	fmt.Println(a.R)
	//c.Close()
	cls.Close()
	a.I1.Close()
	a.I2.Close()
	return []*Response{&Response{Val: "123"}}, nil
}

func (a *ApiTest) Test2Info() intf.FuncInfo {
	return intf.FuncInfo{
		Descrition: `测试函数 222`,
		Deprecated: true,
	}

}

func (a ApiTest) Test2(abb string) (*Response, error) {
	return &Response{
		Val:   "123",
		Key:   abb,
		Items: nil,
		UAry:  nil,
	}, nil
}

type Cls struct {
	V string
}

func (c *Cls) Close() error {
	fmt.Println("OK", c.V)
	return nil
}

func InjectFn(a map[string]interface{}) (io.Closer, error) {
	return &Cls{V: "test"}, nil
}

func InjectFn1(a map[string]interface{}) (*Cls, error) {
	return &Cls{V: "另一个"}, nil
}

func main() {
	ch := runjson.New()
	if err := ch.Inject(InjectFn, InjectFn1); err != nil {
		panic(err)
	}
	ch.Register(&ApiTest{})
	err := ch.Explain()

	if info, err := json.MarshalIndent(ch.ApiInfo, "", "\t"); err == nil {
		data := string(info)
		fmt.Println(data)
	}

	if err != nil {
		panic(err)
	}

	B := "aaaaaaa"

	req := &Req{
		A:    "1230099387747474y44 d",
		B:    &B,
		Req:  ReqID{ID: 100},
		Reqs: []*ReqID{{ID: 11}, {ID: 12}},
	}

	reqs := runjson.Requests{}
	reqs = append(reqs, &runjson.Request{
		Service: "test.Test1",
		Alias:   "ABC",
		Args:    req,
	})
	//reqs = append(reqs, &runjson.Request{
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
