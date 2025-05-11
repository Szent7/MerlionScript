package requests

import (
	"MerlionScript/keeper"
	merlionTypes "MerlionScript/types/soapTypes/merlion"
	"MerlionScript/utils/soap"
	"encoding/xml"
	"fmt"
	"log"
)

func GetItemsImagesByItemId(itemId string) []merlionTypes.ItemImage {
	//fmt.Println(base64.StdEncoding.EncodeToString([]byte(credentials)))
	req := merlionTypes.ItemImageReq{
		//Item_id: merlionTypes.ItemId{{itemId}},
		//Page:   "1",
	}
	//var res = make([]types.ItemMenu, 100)
	decoder, err := soap.SoapCallHandleResponse(keeper.MerlionMainURL, merlionTypes.GetItemsImagesUrl, req)
	if err != nil {
		log.Fatalf("SoapCallHandleResponse error: %s", err)
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
	return items
}
