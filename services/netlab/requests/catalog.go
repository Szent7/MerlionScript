package requests

import (
	netlabTypes "MerlionScript/services/netlab/types"
	"MerlionScript/utils/rest"
	"encoding/xml"
	"fmt"
	"log"
	"strings"
)

func GetAllCategoryCodes(token string) ([]netlabTypes.ItemCategory, error) {
	url := fmt.Sprintf(netlabTypes.CategoryUrl, netlabTypes.CatalogName, token)

	decoder, err := rest.CreateRequestXML("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var codes []netlabTypes.ItemCategory

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		var item netlabTypes.ItemCategory
		switch start := token.(type) {
		case xml.StartElement:
			if start.Name.Local == "category" {
				err := decoder.DecodeElement(&item, &start)
				if err != nil {
					log.Println("ошибка при декодировании item(getcatalog):", err)
					break
				}
				lower := strings.ToLower(item.Name)
				if strings.Contains(lower, "dahua") || strings.Contains(lower, "tenda") {
					codes = append(codes, item)
				}
			}
		}
	}
	return codes, nil
}
