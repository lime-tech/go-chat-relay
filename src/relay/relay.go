package relay

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"go-chat-relay/src/config"
	"go-chat-relay/src/plug"
	"go-chat-relay/src/plugs"
	_ "go-chat-relay/src/plugs/irc"
	"time"
)

func runPlug(k string, p plug.Plug) {
	for {
		log.Infof("Plug '%s' is connecting...", k)
		if err := p.Connect(); err != nil {
			log.Error(err)
			<-time.After(5 * time.Second)
			continue
		}
		log.Infof("Plug '%s' connected", k)
		p.Loop()
		log.Infof("Plug '%s' unlocked after loop", k)
		<-time.After(5 * time.Second)
	}
}

func pipePlug(i plug.Plug, o plug.Plug) {
	for c := range i.Changes() {
		if err := o.Send(c.String()); err != nil {
			log.Error(err)
		}
	}
}

type activePlugs map[string]plug.Plug

func Run() {
	cf := config.Current()
	ps := activePlugs{}
	for k, v := range cf.Plug {
		factory := plugs.Get(v.Type)
		if factory == nil {
			panic(fmt.Sprintf("Plug '%s' is not registered", v.Type))
		}
		log.Debugf("Ready to run %s %+v %+v", k, factory, v)
		p := factory(v)
		go runPlug(k, p)
		ps[k] = p
	}
	for _, v := range cf.Router.Flow {
		i, o := v[0], v[1]
		go pipePlug(ps[i], ps[o])
	}
}
