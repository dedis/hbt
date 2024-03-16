package main

import (
	"fmt"
	"io"
	"os"

	calypso "go.dedis.ch/hbt/server/blockchain/calypso/controller"

	"go.dedis.ch/dela/cli/node"
	access "go.dedis.ch/dela/contracts/access/controller"
	cosipbft "go.dedis.ch/dela/core/ordering/cosipbft/controller"
	kv "go.dedis.ch/dela/core/store/kv/controller"
	pool "go.dedis.ch/dela/core/txn/pool/controller"
	signed "go.dedis.ch/dela/core/txn/signed/controller"
	minogrpc "go.dedis.ch/dela/mino/minogrpc/controller"
	proxy "go.dedis.ch/dela/mino/proxy/http/controller"
	"go.dedis.ch/hbt/server/blockchain/web"
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
		kv.NewController(),
		proxy.NewController(),
		web.NewController(),
		minogrpc.NewController(),
		cosipbft.NewController(),
		signed.NewManagerController(),
		pool.NewController(),
		access.NewController(),
		calypso.NewController(),
	)

	app := builder.Build()

	err := app.Run(args)
	if err != nil {
		return err
	}

	return nil
}
