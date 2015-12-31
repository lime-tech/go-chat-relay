package irc

import (
	"container/ring"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/goirc/state"
	"go-chat-relay/src/globals"
	"go-chat-relay/src/helpers"
	"go-chat-relay/src/plug"
	"go-chat-relay/src/plug/config"
	"go-chat-relay/src/plugs"
	"runtime"
	"time"
)

const (
	PING_FREQ       = 30 * time.Second
	SPLIT_LEN       = 255
	TIMEOUT         = 60 * time.Second
	CHANGE_CAPACITY = 64
)

func init() {
	plugs.Register("irc", New)
}

func nextNick(nicks []string) func(string) string {
	r := ring.New(len(nicks))
	for i := 0; i < len(nicks); i++ {
		r.Value = nicks[i]
		r = r.Next()
	}
	return func(s string) string {
		// start from second nick, first already used in state.Nick.Nick
		r = r.Next()
		return r.Value.(string)
	}
}

func logPanic(conn *client.Conn, line *client.Line) {
	if err := recover(); err != nil {
		_, f, l, _ := runtime.Caller(2)
		log.Error("%s:%d: panic: %+v", f, l, err)
	}
}

type IRCPlug struct {
	config  config.Config
	client  *client.Conn
	endloop chan bool
	changes chan plug.Change
}

func (p *IRCPlug) Connect() error              { return p.client.Connect() }
func (p *IRCPlug) Loop()                       { <-p.endloop }
func (p *IRCPlug) Changes() <-chan plug.Change { return p.changes }
func (p *IRCPlug) Send(msg string) error {
	for _, c := range p.config.Channels {
		p.client.Privmsg(c, msg)
	}
	return nil
}
func (p *IRCPlug) onDisconnected(c *client.Conn, l *client.Line) { p.endloop <- true }
func (p *IRCPlug) onConnected(c *client.Conn, l *client.Line) {
	for _, ch := range p.config.Channels {
		c.Join(ch)
	}
}
func (p *IRCPlug) onMessage(c *client.Conn, l *client.Line) {
	channel, msg := l.Args[0], l.Args[1]
	if !helpers.ContainsString(p.config.Channels, channel) {
		// ignore non channel messages, etc
		return
	}
	if helpers.ContainsString(p.config.Identity.Nick, l.Nick) {
		// do not relay messages from me
		return
	}
	p.changes <- plug.Change{
		User:    l.Nick,
		Channel: channel,
		Server:  p.config.Connection.Addr,
		Data:    msg,
	}
}
func (p *IRCPlug) onTopic(c *client.Conn, l *client.Line) {
	channel, msg := l.Args[0], l.Args[1]
	if !helpers.ContainsString(p.config.Channels, channel) {
		// ignore non channel messages, etc
		return
	}
	p.changes <- plug.Change{
		User:    l.Nick,
		Channel: channel,
		Server:  p.config.Connection.Addr,
		Data:    fmt.Sprintf("Topic changed: %s", msg),
	}
}

func New(pc config.Config) plug.Plug {
	st := &state.Nick{
		Nick:  pc.Identity.Nick[0],
		Ident: pc.Identity.RealName,
		Name:  globals.Name(),
	}
	cf := &client.Config{
		Server:   pc.Connection.Addr,
		Pass:     pc.Identity.Password,
		Me:       st,
		NewNick:  nextNick(pc.Identity.Nick),
		Version:  globals.Version(),
		PingFreq: PING_FREQ,
		SplitLen: SPLIT_LEN,
		Timeout:  TIMEOUT,
		Recover:  logPanic,
		Flood:    true,
	}
	cl := client.Client(cf)
	p := &IRCPlug{
		config:  pc,
		client:  cl,
		endloop: make(chan bool),
		changes: make(chan plug.Change, CHANGE_CAPACITY),
	}
	cl.HandleFunc(client.CONNECTED, p.onConnected)
	cl.HandleFunc(client.DISCONNECTED, p.onDisconnected)
	cl.HandleFunc(client.PRIVMSG, p.onMessage)
	cl.HandleFunc(client.TOPIC, p.onTopic)
	return p
}
