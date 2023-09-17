package utils

import (
	"flag"
	"os"

	"github.com/ethereum/go-ethereum/log"
)

const (
	BzzDefaultNetworkId = 4242
	WSDefaultPort       = 18543
	BzzDefaultPort      = 8542
)

var (
	// custom log, easily grep'able
	Log = log.New("masterlogs", "*")

	// our working directory
	BasePath string

	// out local port for p2p connections
	P2PPort int

	// self-explanatory command line arguments
	verbose = flag.Bool("v", false, "more verbose logs")
	// remoteport   = flag.Int("p", 0, "remote port (enables remote RPC lookup of enode)")
	// remotehost   = flag.String("h", "127.0.0.1", "remote host (RPC, p2p)")
	// enodeconnect = flag.String("e", "", "enode to connect to (overrides remote RPC lookup)")
	// p2plocalport = flag.Int("l", P2pPort, "local port for p2p connections")
)

// setup logging
// set up remote node, if present
func init() {
	var err error

	flag.Parse()

	// get the working directory
	BasePath, err = os.Getwd()
	if err != nil {
		Log.Crit("Could not determine working directory", "err", err)
	}

	// ensure good log formats for terminal
	// handle verbosity flag
	hs := log.StreamHandler(os.Stderr, log.TerminalFormat(true))
	loglevel := log.LvlInfo
	if *verbose {
		loglevel = log.LvlTrace
	}
	hf := log.LvlFilterHandler(loglevel, hs)
	h := log.CallerFileHandler(hf)
	log.Root().SetHandler(h)
}

//func
