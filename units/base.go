package units

import (
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
)

const DEFAULT_HTTP_TIMEOUT = 30

type BaseContext struct {
	cli *http.Client
}

func (c *BaseContext) GetHttpClient() *http.Client {
	if c.cli == nil {
		jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		c.cli = &http.Client{Jar: jar}
	}
	return c.cli
}

type BaseDataProvider struct {
	DataPool chan interface{}
	capacity int
}

func (p *BaseDataProvider) Get() interface{} {
	if p.capacity <= 0 {
		return nil
	}
	return <-p.DataPool
}

func (p *BaseDataProvider) Put(data interface{}) {
	if p.capacity <= 0 {
		return
	}
	p.DataPool <- data
	return
}
func (p *BaseDataProvider) GetCapacity() int {
	return p.capacity
}

func (p *BaseDataProvider) InitPool(datas []interface{}) {
	p.capacity = len(datas)
	p.DataPool = make(chan interface{}, p.capacity)
	for _, v := range datas {
		p.DataPool <- v
	}
	return
}
