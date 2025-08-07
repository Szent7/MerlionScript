package initializer

import (
	"MerlionScript/services/elektronmir"
	"MerlionScript/services/merlion"
	"MerlionScript/services/netlab"
	"MerlionScript/services/sklad"
	"MerlionScript/services/softtronik"
)

func InitServices() {
	sklad.Init()
	merlion.Init()
	netlab.Init()
	softtronik.Init()
	elektronmir.Init()
}
