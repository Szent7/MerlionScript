package merlion

import (
	"MerlionScript/types/restTypes"
	"MerlionScript/types/soapTypes"
	"MerlionScript/utils/db"
	"MerlionScript/utils/rest"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
)

// var code = "I00001"
var counter = 1

type merlionEntity struct {
	No                  string
	Article             string
	Name                string
	Vendor_part         string
	Image               map[string]string // Name - Payload(base64)
	AvailableClient_MSK int
}

func CreateMerlionEntity(MerlionCode string) *merlionEntity {
	return &merlionEntity{No: MerlionCode}
}

func (me *merlionEntity) FillData(MerlionCredentials string) error {
	var err error
	err = me.getNameVendorPart(MerlionCredentials)
	if err != nil {
		return err
	}
	me.getImage(MerlionCredentials)
	err = me.getAvailableMSK(MerlionCredentials)
	if err != nil {
		return err
	}
	return nil
}

func (me *merlionEntity) getNameVendorPart(MerlionCredentials string) error {
	items := GetItemsByItemId(me.No, MerlionCredentials)
	if len(items) == 0 || items[0].No == "" {
		return fmt.Errorf("товар не найден: getnamevendorpart()")
	}
	fmt.Printf("getNameVendorPart: %v\n", items[0])
	me.Name = items[0].Name
	me.Vendor_part = items[0].Vendor_part
	return nil
}

func (me *merlionEntity) getImage(MerlionCredentials string) error {
	me.Image = make(map[string]string)
	items := GetItemsImagesByItemId(me.No, MerlionCredentials)
	if len(items) == 0 || items[0].No == "" {
		return fmt.Errorf("товар не найден: getimage()")
	}
	for _, image := range items {
		payload, err := downloadImage(soapTypes.DownloadImageUrl+"/"+image.FileName, MerlionCredentials)
		if err != nil {
			fmt.Println("изображение не загружено: " + err.Error())
			continue
		}
		me.Image[image.FileName] = payload
	}
	return nil
}

func downloadImage(url string, credentials string) (string, error) {
	//authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
	response, err := rest.CreateRequest("GET", url, "", nil)
	if err != nil || response.StatusCode != 200 {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(response.Body), nil
}

func (me *merlionEntity) getAvailableMSK(MerlionCredentials string) error {
	items := GetItemsAvailByItemId(me.No, MerlionCredentials)
	if len(items) == 0 {
		return fmt.Errorf("товар не найден: getavailablemsk()")
	}
	me.AvailableClient_MSK = items[0].AvailableClient_MSK
	return nil
}

func (me *merlionEntity) SendDataToMoySklad(SkladCredentials string) {
	created, err := me.checkIfExist(SkladCredentials)
	if !created {
		me.setRemains(SkladCredentials)
		if err != nil {
			fmt.Println("SendDataToMoySklad already exists: " + err.Error())
		}
		fmt.Println("SendDataToMoySklad already exists")
		return
	}
	uuid, _ := me.getItemUUID(SkladCredentials)
	me.uploadImage(SkladCredentials, uuid)
	me.setRemains(SkladCredentials)
}

func (me *merlionEntity) checkIfExist(SkladCredentials string) (bool, error) {
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	database, _ := db.GetDBInstance()
	record, exist, err := database.GetCodeRecordByManufacturer(me.Vendor_part)
	if err != nil {
		if !exist {
			return false, nil
		}
		return false, err
	}
	if record.MoySklad == "" {
		response, err := rest.CreateRequest("GET", restTypes.ItemUrl+"?search="+me.Vendor_part, authHeader, nil)
		if err != nil {
			return false, err
		}
		if response.StatusCode != 200 {
			return false, fmt.Errorf("ошибка при запросе (checkifexist), statuscode: %d", response.StatusCode)
		}
		items := restTypes.SearchItem{}
		if err := json.Unmarshal(response.Body, &items); err != nil {
			return false, fmt.Errorf("ошибка при декодировании item (checkifexist): %s", err.Error())
		}
		i, exist := matchSubstring(items.Rows, record.Manufacturer)
		if exist {
			record.MoySklad = items.Rows[i].Article
			database.EditCodeRecord(&record)
			me.Article = items.Rows[i].Article
			return false, nil
			//me.updateItem(SkladCredentials, items.Rows[i].Id)
		} else {
			newid := incrementID()
			me.createItem(SkladCredentials, newid)
			record.MoySklad = newid
			me.Article = newid
			err := database.EditCodeRecord(&record)
			if err != nil {
				fmt.Println("EditCodeRecord:" + err.Error())
			}
		}
		return true, nil
	} else {
		response, err := rest.CreateRequest("GET", restTypes.ItemUrl+"/"+record.MoySklad, authHeader, nil)
		if err != nil {
			return false, err
		}
		if response.StatusCode == 200 {
			return true, nil
		} else if response.StatusCode == 404 {
			newid := incrementID()
			me.createItem(SkladCredentials, newid)
			record.MoySklad = newid
			me.Article = newid
			err := database.EditCodeRecord(&record)
			if err != nil {
				fmt.Println("EditCodeRecord:" + err.Error())
				return false, fmt.Errorf("ошибка при запросе (checkifexist), statuscode: %d", response.StatusCode)
			}
			return true, nil
		}
		return false, fmt.Errorf("ошибка при запросе (checkifexist), statuscode: %d", response.StatusCode)
	}
}

func matchSubstring(rows []restTypes.Rows, substring string) (int, bool) {
	re := regexp.MustCompile(regexp.QuoteMeta(substring) + `(\s|$|\))`)

	for i, row := range rows {
		if re.MatchString(row.Name) {
			return i, true
		}
	}
	return -1, false
}

func (me *merlionEntity) createItem(SkladCredentials string, itemId string) (bool, error) {
	newItem := restTypes.CreateItem{
		Name:    me.Name,
		Article: itemId,
	}
	jsonBody, err := json.MarshalIndent(newItem, "", "  ")
	if err != nil {
		fmt.Println("ошибка при преобразовании структуры в JSON (createitem):", err)
		return false, err
	}
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	//url := "https://api.moysklad.ru/api/remap/1.2/entity/productfolder"
	_, err = rest.CreateRequest("POST", restTypes.ItemUrl, authHeader, bytes.NewBuffer(jsonBody))
	if err != nil {
		return false, err
	}
	/*
		err = json.Unmarshal(body, &meta)
		if err != nil {
			return false, err
		}*/
	return true, nil
}

func (me *merlionEntity) uploadImage(SkladCredentials string, itemId string) error {
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	for name, content := range me.Image {
		jsonBody, err := json.MarshalIndent(restTypes.UploadImage{
			FileName: name,
			Content:  content,
		}, "", "  ")
		if err != nil {
			fmt.Println("ошибка при преобразовании структуры в JSON (uploadimage):", err)
			return err
		}

		_, err = rest.CreateRequest("POST", restTypes.ItemUrl+"/"+itemId+"/images", authHeader, bytes.NewBuffer(jsonBody))
		if err != nil {
			return err
		}
		//fmt.Printf("uploadImage: %d\n Body: %s", response.StatusCode, response.Body)
	}
	return nil
}

func (me *merlionEntity) updateItem(SkladCredentials string, id string) {
	/*upateItem := restTypes.UpdateItem {

	}
	jsonBody, err := json.MarshalIndent(group, "", "  ")
	if err != nil {
		fmt.Println("ошибка при преобразовании структуры в JSON (updateitem):", err)
		return restTypes.TestProductMeta{}, err
	}
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
	//url := "https://api.moysklad.ru/api/remap/1.2/entity/productfolder"
	body, err := rest.CreateRequest("PUT", restTypes.ItemUrl, authHeader, bytes.NewBuffer(jsonBody))
	if err != nil {
		return restTypes.TestProductMeta{}, err
	}
	var meta restTypes.TestProductMeta
	fmt.Println("rawBody:", string(body))
	err = json.Unmarshal(body, &meta)
	if err != nil {
		return restTypes.TestProductMeta{}, err
	}
	return meta, nil*/
}

func (me *merlionEntity) getItemUUID(SkladCredentials string) (string, error) {
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	response, err := rest.CreateRequest("GET", restTypes.ItemUrl+"?search="+me.Vendor_part, authHeader, nil)
	if err != nil {
		return "", err
	}
	if response.StatusCode != 200 {
		return "", err
	}
	items := restTypes.SearchItem{}
	if err := json.Unmarshal(response.Body, &items); err != nil {
		return "", fmt.Errorf("ошибка при декодировании item (getItemUUID): %s", err.Error())
	}
	i, exist := matchSubstring(items.Rows, me.Vendor_part)
	if !exist {
		return "", nil
	}
	return items.Rows[i].Id, nil
}

func incrementID() string {
	newNumPart := counter
	counter++
	newID := fmt.Sprintf("I%05d", newNumPart)
	return newID
}

func (me *merlionEntity) setRemains(SkladCredentials string) error {
	if me.AvailableClient_MSK == 0 {
		return nil
	}
	orgMeta, _ := GetOrganizationMeta(SkladCredentials, "sandbox1250")
	storeMeta, _ := GetStoreMeta(SkladCredentials, "Мерлион")
	/*uuidItem, _ := me.getItemUUID(SkladCredentials)
	itemMeta := restTypes.MetaPositionsAdd{
		Id:       uuidItem,
		Quantity: strconv.Itoa(me.AvailableClient_MSK),
	}*/
	itemMeta, _ := GetItemMeta(SkladCredentials, me.Article)
	acceptance := restTypes.Acceptance{
		Organization: restTypes.MetaMiddle{Meta: orgMeta},
		Agent:        restTypes.MetaMiddle{Meta: orgMeta},
		Store:        restTypes.MetaMiddle{Meta: storeMeta},
		Description:  "Automatically generated by the script",
		Applicable:   true,
		Positions: []restTypes.PositionsAdd{{
			Quantity:   me.AvailableClient_MSK,
			Assortment: restTypes.MetaMiddle{Meta: itemMeta},
		}},
	}

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	jsonBody, err := json.MarshalIndent(acceptance, "", "  ")
	if err != nil {
		fmt.Println("ошибка при преобразовании структуры в JSON (setremains):", err)
		return err
	}
	_, err = rest.CreateRequest("POST", restTypes.AcceptanceUrl, authHeader, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	//fmt.Printf("setRemains: %d\n Body: %s", response.StatusCode, response.Body)
	return nil
}
