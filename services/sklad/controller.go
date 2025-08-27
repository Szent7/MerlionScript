package sklad

import (
	"MerlionScript/services/common"
	skladReq "MerlionScript/services/sklad/requests"
	skladTypes "MerlionScript/services/sklad/types"
	"MerlionScript/utils/db/typesDB"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type SkladERP struct{}

func Init() {
	common.MainERP = new(SkladERP)
	log.Printf("%s: система ERP зарегистрирована\n", skladTypes.ServiceName)
}

func (srv *SkladERP) GetSystemName() string {
	return skladTypes.ServiceName
}
func (srv *SkladERP) GetDBTableName() string {
	return skladTypes.ServiceDBName
}

func (srv *SkladERP) GetItemsByArticle(article string) (*[]common.ItemList, error) {
	manufacturerReplace := strings.ReplaceAll(article, " ", "+")
	response, err := skladReq.GetItemByManufacturer(manufacturerReplace)
	if err != nil || response.StatusCode != 200 {
		log.Printf("ошибка при получении записи из МС (GetItemsByArticle) article = %s StatusCode = %d: %s, \n", article, response.StatusCode, err)
		return nil, err
	}
	msItems := skladTypes.SearchItem{}
	if err := json.Unmarshal(response.Body, &msItems); err != nil {
		log.Printf("Ошибка при декодировании item (createNewPositionsMS) manufacturer = %s: %s", article, err)
		return nil, err
	}

	itemList := make([]common.ItemList, 0, len(msItems.Rows))
	for i := range msItems.Rows {
		itemList = append(itemList, common.ItemList{
			Article:      msItems.Rows[i].Article,
			PositionName: msItems.Rows[i].Name,
		})
	}

	return &itemList, nil
}

func (srv *SkladERP) GetItemID(code string) (string, error) {
	itemMS, err := skladReq.GetItem(code)
	if err != nil || itemMS.Id == "" {
		log.Printf("ошибка при получении ID товара МС (updateRemainsMS) msCode = %s: %s\n", code, err)
		return "", err
	}
	return itemMS.Id, nil
}

func (srv *SkladERP) GetItemAvails(code string, storeUUID string) (common.StockERP, error) {
	itemMS, err := skladReq.GetItem(code)
	if err != nil || itemMS.Id == "" {
		log.Printf("Ошибка при получении товара МС (updateRemainsMS) msCode = %s: %s\n", code, err)
		return common.StockERP{}, err
	}

	itemRemainsMS, err := skladReq.GetItemsAvail(itemMS.Id, storeUUID)
	if err != nil {
		log.Printf("Ошибка при получении остатков с мс (updateRemainsMS) msCode = %s: %s\n", code, err)
		return common.StockERP{}, err
	}

	return common.StockERP{
		ItemMeta:          itemMS.Meta,
		Stock:             itemRemainsMS,
		IsSerialTrackable: itemMS.IsSerialTrackable,
	}, nil
}

func (srv *SkladERP) GetImagesList(id string) (skladTypes.SearchImage, error) {
	response, err := skladReq.GetItemsImagesData(id)
	if err != nil || response.StatusCode != 200 {
		log.Printf("Ошибка при получении записей из МС (UploadAllImages): %s\n", err)
		return skladTypes.SearchImage{}, err
	}

	msImages := skladTypes.SearchImage{}
	if err := json.Unmarshal(response.Body, &msImages); err != nil {
		log.Printf("Ошибка при декодировании msImages (UploadAllImages) manufacturer = %s: %s", id, err)
		return skladTypes.SearchImage{}, err
	}

	return msImages, nil
}

func (srv *SkladERP) CreateItem(item *typesDB.CodesIDs, newId string, itemName string, catalog string) error {
	catalogMeta, err := srv.GetCatMeta(catalog)
	if err != nil {
		return err
	}
	//Если мета каталога пустая
	var newItem skladTypes.CreateItem
	if catalogMeta.Href == "" {
		newItem = skladTypes.CreateItem{
			Name:    itemName,
			Article: newId,
		}
	} else {
		newItem = skladTypes.CreateItem{
			Name:    itemName,
			Article: newId,
			ProductFolder: skladTypes.MetaMiddle{
				Meta: catalogMeta,
			},
		}
	}
	response, err := skladReq.CreateItem(newItem)
	if err != nil || response.StatusCode != 200 {
		log.Printf("Ошибка при создании записи на МС (createNewPositionsMS) erpCode = %s: %s\n", item.MoySkladCode, err)
		fmt.Println(response.StatusCode)
		fmt.Println(string(response.Body))
		return err
	}

	return nil
}

func (srv *SkladERP) CreateAcceptance(acceptanceReq *skladTypes.Acceptance) error {
	return skladReq.IncreaseItemsAvail(acceptanceReq)
}

func (srv *SkladERP) CreateWoff(woffReq *skladTypes.Acceptance) error {
	return skladReq.DecreaseItemsAvail(woffReq)
}

func (srv *SkladERP) UploadImage(imageData skladTypes.UploadImage, id string) error {
	resp, err := skladReq.UploadImage(id, imageData)
	if err != nil || resp.StatusCode != 200 {
		log.Printf("Ошибка при загрузке изображения на МС (UploadAllImages) softtronikCode = %s: %s\n", id, err)
	}
	return err
}

func (srv *SkladERP) GetCatMeta(catalog string) (skladTypes.Meta, error) {
	return skladReq.GetCatMeta(catalog)
}

func (srv *SkladERP) GetOrgMeta(organization string) (skladTypes.Meta, error) {
	return skladReq.GetOrganizationMeta(organization)
}

func (srv *SkladERP) GetStoreMeta(store string) (skladTypes.Meta, error) {
	return skladReq.GetStoreMeta(store)
}

func (srv *SkladERP) GetStoreUUID(store string) (string, error) {
	return skladReq.GetStoreUUID(store)
}
