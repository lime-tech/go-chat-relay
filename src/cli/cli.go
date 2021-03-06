package cli

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"go-chat-relay/src/globals"
	"os"
	"os/signal"
	"syscall"
)

var (
	logLevel log.Level
)

func getLogLevel(c *cli.Context) error {
	log.SetOutput(os.Stderr)
	level, err := log.ParseLevel(c.String("log-level"))
	if err != nil {
		return err
	}

	log.SetLevel(level)
	logLevel = level
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	return nil
}

func toggleDebugLevel() {
	currLevel := log.GetLevel()
	if currLevel == log.DebugLevel {
		log.Info("Disabling debug log level")
		log.SetLevel(logLevel)
	} else {
		log.Info("Enabling debug log level")
		log.SetLevel(log.DebugLevel)
	}
}

func signalingLoop(sigChan chan os.Signal) {
MainLoop:
	for {
		sig := <-sigChan
		switch sig {
		case syscall.SIGUSR1:
			toggleDebugLevel()
		case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
			defer os.Exit(0)
			log.Warn("Exiting with 0 code")
			close(sigChan)
			break MainLoop
		case syscall.SIGHUP:
			log.Warnf("Ignoring SIGHUP")
		default:
			log.Warnf("Got %+v, but don't know how to handle it", sig)
		}
	}
}

func Run() {
	app := cli.NewApp()
	app.Name = globals.Name()
	app.Usage = "Extendable chat messages relay"
	app.Version = fmt.Sprintf("%s", globals.Version())
	app.Authors = []cli.Author{
		cli.Author{Name: "dubi", Email: "devnull@localhost"},
	}
	app.Flags = rootFlags
	app.Before = getLogLevel
	app.Commands = commands

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan)
	go signalingLoop(sigChan)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Got an error initial run %+v", err)
	}
}
