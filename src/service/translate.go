package service

import (
	"chimera/network"
	"chimera/network/myspace"
	"chimera/network/oscar"
	"chimera/utility/logging"
	"fmt"
)

func ServiceTranslateGetInterOp(msg *network.ServiceMessage) ServiceInterOp {
	var interop ServiceInterOp

	for ix := 0; ix < len(network.Clients); ix++ {
		if network.Clients[ix].ClientAccount.UIN == msg.Data.Sender {
			interop.Client = network.Clients[ix]

			if interop.Client == nil {
				logging.Error("Service/ServiceTranslateGetInterOp", "sender offline! skipping....")
				return ServiceInterOp{}
			}
		}

		switch msg.Service {
		case network.Service_MSIM:
			if myspace.ClientContexts[ix].UIN == msg.Data.Sender {
				interop.MSIM_Context = myspace.ClientContexts[ix]

				if interop.MSIM_Context == nil {
					logging.Error("Service/ServiceTranslateGetInterOp", "sender contextless! skipping....")
					return ServiceInterOp{}
				}
			}

		case network.Service_OSCAR:
			if oscar.ClientContexts[ix].UIN == msg.Data.Sender {
				interop.OSCAR_Context = oscar.ClientContexts[ix]

				if interop.OSCAR_Context == nil {
					logging.Error("Service/ServiceTranslateGetInterOp", "sender contextless! skipping....")
					return ServiceInterOp{}
				}
			}

		default:
			logging.Warn("Service/ServiceTranslateGetInterOp", "unknown messenger, skipping....")
			return ServiceInterOp{}
		}
	}

	return interop
}

func ServiceTranslateToMsimStatus(code int, message string) (int, string) {

	switch code {
	case StatusCode_Offline:
		return 0x0000, message
	case StatusCode_Hidden:
		return 0x0000, message
	case StatusCode_Online:
		return 0x0001, message
	case StatusCode_Idle:
		return 0x0002, message
	case StatusCode_Away:
		return 0x0005, message
	// this is where the fun begins
	case StatusCode_DoNotDisturb:
		return 0x0005, fmt.Sprintf("(DND) %s", message)
	case StatusCode_OutToLunch:
		return 0x0002, fmt.Sprintf("(Out on Lunch) %s", message)
	case StatusCode_OnThePhone:
		return 0x0001, fmt.Sprintf("(Mobile) %s", message)
	}

	return 0x0000, "! Status Error !"
}
