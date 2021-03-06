package helpers

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"runtime"
)

func Homepath() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	} else {
		return os.Getenv("HOME")
	}
}

func Unpanic(onRecover func(interface{})) {
	if err := recover(); err != nil {
		trace := make([]byte, 4096)
		count := runtime.Stack(trace, true)

		log.Printf("Recovered from panic: %v\nStack has %d bytes ->\n%s\n", err, count, trace)

		if onRecover != nil {
			onRecover(err)
		}
	}
}

func ContainsString(s []string, k string) bool {
	for _, v := range s {
		if v == k {
			return true
		}
	}
	return false
}
