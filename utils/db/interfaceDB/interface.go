package interfaceDB

import (
	"MerlionScript/utils/db/typesDB"
)

type DB interface {
	InsertCodesIDs(record typesDB.CodesIDs) (bool, error)
	GetCodesFilledMSNoImage(tableName string) (*[]typesDB.General, error)
	UpdateCodesIDs(record typesDB.CodesIDs) error
	GetCodesIDs(article string) (typesDB.CodesIDs, error)
	GetLastOwnIdMS() (int64, error)
	GetCodesFilledMS(tableName string) (*[]typesDB.General, error)

	InsertService(record typesDB.CodesService, tableName string) (bool, error)
	GetService(article string, tableName string) (typesDB.CodesService, error)
	GetServiceIDsNoMS(tableName string) (*[]typesDB.CodesService, error)
	GetServiceIDsFilledMS(tableName string) (*[]typesDB.CodesService, error)
	GetServiceIDsFilledMSWithNoImage(tableName string) (*[]typesDB.CodesService, error)
	GetAllServicesByArticle(article string) (*[]typesDB.CodesService, error)
	UpdateService(record typesDB.CodesService, tableName string) error
}
