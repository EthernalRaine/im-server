package bridge

import (
	"chimera/utility/logging"
)

func Initialize() {
	go BridgeNewCommandListener()
}

func SignOnService(svcname string, svcrev string, svccfg bool, svc func()) {

	if !svccfg {
		return
	}

	service := BridgeClient{
		SerivceName: svcname,
		ServiceRev:  svcrev,
	}

	clients = append(clients, &service)

	logging.Info("Service Bridge", "%s Service SignOn (Rev: %s)", svcname, svcrev)

	go svc()
}
