package src

import (
	"MerlionScript/services/common"
	"MerlionScript/utils/db/interfaceDB"
	"MerlionScript/utils/db/typesDB"
	"context"
	"fmt"
	"log"
	"strings"
)

// Заносим новые артикула в БД
func AddNewRecords(ctx context.Context, dbInstance interfaceDB.DB, service common.Service) error {
	var ServiceName string = service.GetSystemName()
	var DBTableName string = service.GetDBTableName()
	select {
	case <-ctx.Done():
		fmt.Printf("%s (AddNewRecords): работу закончил из-за контекста\n", ServiceName)
		return nil
	default:
		fmt.Printf("%s: начал проверку новых позиций\n", ServiceName)
		articleItems, err := service.GetArticlesList(ctx)
		if err != nil || articleItems == nil {
			log.Printf("%s (AddNewRecords): ошибка при получении артикулов | err = %s\n", ServiceName, err)
			return err
		}
		var countRecords int = 0
		//цикл по списку артикулов
		for i := range *articleItems {
			lowerBrand := strings.ToLower((*articleItems)[i].Brand)
			if strings.Contains(lowerBrand, "dahua") || strings.Contains(lowerBrand, "tenda") {
				newCodesIDs := typesDB.CodesIDs{
					MsOwnId:      0,
					MoySkladCode: "",
					Article:      (*articleItems)[i].Article,
					Manufacturer: lowerBrand,
				}
				newService := typesDB.CodesService{
					Article:        (*articleItems)[i].Article,
					ServiceCode:    (*articleItems)[i].ServiceCode,
					TryUploadImage: 0,
				}
				//Добавляем CodesIDs
				_, err := dbInstance.InsertCodesIDs(newCodesIDs)
				if err != nil {
					log.Printf("%s (AddNewRecords): ошибка при добавлении CodesIDs в БД | article = %s: %s\n", ServiceName, (*articleItems)[i].Article, err)
					continue
				}
				//Добавляем Service
				added, err := dbInstance.InsertService(newService, DBTableName)
				if err != nil {
					log.Printf("%s (AddNewRecords): ошибка при добавлении Service в БД | article = %s: %s\n", ServiceName, (*articleItems)[i].Article, err)
					continue
				}
				if added {
					countRecords++
				}
			}
		}
		fmt.Printf("%s: закончил проверку новых позиций\n", ServiceName)
		fmt.Printf("%s: добавлено %d записей при проверке позиций\n", ServiceName, countRecords)
	}
	return nil
}
