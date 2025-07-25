package merlion

import (
	"MerlionScript/services/common"
	"context"
	"log"
)

func Init() {
	common.RegisteredServices["merlion"] = common.ServiceInfo{
		DBTableName: "codes_merlion",
		Controller:  Controller,
	}
	log.Println("merlion: сервис зарегистрирован")
}

func Controller(ctx context.Context) {

}
