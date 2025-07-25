package typesDB

const (
	IDsTable = "codes_ids"
	PathDB   = `./data`

	TableIDsSQL = `CREATE TABLE IF NOT EXISTS "%s" (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"ms_own_id" INTEGER NOT NULL,
			"moy_sklad" TEXT NOT NULL,
			"article" TEXT NOT NULL UNIQUE,
			"manufacturer" TEXT NOT NULL
		);`

	TableServiceSQL = `CREATE TABLE IF NOT EXISTS "%s" (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"article" TEXT NOT NULL UNIQUE,
			"service" TEXT NOT NULL,
			"try_upload_image" INTEGER NOT NULL,
			FOREIGN KEY ("article") REFERENCES "%s"("article")
			ON UPDATE CASCADE ON DELETE CASCADE
		);`
)

type CodesIDs struct {
	Id           int
	MsOwnId      int64
	MoySkladCode string
	Article      string
	Manufacturer string
}

type CodesService struct {
	Id             int
	Article        string
	ServiceCode    string
	TryUploadImage int
}

type General struct {
	Codes   CodesIDs
	Service CodesService
}
