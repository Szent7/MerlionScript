package src

import (
	"MerlionScript/services/common"
	skladTypes "MerlionScript/services/sklad/types"
	"MerlionScript/utils/db"
	"MerlionScript/utils/rest"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"
)

// Проверка на существование картинок
// Вытягивание картинок из сервиса
// Заливка картинок в ERP и обновление записи
func UploadAllImages(ctx context.Context, dbInstance *db.DB, service common.Service, erpSystem common.ERPSystem) error {
	var ServiceName string = service.GetSystemName()
	var DBTableName string = service.GetDBTableName()
	select {
	case <-ctx.Done():
		fmt.Printf("%s (UploadAllImages): работу закончил из-за контекста\n", ServiceName)
		return nil
	default:
		fmt.Printf("%s: начал обновлять изображения\n", ServiceName)
		//Записи из БД
		items, err := dbInstance.GetCodesFilledMSNoImage(DBTableName)
		if err != nil {
			log.Printf("%s (UploadAllImages): ошибка при получении записей из БД | err = %s\n", ServiceName, err)
			return err
		}

		totalUploadedImages := 0
		for _, item := range *items {
			select {
			case <-ctx.Done():
				return nil
			default:
				//ID товара ERP
				erpItemID, err := erpSystem.GetItemID(item.Codes.MoySkladCode)
				if err != nil {
					log.Printf("%s (UploadAllImages): ошибка при получении ID товара из ERP | err = %s\n", ServiceName, err)
					continue
				}

				//Проверка на существование изображений в ERP
				erpImages, err := erpSystem.GetImagesList(erpItemID)
				if err != nil {
					log.Printf("%s (UploadAllImages): ошибка при получении списка изображений из ERP | err = %s\n", ServiceName, err)
					continue
				}
				if len(erpImages.Rows) == 10 { //Если в ERP уже есть 10 картинок
					item.Service.TryUploadImage = 1
					if err = dbInstance.UpdateService(item.Service, DBTableName); err != nil {
						log.Printf("%s (UploadAllImages): ошибка при изменении записи в БД | serviceCode = %s: %s\n", ServiceName, item.Service, err)
					}
					continue
				}
				//Если изображений нет на МС, то заливаем новые из сервиса
				//Вытягиваем изображения из сервиса
				serviceImages, err := service.GetImagesList(item.Service.ServiceCode)
				if err != nil {
					log.Printf("%s (UploadAllImages): ошибка при получении списка изображений из сервиса | article = %s; err = %s\n", ServiceName, item.Codes.Manufacturer, err)
					continue
				}
				uploadedImages := 0
				for i := range *serviceImages {
					//Загружаем изображение из сервиса
					response, _, err := rest.CreateRequestImageHeader("GET", (*serviceImages)[i].DownloadUrl, nil, "")
					if err != nil || response.StatusCode != 200 {
						log.Printf("%s (UploadAllImages): ошибка при получении изображения | url = %s; err = %s\n", ServiceName, (*serviceImages)[i].DownloadUrl, err)
						continue
					}
					newImage := skladTypes.UploadImage{
						Filename: (*serviceImages)[i].Filename,
						Content:  base64.StdEncoding.EncodeToString(response.Body),
					}
					//Загружаем изображение в ERP
					err = erpSystem.UploadImage(newImage, erpItemID)
					if err != nil {
						log.Printf("%s (UploadAllImages): ошибка при загрузке изображения в ERP | serviceCode = %s; err = %s\n", ServiceName, item.Service, err)
						continue
					}
					uploadedImages++
				}
				totalUploadedImages += uploadedImages
				//Если загружено хотя бы одно изображение или их нет в сервисе
				if uploadedImages > 0 || len(*serviceImages) == 0 {
					item.Service.TryUploadImage = 1
					if err = dbInstance.UpdateService(item.Service, DBTableName); err != nil {
						log.Printf("%s (UploadAllImages): ошибка при изменении записи в БД | serviceCode = %s: %s\n", ServiceName, item.Service, err)
						continue
					}
				}
				time.Sleep(time.Millisecond * 150)
			}
		}
		fmt.Printf("%s: добавлено %d изображений\n", ServiceName, totalUploadedImages)
		fmt.Printf("%s: закончил обновлять изображения\n", ServiceName)
		return nil
	}
}
