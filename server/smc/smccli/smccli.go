package main

import (
	"fmt"
	"io"
	"os"

	"go.dedis.ch/dela/cli/node"
	dkg "go.dedis.ch/dela/dkg/pedersen/controller"
	minogrpc "go.dedis.ch/dela/mino/minogrpc/controller"
	proxy "go.dedis.ch/dela/mino/proxy/http/controller"
	smc "go.dedis.ch/hbt/server/smc/smccli/controller"
	"go.dedis.ch/hbt/server/smc/smccli/web"
	kv "go.dedis.ch/purb-db/store/kv/controller"
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
		web.NewController(),
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
