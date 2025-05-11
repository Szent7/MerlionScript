package requests

import (
	netlabTypesRest "MerlionScript/types/restTypes/netlab"
	netlabTypesSoap "MerlionScript/types/soapTypes/netlab"
	"MerlionScript/utils/rest"
	"encoding/xml"
	"fmt"
	"strconv"
)

func GetCurrency(token string) (float64, error) {
	url := fmt.Sprintf(netlabTypesRest.CurrencyUrl, token)

	decoder, err := rest.CreateRequestXML("GET", url, nil)
	if err != nil {
		return 0, err
	}

	var currency float64 = 0
	var item netlabTypesSoap.Property
	var status netlabTypesSoap.Status

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		switch start := token.(type) {
		case xml.StartElement:
			if start.Name.Local == "status" {
				err := decoder.DecodeElement(&status, &start)
				if err != nil {
					fmt.Println("Ошибка при декодировании status(GetCurrency):", err)
					break
				}
			}

			if start.Name.Local == "property" {
				err := decoder.DecodeElement(&item, &start)
				if err != nil {
					fmt.Println("Ошибка при декодировании item(GetCurrency):", err)
					break
				}
				if item.Name == "usdRateNonCash" {
					currency, err = strconv.ParseFloat(item.Value, 64)
					if err != nil {
						return 0, fmt.Errorf("ошибка при получении курса валют: %s\n", err.Error())
					}
				}
			}
		}
	}
	if currency == 0 || status.Code != "200" {
		return 0, fmt.Errorf("ошибка при получении курса валют: %s\n", status.Message)
	}
	return currency, nil
}
