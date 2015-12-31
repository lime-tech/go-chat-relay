package plug

import (
	"fmt"
)

type Plug interface {
	Connect() error
	Loop()
	Changes() <-chan Change
	Send(string) error
}

type Change struct {
	User    string
	Channel string
	Server  string
	Data    string
}

func (c Change) String() string {
	return fmt.Sprintf(
		"%s%s@%s> %s",
		c.User,
		c.Channel,
		c.Server,
		c.Data,
	)
}
