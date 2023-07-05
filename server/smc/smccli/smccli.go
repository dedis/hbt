package main

import (
	"fmt"
	"go.dedis.ch/dela/cli/node"
	db "go.dedis.ch/dela/core/store/kv/controller"
	dkg "go.dedis.ch/dela/dkg/pedersen/controller"
	minogrpc "go.dedis.ch/dela/mino/minogrpc/controller"
	smcctl "go.dedis.ch/hbt/smc/smccli/controller"
	"io"
	"os"
)

func main() {
	err := run(os.Args)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func run(args []string) error {
	return runWithCfg(args, config{Writer: os.Stdout})
}

type config struct {
	Channel chan os.Signal
	Writer  io.Writer
}

func runWithCfg(args []string, cfg config) error {
	builder := node.NewBuilderWithCfg(
		cfg.Channel,
		cfg.Writer,
		db.NewController(),
		minogrpc.NewController(),
		dkg.NewMinimal(),
		smcctl.NewSmcController(),
	)

	app := builder.Build()

	err := app.Run(args)
	if err != nil {
		return err
	}

	return nil
}
