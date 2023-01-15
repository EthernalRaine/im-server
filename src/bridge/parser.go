package bridge

import (
	"chimera/utility/encryption"
	"chimera/utility/logging"

	"golang.org/x/exp/slices"
)

func BridgeNewCommandListener() {
	for {
		for ix := 0; ix < len(messages); ix++ {
			logging.Trace("Bridge/Listener", "Data Len: %d, Data Content: %s", len(messages[ix]), messages[ix])
			if BridgeVerifyCommand(messages[ix]) {
				BridgeCommandParser(messages[ix])
			} else {
				logging.Error("Bridge/Listener", "failed to verify data crc, skipping...")
			}
			messages = slices.Delete(messages, ix, ix+1)
		}
	}
}

func BridgeVerifyCommand(data string) bool {
	datasplice := BridgeSplitDataPacket(data)[0]
	datasplice += "\x1c"
	datasplice += BridgeSplitDataPacket(data)[1]
	datasplice += "\x1c"

	datacrc := BridgeRetrieveDataCRC(data)
	newcrc := encryption.GetCRCHash([]byte(datasplice))

	logging.Debug("Bridge/CRC", "CommandData CRC: %d", datacrc)
	logging.Debug("Bridge/CRC", "SplicedData CRC: %d", newcrc)

	return datacrc == newcrc
}

/*
	supply 0 for RecvUIN if the action does not require it

	---------------------------------------

	ac - action <- This *needs* to be the Key for any of the following

	---------------------------------------

	ns - new status
	im - instant message
	sb - signon broadcast
	lb - logoff broadcast
	af - add friend
	df - delete friend
	bf - block user
	uf - unblock user
	ic - information change (pfp or similar)
	fc - fetch client
	fi - fetch information (userdetails)

*/

func BridgeCommandParser(data string) {

	actionkey := BridgeRetrieveValueFromActionDataPair("ac", data)
	actiondata := BridgeSplitDataPacket(data)[1]
	actionheader := BridgeRetrieveTransportHeaderInformation(data)

	logging.Debug("Bridge/Parser", "Command Header: %+v", actionheader)
	logging.Debug("Bridge/Parser", "Command Action Data: %s", BridgeFormatActionData(actiondata))

	switch actionkey {
	case "ns":
		BridgeHandleActionNewStatusMessage(actiondata)
	case "im":
		BridgeHandleActionSendInstantMessage(actiondata)
	case "sb":
		BridgeHandleActionBroadcastSignon(actionheader, actiondata)
	case "lb":
		BridgeHandleActionBroadcastLogoff(actionheader, actiondata)
	case "af":
		BridgeHandleActionAddContact(actiondata)
	case "df":
		BridgeHandleActionDeleteContact(actiondata)
	case "bf":
		BridgeHandleActionBlockUser(actiondata)
	case "uf":
		BridgeHandleActionUnblockUser(actiondata)
	case "ic":
		BridgeHandleActionChangeServerInformation(actiondata)
	case "fc":
		BridgeHandleActionFetchClientDetails(actiondata)
	case "fi":
		BridgeHandleActionFetchUserDetails(actiondata)
	default:
		logging.Warn("Bridge/Parser", "unknown command detected, ignoring")
	}
}
