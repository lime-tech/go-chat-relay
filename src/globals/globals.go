package globals

import (
	"os"
	"path"
)

var (
	version string
	name    = path.Base(os.Args[0])
)

func Version() string { return version }
func Name() string    { return name }
