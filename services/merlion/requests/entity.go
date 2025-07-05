package requests

import "MerlionScript/utils/db/typesDB"

type merlionEntity struct {
	MoySkladCode     string
	ManufacturerCode string
	MerlionCode      string
}

func CreateMerlionEntity(codes typesDB.Codes) *merlionEntity {
	return &merlionEntity{
		MoySkladCode:     codes.MoySklad,
		ManufacturerCode: codes.Manufacturer,
		MerlionCode:      codes.Service,
	}
}
