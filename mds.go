package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/hashicorp/go-version"
	"github.com/urfave/cli/v2"

	"github.com/shimingyah/mds/cmd"
	"github.com/shimingyah/mds/utils"
)

const (
	mdsGoVersion        = "1.16.0"
	goVersionConstraint = ">= " + mdsGoVersion
)

func flags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"debug", "v"},
			Usage:   "enable debug log",
		},
		&cli.BoolFlag{
			Name:    "quiet",
			Aliases: []string{"q"},
			Usage:   "only warning and errors",
		},
		&cli.BoolFlag{
			Name:  "trace",
			Usage: "enable trace log",
		},
	}
}

// Check if this binary is compiled with at least minimum Go version.
func checkGoVersion(goVersionStr string) error {
	constraint, err := version.NewConstraint(goVersionConstraint)
	if err != nil {
		return fmt.Errorf("'%s': %s", goVersionConstraint, err)
	}

	goVersion, err := version.NewVersion(goVersionStr)
	if err != nil {
		return err
	}

	if !constraint.Check(goVersion) {
		return fmt.Errorf("MDS is not compiled by go %s. Minimum required version is %s", goVersionStr, mdsGoVersion)
	}

	return nil
}

func main() {
	if err := checkGoVersion(runtime.Version()[2:]); err != nil {
		log.Fatalln(err)
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name: "version", Aliases: []string{"V"},
		Usage: "print only the version",
	}

	app := &cli.App{
		Name:      "mds",
		Usage:     "A high-performance metadata server for filesystem.",
		Version:   utils.Version(),
		Copyright: "Apache-2.0",
		Flags:     flags(),
		Commands: []*cli.Command{
			cmd.Server(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
