package service

import (
	"chimera/network"
	"chimera/network/myspace"
	"chimera/network/oscar"
)

type ServiceInterOp struct {
	Client        *network.Client
	MSIM_Context  *myspace.MySpaceContext
	OSCAR_Context *oscar.OSCARContext
}

const (
	StatusCode_Offline      = 0x0000
	StatusCode_Hidden       = 0x0001
	StatusCode_Online       = 0x0002
	StatusCode_Idle         = 0x0003
	StatusCode_Away         = 0x0004
	StatusCode_DoNotDisturb = 0x0005
	StatusCode_OutToLunch   = 0x0006
	StatusCode_OnThePhone   = 0x0007
)
