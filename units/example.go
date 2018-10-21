package units

import (
	"github.com/Houjingchao/gopressure/data"
	"github.com/Houjingchao/gopressure/interfaces"
	"github.com/Houjingchao/gopressure/manager"
	"github.com/cocotyty/httpclient"
)

const HOST = "http://baidu.com"

type ExampleContext struct {
	accountId int
	token     string
	BaseContext
}

type TxExample struct {
	result data.ReqResult
}

func (tx *TxExample) Request(ctx interfaces.Context, data interface{}) (err error) {
	resp := httpclient.Get("https://baidu.com").Send()
	tx.result.Code,_=resp.Code()
	tx.result.Success = tx.result.Code == 200
	return
}

func (tx *TxExample) Report() (res data.ReqResult) {
	return tx.result
}

type testLink struct {
	txlist [][]interfaces.Tx
	index  int
	ctx    interfaces.Context
}

func newLink() (l *testLink, err error) {
	l = &testLink{index: 0}
	l.txlist = make([][]interfaces.Tx, 3)

	var (
		txExample = &TxExample{
			data.ReqResult{},
		}
	)
	l.txlist[0] = []interfaces.Tx{txExample}
	//tx3 tx4 并发执行
	/*
	l.txlist[2] = []interfaces.Tx{tx3, tx4}
	*/
	return l, nil
}

func (link *testLink) Setup(ctx interfaces.Context) (err error) {
	link.ctx = ctx
	return
}

func (link *testLink) Next(ctx interfaces.Context) (txs []interfaces.Tx, err error) {
	if link.index > len(link.txlist)-1 {
		return
	}

	txs = link.txlist[link.index]
	link.index++
	return

}

func (link *testLink) Report() (res data.ReqResult) {
	//判断此测试链路是否成功
	if false {
		return data.ReqResult{
			Name:      "just a test link",
			Code:      404,
			Success:   false,
			Extradata: "some error info or something need to report",
		}
	}
	return data.ReqResult{
		Name:      "link",
		Code:      200,
		Success:   true,
		Extradata: "something need to report",
	}
}

type ExampleFactory struct {
	dataProvider interfaces.Provider
}

func (factory *ExampleFactory) GetDataProvider() (p interfaces.Provider, err error) {
	return factory.dataProvider, nil
}

func (factory *ExampleFactory) GetLink() (link interfaces.Link, err error) {
	return newLink()

}

func (factory *ExampleFactory) GetContext() (context interfaces.Context, err error) {
	return &ExampleContext{}, nil

}
func (factory *ExampleFactory) InitFactory(args string) (err error) {
	return nil
}

func prepareProviderData() []interface{} {
	return nil
}

func init() {
	provider := &BaseDataProvider{}
	provider.InitPool(prepareProviderData())
	manager.RegisterFactory("example", &ExampleFactory{provider})
}
