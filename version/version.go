package version

import (
	"fmt"
	"runtime"
)

//needs to be a var (no const)
//so that we can overwrite during linking with -X main.GITCOMMIT ...
var (
	Version   = "20160601.01"
	GITCOMMIT = "HEAD"
)

func String() string {
	return fmt.Sprintf("%s (%s), %s", Version, GITCOMMIT, runtime.Version())
}
