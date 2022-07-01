package helper

import (
	"log"
	"os"

	"github.com/georzaza/go-ethereum-v0.7.10_official/ethutil"
	logpkg "github.com/georzaza/go-ethereum-v0.7.10_official/logger"
)

var Logger logpkg.LogSystem
var Log = logpkg.NewLogger("TEST")

func init() {
	Logger = logpkg.NewStdLogSystem(os.Stdout, log.LstdFlags, logpkg.InfoLevel)
	logpkg.AddLogSystem(Logger)

	ethutil.ReadConfig(".ethtest", "/tmp/ethtest", "")
	ethutil.Config.Db, _ = NewMemDatabase()
}
