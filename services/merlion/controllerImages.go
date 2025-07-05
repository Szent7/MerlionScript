package merlion

import (
	"MerlionScript/cache"
	merlionReq "MerlionScript/services/merlion/requests"
	skladReq "MerlionScript/services/sklad/requests"
	skladTypes "MerlionScript/types/restTypes/sklad"
	merlionTypes "MerlionScript/types/soapTypes/merlion"
	"MerlionScript/utils/db"
	"MerlionScript/utils/db/typesDB"
	"MerlionScript/utils/rest"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Проверка на существование картинок
// Вытягивание картинок из мерлиона
// Заливка картинок на МС и обновление записи
func UploadAllImages(ctx context.Context) error {
	fmt.Println("Начал обновлять изображения на мс по мерлиону")
	select {
	case <-ctx.Done():
		fmt.Println("UploadAllImages работу закончил из-за контекста")
		return nil
	default:
		//Экземпляр БД
		dbInstance, err := db.GetDBInstance()
		if err != nil {
			log.Printf("Ошибка при получении экземпляра БД (UploadAllImages): %s\n", err)
			return err
		}
		//Записи из БД
		items, err := dbInstance.GetCodeRecordsFilledMSWithNoImage(typesDB.MerlionTable)
		if err != nil {
			log.Printf("Ошибка при получении записей из БД (UploadAllImages): %s\n", err)
			return err
		}
		totalUploadedImages := 0
		for _, item := range *items {
			select {
			case <-ctx.Done():
				return nil
			default:
				//Данные МС
				var itemMS skladTypes.Rows
				var err1 error
				itemMSRaw, err := cache.Cache.Get(item.MoySklad)
				//! нужно отрефакторить
				if err != nil {
					itemMS, err1 = skladReq.GetItem(item.MoySklad)
					if err1 != nil || itemMS.Id == "" {
						log.Printf("Ошибка при получении товара МС (updateRemainsMS) msCode = %s: %s\n", item.MoySklad, err)
						continue
					}
					if sItemMS, err1 := cache.Serialize(itemMS); err1 == nil {
						cache.Cache.Append(item.MoySklad, sItemMS)
					}
				} else {
					itemMS, err1 = cache.Deserialize[skladTypes.Rows](itemMSRaw)
					if err1 != nil {
						cache.Cache.Delete(item.MoySklad)
						itemMS, err1 = skladReq.GetItem(item.MoySklad)
						if err1 != nil || itemMS.Id == "" {
							log.Printf("Ошибка при получении товара МС (updateRemainsMS) msCode = %s: %s\n", item.MoySklad, err)
							continue
						}
						if sItemMS, err1 := cache.Serialize(itemMS); err1 == nil {
							cache.Cache.Append(item.MoySklad, sItemMS)
						}
					}
				}
				//Проверка на существование изображений на МС
				response, err := skladReq.GetItemsImagesData(itemMS.Id)
				if err != nil || response.StatusCode != 200 {
					log.Printf("Ошибка при получении записей из МС (UploadAllImages): %s\n", err)
					continue
				}
				msImages := skladTypes.SearchImage{}
				if err := json.Unmarshal(response.Body, &msImages); err != nil {
					log.Printf("Ошибка при декодировании msImages (UploadAllImages) manufacturer = %s: %s", item.Manufacturer, err)
					continue
				}
				//Если изображения есть на МС
				if len(msImages.Rows) != 0 {
					item.LoadedImage = 1
				} else if len(msImages.Rows) == 10 { //Если на МС уже есть 10 картинок
					item.LoadedImage = 1
					item.TryLoadImage = 1
					if err = dbInstance.EditCodeRecord(&item, typesDB.MerlionTable); err != nil {
						log.Printf("Ошибка при изменении записи в БД (createNewPositionsMS) merlionCode = %s: %s\n", item.Service, err)
					}
					continue
				}
				//Если изображений нет на МС, то заливаем новые из мерлиона
				//Вытягиваем изображения из мерлиона
				listImages, err := merlionReq.GetImagesByItemIdFormatted(item.Service)
				if err != nil {
					log.Printf("Ошибка при получении записей из мерлиона (UploadAllImages) manufacturer = %s: %s\n", item.Manufacturer, err)
					continue
				}
				uploadedImages := len(listImages)
				for i := range listImages {
					url := merlionTypes.DownloadImageUrl + "/" + listImages[i]
					response, _, err := rest.CreateRequestImageHeader("GET", url, nil, "")
					if err != nil || response.StatusCode != 200 {
						log.Printf("Ошибка при получении изображений из мерлиона (UploadAllImages) url = %s: %s\n", url, err)
						continue
					}
					newImage := skladTypes.UploadImage{
						FileName: listImages[i],
						Content:  base64.StdEncoding.EncodeToString(response.Body),
					}
					//Заливаем на МС
					totalUploadedImages++
					resp, err := skladReq.UploadImage(itemMS.Id, newImage)
					if err != nil || resp.StatusCode != 200 {
						log.Printf("Ошибка при загрузке изображения на МС (UploadAllImages) merlionCode = %s: %s\n", item.Service, err)
						uploadedImages--
						totalUploadedImages--
					}
				}
				//если загружено хотя бы одно изображение или их нет на мерлионе
				if uploadedImages > 0 || len(listImages) == 0 {
					item.TryLoadImage = 1
					if err = dbInstance.EditCodeRecord(&item, typesDB.MerlionTable); err != nil {
						log.Printf("Ошибка при изменении записи в БД (createNewPositionsMS) merlionCode = %s: %s\n", item.Service, err)
						continue
					}
				}
				time.Sleep(time.Millisecond * 150)
			}
		}
		log.Printf("Добавлено %d изображений из Мерлиона\n", totalUploadedImages)
		fmt.Println("Закончил обновлять изображения на мс по мерлиону")
		return nil
	}
}
