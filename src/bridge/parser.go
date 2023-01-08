package bridge

import "chimera/utility/logging"

func BridgeNewCommandListener() {
	for ix := 0; ix < len(messages); ix++ {
		logging.Trace("Bridge/BridgeNewCommandListener", "")
	}
}
