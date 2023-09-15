package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	portFlag = &cli.IntFlag{
		Name:  "port",
		Value: 8888,
		Usage: "Port setting",
	}
	nodeNameFlag = &cli.StringFlag{
		Name:  "node",
		Value: "server-0",
		Usage: "Server",
	}
	targetFlag = &cli.StringFlag{
		Name:  "target",
		Value: "",
		Usage: "Target",
	}
)

var app = cli.NewApp()

func init() {
	app.Usage = "the go-ethereum devp2p hand-on command line interface"
	app.Copyright = "Copyright 2023 The An Quach"
	app.Flags = []cli.Flag{
		portFlag,
		nodeNameFlag,
		targetFlag,
	}
	app.Action = node
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func node(ctx *cli.Context) error {
	srv, err := newServer(ctx.String(nodeNameFlag.Name), ctx.Int(portFlag.Name))
	if err != nil {
		log.Printf("Failed to create new server with err: %v", err)
	}

	log.Println(srv.NodeInfo().Enode)

	if ctx.String(targetFlag.Name) != "" {
		if err := connectPeer(srv, ctx.String(targetFlag.Name)); err != nil {
			log.Printf("Failed to connect to peer with err: %v", err)
		}
	}

	select {}
}
