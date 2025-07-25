package types

type CategoryItem struct {
	Name        string `json:"Name"`
	ID          string `json:"ID"`
	Code        string `json:"Code"`
	TradeMark   string `json:"TradeMark"`
	Description string `json:"Description"`
	ParentName  string `json:"ParentName"`
}

type ProductItem struct {
	Article         string `json:"Article"`
	Code            string `json:"Code"`
	Name            string `json:"Name"`
	UnitMeasureName string `json:"UnitMeasureName"`
	UnitMeasureCode string `json:"UnitMeasureCode"`
}

type StocksItem struct {
	Header header `json:"Header"`
	Body   body   `json:"Body"`
}

type header struct {
	Code int `json:"Code"`
}

type body struct {
	Currency                          []currency `json:"currency"`
	ProductsDataWithPricesAndBalances []product  `json:"ProductsDataWithPricesAndBalances"`
	Date                              string     `json:"_date"`
}

type currency struct {
	ID   string `json:"_id"`
	Rate any    `json:"_rate"`
}

type product struct {
	Article      string         `json:"Article"`
	Code         string         `json:"Code"`
	Prices       []price        `json:"Prices"`
	StockBalance []stockBalance `json:"StockBalance"`
}

type price struct {
	PriceType  string  `json:"PriceType"`
	CurrencyId string  `json:"CurrencyId"`
	Price      float64 `json:"Price"`
}

type stockBalance struct {
	Warehouse string `json:"Warehouse"`
	RestTotal int    `json:"RestTotal"`
	Rest      int    `json:"Rest"`
	Transit   int    `json:"Transit"`
}

type ImageItem struct {
	Article   string     `json:"Article"`
	Code      string     `json:"Code"`
	DataFiles []dataFile `json:"DataFiles"`
}

type dataFile struct {
	NameFile             string `json:"NameFile"`
	Category             string `json:"Category"`
	MainImage            string `json:"MainImage"`
	DateModificationFile string `json:"DateModificationFile"`
	ExtensionFile        string `json:"ExtensionFile"`
	SizeFile             string `json:"SizeFile"`
	LinkToDataFile       string `json:"LinkToDataFile"`
}

type ItemStockPrice struct {
	Stocks int
	Price  float64
}
