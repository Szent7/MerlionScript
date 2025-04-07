package restTypes

type SearchItem struct {
	Rows []Rows `json:"rows"`
}

type Rows struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Group   Group  `json:"group"`
	Article string `json:"article"`
	Meta    Meta   `json:"meta"`
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
