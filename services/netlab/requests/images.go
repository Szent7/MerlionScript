package requests

import (
	netlabTypesRest "MerlionScript/types/restTypes/netlab"
	netlabTypesSoap "MerlionScript/types/soapTypes/netlab"
	"MerlionScript/utils/rest"
	"encoding/xml"
	"fmt"
)

func getImagesByItemId(id string, token string) ([]netlabTypesSoap.ItemImage, error) {
	url := fmt.Sprintf(netlabTypesRest.ImageUrl, id, token)

	decoder, err := rest.CreateRequestXML("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var item netlabTypesSoap.ItemImage
	var items []netlabTypesSoap.ItemImage

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		switch start := token.(type) {
		case xml.StartElement:
			//fmt.Println("Xml tag:", start.Name.Local)
			if start.Name.Local == "item" {
				err := decoder.DecodeElement(&item, &start)
				//fmt.Println(item)
				//res = append(res, item)
				if err != nil {
					fmt.Println("Ошибка при декодировании item(GetItemsByCatId):", err)
					break
				}

				items = append(items, item)
			}
		}
	}
	return items, nil
}

func GetImagesByItemIdFormatted(id string, token string) (map[string]string, error) {
	rawImages, err := getImagesByItemId(id, token)
	if err != nil {
		return nil, err
	}

	var images map[string]string = make(map[string]string, 5)
	for i := range rawImages {
		for _, property := range rawImages[i].Properties.Property {
			switch property.Name {
			case "Url":
				filename := fmt.Sprintf("%s_%d", id, (i + 1))
				images[filename] = property.Value
			}
		}
	}

	return images, nil
}
