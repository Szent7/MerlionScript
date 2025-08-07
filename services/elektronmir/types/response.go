package types

type TokenResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

type CatalogResponse struct {
	Status string `json:"status"`
	Total  int    `json:"total"`
	Data   []Item `json:"data"`
}

type Item struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Count    int    `json:"count"`
	ParentID *int   `json:"parent_id"`
}

type ItemResponse struct {
	Status string `json:"status"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Total  int    `json:"total"`
	Data   []Data `json:"data"`
}

type Data struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Article     string   `json:"article"`
	Description string   `json:"description"`
	Vendor      string   `json:"vendor"`
	Category    Category `json:"category"`
	Color       string   `json:"color"`
	EAN         string   `json:"ean"`
	Weight      float64  `json:"weight"`
	Length      float64  `json:"length"`
	Width       float64  `json:"width"`
	Height      float64  `json:"height"`
	Volume      float64  `json:"volume"`
	Country     string   `json:"country"`
	Properties  []any    `json:"properties"`
	Photos      []string `json:"photos"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type StockResponse struct {
	Status string      `json:"status"`
	Total  int         `json:"total"`
	Data   []StockData `json:"data"`
}

type StockData struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Article     string   `json:"article"`
	EAN         string   `json:"ean"`
	Category    Category `json:"category"`
	Offers      []Offer  `json:"offers"`
	VendorName  string   `json:"vendor_name"`
	OnlinerID   string   `json:"onliner_id"`
	OnlinerName string   `json:"onliner_name"`
}

type Offer struct {
	ExtID                     string  `json:"ext_id"`
	Price                     float64 `json:"price"`
	PriceRRC                  float64 `json:"price_rrc"`
	PriceOld                  float64 `json:"price_old"`
	Currency                  string  `json:"currency"`
	Quantity                  int     `json:"quantity"`
	Status                    string  `json:"status"`
	Warranty                  int     `json:"warranty"`
	CompanyUNP                string  `json:"company_unp"`
	WarehouseID               int     `json:"warehouse_id"`
	DeliveryTime              *string `json:"delivery_time,omitempty"`
	Promotion                 *string `json:"promotion,omitempty"`
	OrganisationCode          int     `json:"organisation_code"`
	WarehouseOrganisationCode int     `json:"warehouse_organisation_code"`
	WarehouseCode             string  `json:"warehouse_code"`
}
