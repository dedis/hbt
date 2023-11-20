package main

import (
	"fmt"
	"io"
	"os"

	"go.dedis.ch/dela/cli/node"

	dkg "go.dedis.ch/dela/dkg/pedersen/controller"
	minogrpc "go.dedis.ch/dela/mino/minogrpc/controller"
	smcctl "go.dedis.ch/hbt/server/smc/smccli/controller"
	db "go.dedis.ch/purb-db/store/kv/controller"
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
