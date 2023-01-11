package bridge

import (
	"chimera/utility/encryption"
	"chimera/utility/logging"

	"golang.org/x/exp/slices"
)

func BridgeNewCommandListener() {
	for {
		for ix := 0; ix < len(messages); ix++ {
			logging.Trace("Bridge/BridgeNewCommandListener", "Data Len: %d, Data Content: %s", len(messages[ix]), messages[ix])
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

	return datacrc == newcrc
}

/*
	ac - action <- This *needs* to be the Key for any of the following

	---------------------------------------

	ns - new status
	im - instant message
	sb - signon broadcast
	lb - logoff broadcast
	af - add friend
	df - delete friend
	bf - block friend
	uf - unblock friend
	ic - information change (pfp or similar)
	fc - fetch client
	fi - fetch information (userdetails)

	---------------------------------------

	er - error key followed by error code
*/

func BridgeCommandParser(data string) {

	actionkey := BridgeRetrieveValueFromActionDataPair("ac", data)

	logging.Debug("Bridge/BridgeCommandParser", "Command Header: %+v", BridgeRetrieveTransportHeaderInformation(data))
	logging.Debug("Bridge/BridgeCommandParser", "Command Action Data: %s", BridgeSplitDataPacket(data)[1])
	logging.Debug("Bridge/BridgeCommandParser", "Command CRC32: %d", BridgeRetrieveDataCRC(data))

	switch actionkey {
	case "ns":
	case "im":
	case "sb":
	case "lb":
	case "af":
	case "df":
	case "bf":
	case "uf":
	case "ic":
	case "fc":
	case "fi":
	default:
		logging.Warn("Bridge/Parser", "unknown command detected, redirecting to error")
	}
}
