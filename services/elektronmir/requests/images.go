package requests

import (
	netlabTypes "MerlionScript/services/netlab/types"
	"MerlionScript/utils/rest"
	"encoding/xml"
	"fmt"
)

func getImagesByItemId(id string, token string) ([]netlabTypes.ItemImage, error) {
	url := fmt.Sprintf(netlabTypes.ImageUrl, id, token)

	decoder, err := rest.CreateRequestXML("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var item netlabTypes.ItemImage
	var items []netlabTypes.ItemImage

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		switch start := token.(type) {
		case xml.StartElement:
			if start.Name.Local == "item" {
				err := decoder.DecodeElement(&item, &start)
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
