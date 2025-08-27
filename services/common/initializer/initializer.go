package initializer

import (
	"MerlionScript/services/elektronmir"
	"MerlionScript/services/infocom"
	"MerlionScript/services/merlion"
	"MerlionScript/services/netlab"
	"MerlionScript/services/sklad"
	"MerlionScript/services/softtronik"
	"MerlionScript/services/torus"
)

func InitServices() {
	sklad.Init()
	merlion.Init()
	netlab.Init()
	softtronik.Init()
	elektronmir.Init()
	infocom.Init()
	torus.Init()
}
