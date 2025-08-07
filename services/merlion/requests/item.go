package requests

import (
	merlionTypes "MerlionScript/services/merlion/types"
	"MerlionScript/utils/soap"
	"encoding/xml"
	"fmt"
	"log"
	"time"
)

func GetItemsByCatId(catId string) ([]merlionTypes.ItemCatalog, error) {
	req := merlionTypes.ItemCatalogReq{
		Cat_id: catId,
	}

	decoder, err := soap.SoapCallHandleResponse(merlionTypes.MerlionMainURL, merlionTypes.GetItemsUrl, req)
	if err != nil {
		log.Printf("Ошибка при SOAP-запросе (GetItemsByCatId): %s\n", err)
		return nil, err
	}
	var item merlionTypes.ItemCatalog
	var items []merlionTypes.ItemCatalog
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

func GetItemsAvailByItemIdBatch(itemId []string) (*[]merlionTypes.ItemAvail, error) {
	req := merlionTypes.ItemAvailReq{
		Item_id:         merlionTypes.ItemId{Item: itemId},
		Shipment_method: "ДОСТАВКА",
		Shipment_date:   GetNextDate(),
		//Page:   "1",
	}

	decoder, err := soap.SoapCallHandleResponse(merlionTypes.MerlionMainURL, merlionTypes.GetItemsAvailUrl, req)
	if err != nil {
		log.Printf("Ошибка при SOAP-запросе (GetItemsAvailByItemId): %s\n", err)
		return nil, err
	}
	var item merlionTypes.ItemAvail
	var items = make([]merlionTypes.ItemAvail, 0, len(itemId))
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
					fmt.Println("ошибка при декодировании item(getitemsavail):", err)
					break
				}

				items = append(items, item)
			}
		}
	}
	return &items, nil
}

func GetItemsByItemIdBatch(itemId []string) (*[]merlionTypes.ItemCatalog, error) {
	req := merlionTypes.ItemCatalogReq{
		Item_id: merlionTypes.ItemId{Item: itemId},
		//Page:   "1",
	}

	decoder, err := soap.SoapCallHandleResponse(merlionTypes.MerlionMainURL, merlionTypes.GetItemsUrl, req)
	if err != nil {
		log.Printf("Ошибка при SOAP-запросе (GetItemsByItemId): %s\n", err)
		return nil, err
	}
	var item merlionTypes.ItemCatalog
	var items = make([]merlionTypes.ItemCatalog, 0, len(itemId))
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
					fmt.Println("ошибка при декодировании item(getitems):", err)
					break
				}

				items = append(items, item)
			}
		}
	}
	return &items, nil
}

func GetNextDate() string {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	return tomorrow.Format("2006-01-02")
}
