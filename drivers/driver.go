package drivers

import (
	"context"
	"fmt"
	"os"
	"plugin"
)

type Driver interface {
	Name() string
	Version() string
	Start(context.Context) error
	Stop(context.Context) error
	TxChan(chan []byte) error
	RxChan() chan []byte
}

func LoadDriver(driverPath string) (Driver, error) {
	finfo, err := os.Stat(driverPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat driver %s: %w", driverPath, err)
	}
	if finfo.IsDir() {
		return nil, fmt.Errorf("driver %s is a directory", driverPath)
	}

	code, err := plugin.Open(driverPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open driver %s: %w", driverPath, err)
	}

	constructorSymbol, err := code.Lookup("New")
	if err != nil {
		return nil, fmt.Errorf("no constructor function on driver %s: %w", driverPath, err)
	}

	constructor, ok := constructorSymbol.(func() interface{})
	if !ok {
		return nil, fmt.Errorf("wrong constructor signature on driver %s: %w", driverPath, err)
	}
	return constructor().(Driver), nil
}
