package typesDB

type Codes struct {
	Id               int
	MsOwnId          int64
	MoySklad         string
	Manufacturer     string
	ManufacturerName string
	Service          string
	LoadedImage      int
	TryLoadImage     int
}

const (
	MerlionTable = "codes_merlion"
	NetlabTable  = "codes_netlab"
	OwnIDsTable  = "codes_ids"

	MerlionService = "merlion"
	NetlabService  = "netlab"
)
