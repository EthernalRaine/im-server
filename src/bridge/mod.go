package bridge

import (
	"chimera/utility/logging"
	"sync"
)

func Initialize() {
	logging.Info("Bridge Service", "Started Listener and Router")
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

var sendDataMutex sync.RWMutex

func BridgeSendData(data string) {
	sendDataMutex.Lock()
	messages = append(messages, data)
	sendDataMutex.Unlock()
}
