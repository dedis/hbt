package main

import (
	"fmt"
	"io"
	"os"

	"go.dedis.ch/dela/cli/node"
	kv "go.dedis.ch/dela/core/store/kv/controller"
	dkg "go.dedis.ch/dela/dkg/pedersen/controller"
	minogrpc "go.dedis.ch/dela/mino/minogrpc/controller"
	proxy "go.dedis.ch/dela/mino/proxy/http/controller"
	smc "go.dedis.ch/hbt/server/smc/smccli/controller"
)

type config struct {
	Channel chan os.Signal
	Writer  io.Writer
}

func main() {
	err := run(os.Args)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func run(args []string) error {
	return runWithCfg(args, config{Writer: os.Stdout})
}

func runWithCfg(args []string, cfg config) error {
	builder := node.NewBuilderWithCfg(
		cfg.Channel,
		cfg.Writer,
		kv.NewController(),
		proxy.NewController(),
		proxy.NewController(),
		minogrpc.NewController(),
		dkg.NewMinimal(),
		smc.NewSmcController(),
	)

	app := builder.Build()

	err := app.Run(args)
	if err != nil {
		return err
	}

	return nil
}
