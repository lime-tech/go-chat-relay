package cli

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"go-chat-relay/src/config"
	"go-chat-relay/src/relay"
)

func startAction(c *cli.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()
	config.FromFile(c.String("config"))
	go relay.Run()
	select {}
}
