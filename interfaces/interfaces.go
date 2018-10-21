package interfaces

import "github.com/Houjingchao/gopressure/data"

type Context interface{}

// 压测试单元
type Tx interface {
	Request(ctx Context, d interface{}) error
	Report() data.ReqResult
}

// 压测链路
type Link interface {
	Setup(ctx Context) (err error)
	Next(ctx Context) (txs []Tx, err error)
	Report() data.ReqResult
}

type Provider interface {
	Get() interface{}
	Put(interface{})
	GetCapacity() int
}

type Factory interface {
	InitFactory(args string) error
	GetDataProvider() (p Provider, err error)
	GetLink() (link Link, err error)
	GetContext() (context Context, err error)
}
