package types

const (
	ServiceName   = "softtronik"
	ServiceDBName = "codes_softtronik"

	CategoryUrl = "https://api.soft-tronik.ru/API/hs/Products/GetCommodityGroupsList/%s"
	ItemUrl     = "https://api.soft-tronik.ru/API/hs/Products/GetProductList/%s?CommodityGroupKey=%s"
	StocksUrl   = "https://api.soft-tronik.ru/API/hs/Products/ProductsPriceAndRest?ContactPersonKey=%s&CommodityGroupKey=%s&ContractKey=%s"
	ImageUrl    = "https://api.soft-tronik.ru/API/hs/Products/GetProductDataFilesList/%s/%s?TypeFiles=jpg,png"
)
