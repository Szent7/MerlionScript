package src

import (
	"MerlionScript/services/common"
	"MerlionScript/services/sklad"
	"MerlionScript/utils/db"
	"MerlionScript/utils/db/typesDB"
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// Создаем отсутствующие позиции в ERP или создаем привязки
func CreateNewPositionsERP(ctx context.Context, dbInstance *db.DB, service common.Service, erpSystem common.ERPSystem) error {
	var ServiceName string = service.GetSystemName()
	var DBTableName string = service.GetDBTableName()
	catName := erpSystem.GetCatName()
	catMeta, err := erpSystem.GetCatMeta()
	if err != nil {
		log.Printf("%s (CreateNewPositionsERP): ошибка при получении метаданных каталога | err = %s\n", ServiceName, err)
		return err
	}
	Items, err := service.GetItemsList()
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		fmt.Printf("%s (CreateNewPositionsERP): работу закончил из-за контекста\n", ServiceName)
		return nil
	default:
		fmt.Printf("%s: начал сверять новые позиции с мс\n", ServiceName)
		//TODO Необходимо добавить парсинг последнего кода из базы мс, чтобы счетчик не задвоился
		counter, err := dbInstance.GetLastOwnIdMS()
		if err != nil {
			log.Printf("%s (CreateNewPositionsERP): ошибка при получении counter из БД | err = %s\n", ServiceName, err)
			return err
		}
		var createdItems int = 0
		for i := range *Items {
			select {
			case <-ctx.Done():
				return nil
			default:
				//Поиск по товарам МС
				articleReplace := strings.Replace((*Items)[i].Article, " ", "+", -1)
				erpItems, err := erpSystem.GetItemsByArticle(articleReplace)
				if err != nil {
					log.Printf("%s (CreateNewPositionsERP): ошибка при получении записи из ERP | article = %s; err = %s\n", ServiceName, (*Items)[i].Article, err)
					continue
				}
				//Запись из БД
				erpRecord, err := dbInstance.GetCodesIDs((*Items)[i].Article)
				if err != nil || erpRecord == (typesDB.CodesIDs{}) {
					log.Printf("%s (CreateNewPositionsERP): ошибка при получении записи из БД | article = %s; err = %s\n", ServiceName, (*Items)[i].Article, err)
					continue
				}
				//Если этого товара не существует на МС
				if len(*erpItems) == 0 {
					log.Printf("%s (CreateNewPositionsERP): полных соответствий не найдено (создаю новую позицию) | article = %s ", ServiceName, (*Items)[i].Article)
					if sklad.CreateNewItemMS(&erpRecord, catName, catMeta, &counter, dbInstance, DBTableName) {
						createdItems++
					}
				} else { //Если найдены совпадения/похожие
					//Проверяем существование позиции на МС
					foundedItem, err := CompareArticle(erpItems, &(*Items)[i])
					if err != nil {
						log.Printf("%s: %s\n", ServiceName, err.Error())
						continue
					}
					//Если позиция существует
					if foundedItem != (common.ItemList{}) {
						//foundedItem.Article - код МС, а не артикул позиции
						/*_, exists, err := dbInstance.GetCodeRecordByMS(foundedItem.Article, DBTableName)
						if err != nil {
							log.Printf("%s (CreateNewPositionsERP): ошибка при получении записи из БД | err = %s\n", ServiceName, err)
							continue
						}
						if exists {
							log.Printf("%s (CreateNewPositionsERP): вероятное задвоение article = %s | соответствие на мс = %s\n", ServiceName, (*Items)[i].Article, foundedItem.Article)
							continue
						}*/
						erpRecord.MoySklad = foundedItem.Article
						// Извлекаем counter из артикула
						ownId, _ := db.ExtractCounterFromOwnID(erpRecord.MoySklad)
						if ownId > 0 {
							erpRecord.MsOwnId = ownId
						}
						if err = dbInstance.UpdateCodesIDs(erpRecord); err != nil {
							log.Printf("%s (CreateNewPositionsERP): ошибка при изменении записи в БД | serviceCode = %s; err = %s\n", ServiceName, ServiceName, err)
						}
					} else { //Если совпадения не найдены
						log.Printf("%s (CreateNewPositionsERP): полных соответствий не найдено (создаю новую позицию) | article = %s ", ServiceName, (*Items)[i].Article)
						if sklad.CreateNewItemMS(&erpRecord, catName, catMeta, &counter, dbInstance, DBTableName) {
							createdItems++
						}
					}
				}
			}
			time.Sleep(time.Millisecond * 100)
		}
		fmt.Printf("%s: закончил сверять новые позиции с мс\n", ServiceName)
		fmt.Printf("%s: добавлено %d записей на МС\n", ServiceName, createdItems)
		return nil
	}
}

//!92 - нихуя не понятно
