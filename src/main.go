package main

import (
	"chimera/bridge"
	"chimera/network/myspace"
	"chimera/network/oscar"
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
	logging.Info("Main", "Build Info: %s (%s)", utility.GetBuildInfo(), utility.GetBuildSpecifics())

	database.Initialize()

	svc := configuration.GetConfiguration().Services
	bridge.SignOnService("MySpace", bridge.ServiceMySpace, "2.0", svc.MySpace, myspace.LogonMySpace)
	bridge.SignOnService("AIM", bridge.ServiceAIM, "1.0", svc.AIM, oscar.LaunchOSCAR)
	//bridge.SignOnService("MSN", bridge.ServiceMSN, "1.0", svc.MSN, msn.LogonMSN)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	for sig := range c {
		logging.Info("Exit Handler", "Captured %v! Stopping Server...", sig)
		os.Exit(0)
	}
}
