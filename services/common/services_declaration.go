package common

import (
	"context"
)

type ServiceInfo struct {
	DBTableName string
	Controller  func(context.Context)
}

//K - ServiceName
//V - ServiceInfo
var RegisteredServices map[string]ServiceInfo

func GetTableNames() []string {
	if len(RegisteredServices) == 0 {
		return nil
	}
	tableNames := make([]string, len(RegisteredServices))
	for _, v := range RegisteredServices {
		tableNames = append(tableNames, v.DBTableName)
	}
	return tableNames
}

/*
var Controllers []func(context.Context) = []func(context.Context){
	merlion.Controller,
	netlab.Controller,
	softtronik.Controller,
}

//Название сервиса : Название таблицы БД
var ServicesDec map[string]string = map[string]string{
	"ownIDs":     "codes_ids",
	"merlion":    "codes_merlion",
	"netlab":     "codes_netlab",
	"softtronik": "codes_softtronik",
}
*/
