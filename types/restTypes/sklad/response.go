package sklad

type SearchItem struct {
	Rows []Rows `json:"rows"`
}

type Rows struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Group             Group  `json:"group"`
	Article           string `json:"article"`
	IsSerialTrackable bool   `json:"isSerialTrackable"`
	Meta              Meta   `json:"meta"`
}

type Group struct {
	Meta Meta `json:"meta"`
}

type SearchStoreOrganization struct {
	Rows []RowsStoreOrganization `json:"rows"`
}

type RowsStoreOrganization struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	StoreMeta Meta   `json:"meta"`
}

type SearchStock struct {
	AssortmentId string `json:"assortmentId"`
	StoreId      string `json:"storeId"`
	Stock        int    `json:"stock"`
}

type SearchImage struct {
	Rows []RowsImages `json:"rows"`
}

type RowsImages struct {
	ImageMeta Meta   `json:"meta"`
	Title     string `json:"title"`
	Filename  string `json:"filename"`
	Size      int    `json:"size"`
	Updated   string `json:"updated"`
	Miniature Meta   `json:"miniature"`
	Tiny      Meta   `json:"tiny"`
}
