package cmd

import (
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/urfave/cli/v2"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/shimingyah/mds/logrus"
)

var logger = logrus.GetLogger("mds")

func initLog(c *cli.Context) {
	setLogLevel(c)
	setLogOutput(c)
}

func setLogLevel(c *cli.Context) {
	if c.Bool("trace") {
		logrus.SetLevel(logger, "trace")
	} else if c.Bool("debug") {
		logrus.SetLevel(logger, "debug")
	} else if c.Bool("warn") {
		logrus.SetLevel(logger, "warn")
	}
	logger.SetReportCaller(true)
}

func setLogOutput(c *cli.Context) {
	fileName := filepath.Join(c.String("log-dir"), c.String("log-name"))
	logger.SetOutput(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    c.Int("log-max-size"), // MB
		MaxBackups: c.Int("log-max-file"),
		MaxAge:     c.Int("log-max-age"), //days
		LocalTime:  true,
	})
}

func commonFlags() []cli.Flag {
	var defaultLogDir = "/tmp/logs"
	var defaultLogFileName = "mds.log"
	var defaultHostName = "localhost"
	switch runtime.GOOS {
	case "darwin":
		fallthrough
	case "windows":
		homeDir, err := os.UserHomeDir()
		if err == nil {
			defaultLogDir = path.Join(homeDir, ".mds", "log")
		}
		hostname, err := os.Hostname()
		if err == nil {
			defaultHostName = hostname
		}
	}
	return []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "for debug, umount old volume at mountPoint/debug log",
		},
		&cli.StringFlag{
			Name:  "conf",
			Usage: "conf file path",
		},
		&cli.StringFlag{
			Name:  "log-dir",
			Value: defaultLogDir,
			Usage: "log path",
		},
		&cli.StringFlag{
			Name:  "log-name",
			Value: defaultLogFileName,
			Usage: "log file name",
		},
		&cli.IntFlag{
			Name:  "log-max-size",
			Value: 1024, // MiB
			Usage: "log file max size",
		},
		&cli.IntFlag{
			Name:  "log-max-file",
			Value: 10,
			Usage: "log max file num",
		},
		&cli.IntFlag{
			Name:  "log-max-age",
			Value: 7, // day
			Usage: "log save time",
		},
		&cli.StringFlag{
			Name:  "host",
			Value: defaultHostName,
			Usage: "hostname / odin metric host",
		},
	}
}
