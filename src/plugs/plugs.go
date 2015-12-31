package plugs

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"go-chat-relay/src/plug"
	"go-chat-relay/src/plug/config"
	"strings"
)

type AlreadyRegisteredError struct{ s string }

func (e *AlreadyRegisteredError) Error() string { return e.s }

type PlugRegistry map[string]func(config.Config) plug.Plug

var (
	plugs = PlugRegistry{}
)

func Register(n string, factory func(config.Config) plug.Plug) {
	name := strings.ToLower(n)
	if _, ok := plugs[name]; ok {
		panic(AlreadyRegisteredError{
			fmt.Sprintf("Plug with name '%s' already registered", name),
		})
	}
	log.Infof("Registering plug '%s'", name)
	plugs[name] = factory
}

func Unregister(n string) {
	name := strings.ToLower(n)
	log.Infof("Unregistering plug '%s'", name)
	delete(plugs, name)
}

func Get(key string) func(config.Config) plug.Plug {
	if val, ok := plugs[strings.ToLower(key)]; !ok {
		return nil
	} else {
		return val
	}
}

func GetAll() PlugRegistry {
	res := PlugRegistry{}
	for k, v := range plugs {
		res[k] = v
	}
	return res
}
