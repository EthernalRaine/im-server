package main

import (
	"chimera/bridge"
	"chimera/network/aim"
	"chimera/network/msn"
	"chimera/network/myspace"
	"chimera/utility"
	"chimera/utility/configuration"
	"chimera/utility/database"
	"chimera/utility/logging"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logging.Info("Main", "Starting Chimera...")
	logging.Info("Main", "Build Info: [%s]", utility.GetBuild())

	database.Initialize()
	bridge.Initialize()

	svc := configuration.GetConfiguration().Services
	bridge.SignOnService("MySpace", "2.0", svc.MySpace, myspace.LogonMySpace)
	bridge.SignOnService("MSN", "1.0", svc.MSN, msn.LogonMSN)
	bridge.SignOnService("AIM", "1.0", svc.AIM, aim.LogonAIM)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	for sig := range c {
		logging.Info("Exit Handler", "Captured %v! Stopping Server...", sig)
		os.Exit(0)
	}
}
