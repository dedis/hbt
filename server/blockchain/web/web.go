package web

import (
	"os"
	"time"

	"go.dedis.ch/dela"
	"go.dedis.ch/dela/cli"
	"go.dedis.ch/dela/cli/node"
	"go.dedis.ch/dela/mino/proxy"
	"go.dedis.ch/dela/mino/proxy/http"
	"golang.org/x/xerrors"
)

var defaultRetry = 10
var proxyFac func(string) proxy.Proxy = http.NewHTTP

const defaultProxyAddr = "127.0.0.1:0"

// NewController returns a new controller initializer
func NewController() node.Initializer {
	return controller{}
}

// controller is an initializer with a set of commands.
//
// - implements node.Initializer
type controller struct{}

// Build implements node.Initializer.
func (m controller) SetCommands(builder node.Builder) {
	builder.SetStartFlags(
		cli.StringFlag{
			Name:     "proxyaddr",
			Usage:    "the proxy address",
			Required: false,
			Value:    defaultProxyAddr,
		},
	)
}

// OnStart implements node.Initializer. It creates and registers a pedersen DKG.
func (m controller) OnStart(ctx cli.Flags, inj node.Injector) error {
	dela.Logger.Info().Msg("Installing Blockchain proxy")

	proxyAddr := ctx.String("proxyaddr")

	proxyhttp := proxyFac(proxyAddr)

	inj.Inject(proxyhttp)

	go proxyhttp.Listen()

	for i := 0; i < defaultRetry && proxyhttp.GetAddr() == nil; i++ {
		time.Sleep(time.Second)
	}

	if proxyhttp.GetAddr() == nil {
		return xerrors.Errorf("failed to start proxy server")
	}

	// We assume the listen worked proprely, however it might not be the case.
	// The log should inform the user about that.
	dela.Logger.Info().Msgf("started proxy server on %s", proxyhttp.GetAddr().String())

	//
	// Register the smc proxy handlers
	//

	register := RegisterAction{}
	err := register.Execute(node.Context{
		Injector: inj,
		Flags:    node.FlagSet{},
		Out:      os.Stdout,
	})

	if err != nil {
		return xerrors.Errorf("failed to register evoting handlers: %v", err)
	}

	return nil
}

// OnStop implements node.Initializer.
func (controller) OnStop(_ node.Injector) error {
	return nil
}
