package merlion

import "MerlionScript/utils/db/typesDB"

// var code = "I00001"
var counter = 1

type merlionEntity struct {
	MoySkladCode     string
	ManufacturerCode string
	MerlionCode      string
}

func CreateMerlionEntity(codes typesDB.Codes) *merlionEntity {
	return &merlionEntity{
		MoySkladCode:     codes.MoySklad,
		ManufacturerCode: codes.Manufacturer,
		MerlionCode:      codes.Merlion,
	}
}
