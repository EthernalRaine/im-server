package service

import (
	"fmt"
)

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
