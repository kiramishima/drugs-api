package main

import (
	_ "go.uber.org/automaxprocs"
	"go.uber.org/fx"
	"kiramishima/ionix/bootstrap"
)

func main() {
	fx.New(bootstrap.Module).Run()
}
