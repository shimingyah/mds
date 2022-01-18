package cmd

import (
	"fmt"
	"runtime"

	"github.com/shimingyah/mds/server"
	"github.com/urfave/cli/v2"
)

// Server run mds server
func Server() *cli.Command {
	cmd := &cli.Command{
		Name:      "server",
		Usage:     "start a mds",
		ArgsUsage: "mds server -port=8000",
		Action:    serverAction,
	}
	cmd.Flags = append(cmd.Flags, serverFlags()...)
	cmd.Flags = append(cmd.Flags, commonFlags()...)
	return cmd
}

func serverFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "ip",
			Value: "0.0.0.0",
			Usage: "server ip address, default is 0.0.0.0",
		},
		&cli.IntFlag{
			Name:  "port",
			Value: 8000,
			Usage: "http or grpc listen port",
		},
		&cli.IntFlag{
			Name:  "maxcpu",
			Value: 0,
			Usage: "maximum number of CPUs. 0 means all available CPUs",
		},
		&cli.IntFlag{
			Name:  "read-timeout",
			Value: 60,
			Usage: "read timeout in seconds",
		},
		&cli.IntFlag{
			Name:  "write-timeout",
			Value: 120,
			Usage: "write timeout in seconds",
		},
	}
}

func serverAction(c *cli.Context) error {
	initLog(c)

	logger.Infof("%s module", c.Command.Name)

	maxcpu := c.Int("maxcpu")
	if maxcpu < 1 {
		maxcpu = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(maxcpu)

	endpoint := fmt.Sprintf("%s:%d", c.String("ip"), c.Int("port"))

	mds := server.NewMDS(endpoint, c.Int("read-timeout"), c.Int("write-timeout"))
	mds.Run()

	return nil
}
