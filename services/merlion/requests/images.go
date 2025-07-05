package requests

import (
	"MerlionScript/keeper"
	merlionTypes "MerlionScript/types/soapTypes/merlion"
	"MerlionScript/utils/soap"
	"encoding/xml"
	"fmt"
	"log"
	"path"
	"strings"
)

func getImagesByItemId(itemId string) ([]merlionTypes.ItemImage, error) {
	req := merlionTypes.ItemImageReq{
		Item_id: merlionTypes.ItemId{Item: []string{itemId}},
	}
	decoder, err := soap.SoapCallHandleResponse(keeper.MerlionMainURL, merlionTypes.GetItemsImagesUrl, req)
	if err != nil {
		log.Printf("Ошибка при SOAP-запросе (getImagesByItemId): %s\n", err)
		return nil, err
	}
	var item merlionTypes.ItemImage
	var items []merlionTypes.ItemImage
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
					fmt.Println("ошибка при декодировании item(getitemsimages):", err)
					break
				}

				items = append(items, item)
			}
		}
	}
	return items, nil
}

func GetImagesByItemIdFormatted(id string) ([]string, error) {
	rawImages, err := getImagesByItemId(id)
	if err != nil {
		return nil, err
	}

	var images []string = make([]string, 0, len(rawImages)/3)
	for i := range rawImages {
		ext := path.Ext(rawImages[i].FileName)
		if ext == ".png" || ext == ".jpg" {

			base := rawImages[i].FileName[:len(rawImages[i].FileName)-len(ext)]

			parts := strings.Split(base, "_")
			if len(parts) < 3 {
				continue
			}

			if parts[2] == "m" {
				images = append(images, rawImages[i].FileName)
			}
		}
	}

	return images, nil
}
