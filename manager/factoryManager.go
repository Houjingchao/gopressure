package manager

import "github.com/Houjingchao/gopressure/interfaces"

var factoryMapping = make(map[string]interfaces.Factory)

func RegisterFactory(name string, factory interfaces.Factory) {
	factoryMapping[name] = factory
}

func GetFactory(name string) (factory interfaces.Factory, ok bool) {
	factory, ok = factoryMapping[name]
	return
}
