package merlion

import (
	"MerlionScript/types/soapTypes"
	"MerlionScript/utils/soap"
	"encoding/xml"
	"fmt"
	"log"
)

func GetItemsImagesByItemId(itemId string) []soapTypes.ItemImage {
	//fmt.Println(base64.StdEncoding.EncodeToString([]byte(credentials)))
	req := soapTypes.ItemImageReq{
		Item_id: []soapTypes.ItemId{{Item: itemId}},
		//Page:   "1",
	}
	//var res = make([]types.ItemMenu, 100)
	decoder, err := soap.SoapCallHandleResponse("https://apitest.merlion.com/rl/mlservice3", soapTypes.GetItemsImagesUrl, req)
	if err != nil {
		log.Fatalf("SoapCallHandleResponse error: %s", err)
	}
	var item soapTypes.ItemImage
	var items []soapTypes.ItemImage
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
	return items
}
