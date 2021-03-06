package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/seerx/runjson/pkg/context"

	"github.com/seerx/runjson"

	"github.com/seerx/runjson/examples/test/tt"

	"github.com/seerx/runjson/pkg/graph"

	"github.com/seerx/runjson/pkg/rj"
)

// GA 分组
var GA = rj.Group{
	Name:        "test",
	Description: "测试组",
}

type ApiTest struct {
	AN string    `json:"an"`
	V  int       `json:"v"`
	D  float32   `json:"d"`
	T  time.Time `json:"t"`

	R   rj.Require
	I1  Cls
	I2  *Cls
	Rsp rj.ResponseContext
}

func (a ApiTest) Group() *rj.Group {
	return &GA
}

type N struct {
	II bool
}

type U struct {
	ID   int    `json:"id" rj:"desc:ID,require"`
	Name string `json:"name" rj:"desc:名称"`

	US []*U `json:"us" rj:"desc:yes,deprecated"`
	Ni N    `json:"ni"`
}

type Response struct {
	Val   string   `json:"val" rj:"desc:123,deprecated"`
	Key   string   `json:"key" rj:"desc:键,require"`
	Items []string `json:"items" rj:"desc:数组"`
	UAry  []*U     `json:"uAry" rj:"desc:U数组"`
	//Error error    `json:"error"`
}

type ReqID struct {
	ID int `json:"id" rj:"desc:123"`
}

type Req struct {
	A    string   `json:"a,omitempty" rj:"desc:测试A,require,limit:10<$v"`
	B    *string  `json:"b" rj:"desc:测试B ptr"`
	N    uint64   `json:"n,string" rj:"range:(0, 100)"`
	BA   float64  `json:"ba" rj:"range:[10.10,1200]"`
	Req  ReqID    `json:"req" rj:"desc:测试结构"`
	Reqs []*ReqID `json:"reqs" rj:"desc:啊哈"`
}

func (a *ApiTest) Test1Info() *rj.FuncInfo {
	return &rj.FuncInfo{
		Description: "测试函数1",
		Deprecated:  false,
		//InputIsRequire: true,
		History: nil,
	}
}

func (a ApiTest) Test1(aa Req, cls Cls) ([]*Response, error) {
	fmt.Println("111====", a.R, aa.N)
	//c.Close()
	cls.Close()
	a.I1.Close()
	a.I2.Close()
	func() {
		val, err := a.Rsp.Get((&ApiTest{}).Test2)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(val)
		}
	}()

	val, err := a.Rsp.Get((&tt.TT{}).New)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(val)
	}

	return []*Response{&Response{Val: "123"}}, nil
}

func (a *ApiTest) Test2(abb string) (*Response, error) {
	return &Response{
		Val:   "123",
		Key:   abb,
		Items: nil,
		UAry:  nil,
		//Error: errors.New("Error ............"),
	}, nil
}

func (a *ApiTest) Test3() (*Response, error) {
	return &Response{
		Val:   "123",
		Key:   "abb",
		Items: nil,
		UAry:  nil,
		//Error: errors.New("Error ............"),
	}, nil
}

func (a *ApiTest) Test2Info() rj.FuncInfo {
	//val, err := a.Rsp.Get((&ApiTest{}).Test1)
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println(val)
	//}

	return rj.FuncInfo{
		Description: `测试函数 222`,
		History: []*graph.CR{
			{"2019/12/26", "hyb", "创建"},
		},
	}
}

type Cls struct {
	V string
}

func (c *Cls) Close() error {
	fmt.Println("OK", c.V)
	return nil
}

func InjectFn(arg *rj.InjectArg) (io.Closer, error) {
	return &Cls{V: "test"}, nil
}

func InjectFn1(arg *rj.InjectArg) (*Cls, error) {
	return &Cls{V: "另一个"}, nil
}

func main() {

	//mp := reflects.Parse("limit:0-$v--10,desc:测试")
	//data, _ := json.MarshalIndent(mp, "", "  ")
	//fmt.Println(string(data))

	ch := runjson.New()
	if err := ch.RegisterProvider(InjectFn); err != nil {
		panic(err)
	}
	if err := ch.RegisterAccessController(InjectFn1); err != nil {
		panic(err)
	}
	ch.Register(&ApiTest{}, &tt.TT{})

	//ch.BeforeExecute(func(item *rj.Request) {
	//	fmt.Println("before:", item.Service)
	//}).AfterExecute(func(item *rj.Request, result *rj.ResponseItem, results rj.ResponseContext) {
	//	fmt.Println("after:", item.Service)
	//})

	err := ch.Engage()

	//if info, err := json.MarshalIndent(ch.ApiInfo, "", "\t"); err == nil {
	//	data := string(info)
	//	fmt.Println(data)
	//}

	if err != nil {
		panic(err)
	}

	B := "aaaaaaa"

	req := &Req{
		A:   "1230099387747474y44 d",
		B:   &B,
		N:   9,
		BA:  100,
		Req: ReqID{ID: 100},
		//Reqs: []*ReqID{{ID: 11}, {ID: 12}},
		Reqs: nil,
	}

	reqs := rj.Requests{}
	//reqs = append(reqs, &rj.Request{
	//	Service: "test.Test2",
	//	Args:    B,
	//})
	reqs = append(reqs, &rj.Request{
		Service: "test.Test1",
		Args:    req,
	})
	//reqs = append(reqs, &rj.Request{
	//	Service: "test.Test3",
	//	Args:    req,
	//})

	data, _ := json.Marshal(reqs)
	str := string(data)
	fmt.Println(str)
	rsp, err := ch.RunString(&context.Context{}, str)
	if err != nil {
		panic(err)
	}

	data, _ = json.Marshal(rsp)
	str = string(data)
	fmt.Println(str)

	fmt.Println("DONE ..............")
}
