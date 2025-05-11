package netlab

import (
	"MerlionScript/cache"
	netlabReq "MerlionScript/services/netlab/requests"
	skladReq "MerlionScript/services/sklad/requests"
	skladTypes "MerlionScript/types/restTypes/sklad"
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
// Вытягивание картинок из нетлаба
// Заливка картинок на МС и обновление записи
func UploadAllImages(ctx context.Context, token string) error {
	//fillItemsGlobal(token)
	fmt.Println("Начал обновлять изображения на мс по нетлаб")
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
		items, err := dbInstance.GetCodeRecordsFilledMSWithNoImage(typesDB.NetlabTable)
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
					item.TryLoadImage = 1
					item.LoadedImage = 1
					if err = dbInstance.EditCodeRecord(&item, typesDB.NetlabTable); err != nil {
						log.Printf("Ошибка при изменении записи в БД (createNewPositionsMS) netlabCode = %s: %s\n", item.Service, err)
						continue
					}
				}
				//Если изображений нет на МС, то заливаем новые из нетлаба
				//Вытягиваем изображения из нетлаба
				listImages, err := netlabReq.GetImagesByItemIdFormatted(item.Service, token)
				if err != nil {
					log.Printf("Ошибка при получении записей из нетлаба (UploadAllImages) manufacturer = %s: %s\n", item.Manufacturer, err)
					continue
				}
				uploadedImages := len(listImages)
				for k, v := range listImages {
					response, contentType, err := rest.CreateRequestImageHeader("GET", v, nil, "")
					if err != nil || response.StatusCode != 200 {
						log.Printf("Ошибка при получении изображений из нетлаба (UploadAllImages) url = %s: %s\n", v, err)
						continue
					}
					newImage := skladTypes.UploadImage{
						FileName: k + getExtensionFromContentType(contentType),
						Content:  base64.StdEncoding.EncodeToString(response.Body),
					}
					//Заливаем на МС
					totalUploadedImages++
					resp, err := skladReq.UploadImage(itemMS.Id, newImage)
					if err != nil || resp.StatusCode != 200 {
						log.Printf("Ошибка при загрузке изображения в нетлаб (UploadAllImages) netlabCode = %s: %s\n", item.Service, err)
						uploadedImages--
						totalUploadedImages--
					}
				}
				//если загружено хотя бы одно изображение или их нет на нетлабе
				if uploadedImages > 0 || len(listImages) == 0 {
					item.TryLoadImage = 1
					item.LoadedImage = 0
					if err = dbInstance.EditCodeRecord(&item, typesDB.NetlabTable); err != nil {
						log.Printf("Ошибка при изменении записи в БД (createNewPositionsMS) netlabCode = %s: %s\n", item.Service, err)
						continue
					}
				}
				time.Sleep(time.Millisecond * 150)
			}
		}
		log.Printf("Добавлено %d изображений из Нетлаба\n", totalUploadedImages)
		fmt.Println("Закончил обновлять изображения на мс по нетлаб")
		return nil
	}
}
