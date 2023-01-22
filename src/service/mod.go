package service

import (
	"chimera/network"
	"chimera/utility/logging"

	"golang.org/x/exp/slices"
)

func Launch(name string, enabled bool, runner func()) {
	if enabled {
		logging.Info("Service Launch", "Started service: %s", name)
		go runner()
	}
}

func Initialize() {
	logging.Info("Service Bridge", "Started cross-protocol translator!")
	go ServiceFindPackets()
}

func ServiceFindPackets() {
	for {
		for ix := 0; ix < len(network.MessageCache); ix++ {

			logging.Debug("Service/ServiceFindPackets", "ServiceMessage -> Service: %d", network.MessageCache[ix].Service)
			logging.Debug("Service/ServiceFindPackets", "ServiceMessage -> Type: %d", network.MessageCache[ix].Type)
			logging.Debug("Service/ServiceFindPackets", "ServiceData    -> Sender: %d", network.MessageCache[ix].Data.Sender)
			logging.Debug("Service/ServiceFindPackets", "ServiceData    -> Recv: %d", network.MessageCache[ix].Data.Recv)
			logging.Debug("Service/ServiceFindPackets", "ServiceData    -> Raw: %d", network.MessageCache[ix].Data.Raw)

			ServiceFilterPackets(network.MessageCache[ix])

			network.MessageCache = slices.Delete(network.MessageCache, ix, ix+1)
		}
	}
}

func ServiceFilterPackets(msg *network.ServiceMessage) {
	switch msg.Type {
	case network.MessageType_Status:
	case network.MessageType_SignOn:
		ServiceActionBroadcastSignOnStatus(msg)
	case network.MessageType_LogOff:
	case network.MessageType_IM:
	case network.MessageType_AddFriend:
	case network.MessageType_DelFriend:
	default:
		logging.Warn("Service/ServiceFilterPackets", "Unknown MessageType, skipping...")
	}
}
