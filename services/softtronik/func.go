package softtronik

import (
	softtronikTypesRest "MerlionScript/types/restTypes/softtronik"
	"strings"
)

func getGlobalItemsRecord(id string, GlobalItems []softtronikTypesRest.ProductItem) (record softtronikTypesRest.ProductItem, found bool) {
	for i := range GlobalItems {
		if GlobalItems[i].Code == id {
			return GlobalItems[i], true
		}
	}
	return softtronikTypesRest.ProductItem{}, false
}

func getAvailsItemsRecords(SofttronikItems softtronikTypesRest.StocksItem) (itemAvails map[string]softtronikTypesRest.ItemStockPrice) {
	itemAvails = make(map[string]softtronikTypesRest.ItemStockPrice, len(SofttronikItems.Body.ProductsDataWithPricesAndBalances))
	for _, item := range SofttronikItems.Body.ProductsDataWithPricesAndBalances {
		newItem := softtronikTypesRest.ItemStockPrice{
			Stocks: 0,
			Price:  0,
		}
		for _, price := range item.Prices {
			if price.PriceType == "DillerPrice" {
				newItem.Price = price.Price
			}
		}
		for _, stocks := range item.StockBalance {
			if strings.Contains(strings.ToLower(stocks.Warehouse), "москва") {
				newItem.Stocks += stocks.Rest
			}
		}
		itemAvails[item.Article] = newItem
	}
	return itemAvails
}

func getExtensionFromContentType(contentType string) string {
	switch contentType {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}
