package main

import (
	"os"
	"github.com/cocotyty/summer"
	"log"
	"gopkg.in/macaron.v1"
	"github.com/mougeli/beauty"
	"strconv"
	"github.com/Houjingchao/gopressure/manager"
	"github.com/Houjingchao/gopressure/benchmark"
	_ "github.com/Houjingchao/gopressure/units"
)

func init() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./conf/dev.toml"
	}
	log.Println("配置文件路径:", configPath)
	err := summer.TomlFile(configPath)
	if err != nil {
		panic(err)
	}
}

type WebConfig struct {
	Port int `sm:"#.web.port"`
}

func Test() string {
	return "test"
}

func main() {
	webConfig := &WebConfig{}
	summer.Put(webConfig)
	summer.Start()

	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Use(beauty.Renderer())

	m.Get("/gotest", Handler)
	m.Get("/test", Test)
	m.Run(webConfig.Port)
}

func Handler(ctx *macaron.Context, r beauty.Render) (res string) {
	factoryName := ctx.Query("factory")
	tittle := ctx.Query("testTittle")
	debug := ctx.Query("debug") == "1"
	occupyData := ctx.Query("occupyData") == "1"

	log.Printf("测试FactoryName:%s",factoryName)
	concurrency, err := strconv.Atoi(ctx.Query("concurrency"))
	if err != nil || concurrency <= 0 {
		log.Println(err)
		return
	}

	last, err := strconv.Atoi(ctx.Query("lastTime"))
	if err != nil || last <= 0 {
		log.Println(err)
		return
	}

	args := ctx.Query("args")
	factory, ok := manager.GetFactory(factoryName)

	if !ok {
		log.Println("invalid factory")
		return
	}

	if err := factory.InitFactory(args); err != nil {
		log.Println(err)
		return
	}
	log.Println(factory, concurrency, last)

	bench := benchmark.NewbenchMark(factory, concurrency, int64(last), debug, occupyData)
	res, err = bench.Run()
	if err != nil {
		r.Error(err.Error())
	}

	// send email
	log.Printf("测试%s结束:",tittle)
	return
}