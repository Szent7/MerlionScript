package softtronik

import (
	"MerlionScript/cache"
	"MerlionScript/keeper"
	"MerlionScript/services/merlion"
	"MerlionScript/services/sklad"
	skladReq "MerlionScript/services/sklad/requests"
	softtronikReq "MerlionScript/services/softtronik/requests"
	skladTypes "MerlionScript/types/restTypes/sklad"
	"MerlionScript/types/restTypes/softtronik"
	"MerlionScript/utils/db"
	"MerlionScript/utils/db/typesDB"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
	"sync/atomic"
	"time"
)

var itemsGlobal []softtronik.ProductItem

func CheckSofttronikNewPositions(ctx context.Context) error {
	select {
	case <-ctx.Done():
		fmt.Println("CheckSofttronikNewPositions работу закончил из-за контекста")
		return nil
	default:
		fmt.Println("Начал проверку новых позиций на Софт-тронике")
		catID, err := softtronikReq.GetAllCategoryCodes()
		if err != nil {
			return err
		}
		if len(catID) == 0 {
			return fmt.Errorf("ошибка при получении номеров каталога: len(catID) = 0")
		}
		dbInstance, err := db.GetDBInstance()
		if err != nil {
			log.Printf("Ошибка при получении экземпляра БД (CheckSofttronikNewPositions): %s\n", err)
			return err
		}
		var countRecords int32 = 0
		itemsGlobal = make([]softtronik.ProductItem, 0, 200)
		//цикл по каталогам
		for i, cat := range catID {
			fmt.Printf("Обработка каталога %s %d/%d\n", cat.Name, i+1, len(catID))
			select {
			case <-ctx.Done():
				return nil
			default:
				items, err := softtronikReq.GetItemsByCatId(cat.ID)
				if err != nil {
					log.Printf("Ошибка при получении товаров по каталогу (CheckSofttronikNewPositions) id = %s: %s\n", cat.Name, err)
					continue
				}
				itemsGlobal = append(itemsGlobal, items...)
				//цикл по товарам из каталога
				for _, item := range items {
					//Добавляем каждую позицию, т.к. в БД могут быть только уникальные записи
					newRecord := typesDB.Codes{
						MoySklad:         "",
						Manufacturer:     item.Article,
						ManufacturerName: "dahua",
						Service:          item.Code,
						MsOwnId:          0,
					}
					record, exists, err := dbInstance.GetCodeRecordByManufacturerAll(item.Article)
					if err != nil {
						log.Printf("Ошибка при получении записи из БД (CheckSofttronikNewPositions) manufacturer = %s: %s\n", item.Article, err)
						continue
					}
					if exists {
						newRecord.MoySklad = record.MoySklad
						newRecord.MsOwnId = record.MsOwnId
					}
					added, err := dbInstance.AddCodeRecord(&newRecord, typesDB.SofttronikTable)
					if err != nil {
						log.Printf("Ошибка при добавлении новой записи из Софт-троника (CheckSofttronikNewPositions) manufacturer = %s: %s\n", item.Article, err)
						continue
					}
					if added {
						atomic.AddInt32(&countRecords, 1)
					}
				}
			}
		}
		fmt.Println("Закончил проверку новых позиций в Софт-тронике")
		fmt.Printf("(%s) Добавлено %d записей при проверке позиций в Софт-тронике\n", sklad.GetTimeNow(), countRecords)
		return nil
	}
}

func CreateNewPositionsMS(ctx context.Context) error {
	fmt.Println("Начал сверять новые позиции с мс")
	select {
	case <-ctx.Done():
		fmt.Println("createNewPositionsMS работу закончил из-за контекста")
		return nil
	default:
		dbInstance, err := db.GetDBInstance()
		if err != nil {
			log.Printf("Ошибка при получении экземпляра БД (createNewPositionsMS): %s\n", err)
			return err
		}

		items, err := dbInstance.GetCodeRecordsNoMS(typesDB.SofttronikTable)
		if err != nil {
			log.Printf("Ошибка при получении записей из БД (createNewPositionsMS): %s\n", err)
			return err
		}
		//TODO Необходимо добавить парсинг последнего кода из базы мс, чтобы счетчик не задвоился
		counter, err := dbInstance.GetLastOwnIdMS(typesDB.OwnIDsTable)
		if err != nil {
			log.Printf("Ошибка при получении counter из БД (createNewPositionsMS): %s\n", err)
			return err
		}
		catMeta, err := skladReq.GetCatMeta(keeper.K.GetSkladCat())
		if err != nil || catMeta.Href == "" {
			log.Printf("Ошибка при получении метаданных склада (updateRemainsMS): %s\n", err)
			return err
		}

		var createdItems int = 0
		for _, item := range *items {
			/*if createdItems >= 50 {
				break
			}*/
			select {
			case <-ctx.Done():
				return nil
			default:
				manufacturerReplace := strings.Replace(item.Manufacturer, " ", "+", -1)
				response, err := skladReq.GetItemByManufacturer(manufacturerReplace)
				if err != nil || response.StatusCode != 200 {
					log.Printf("Ошибка при получении записи из МС (createNewPositionsMS) manufacturer = %s StatusCode = %d: %s, \n", item.Manufacturer, response.StatusCode, err)
					continue
				}
				msItems := skladTypes.SearchItem{}
				if err := json.Unmarshal(response.Body, &msItems); err != nil {
					log.Printf("Ошибка при декодировании item (createNewPositionsMS) manufacturer = %s: %s", item.Manufacturer, err)
					continue
				}
				itemCatalog, found := getGlobalItemsRecord(item.Service, itemsGlobal)
				if !found {
					log.Printf("Ошибка при получении записи из Софт-троника (createNewPositionsMS) softtronikCode = %s: %s\n", item.Service, err)
				}
				//Если этого товара не существует на МС
				if len(msItems.Rows) == 0 {
					if sklad.CreateNewItemMS(&item, itemCatalog.Name+" "+itemCatalog.Article, catMeta, &counter, dbInstance, typesDB.SofttronikTable) {
						createdItems++
					}
				} else { //Если найдены совпадения/похожие
					article := ""
					var founded bool = false
					var dhiProblem bool = false
					var emptyArticleProblem bool = false
					var sProblem bool = false
					var softtronikSnum int = -10
					var msSnum int = -11
					var bProblem bool = false
					var softtronikBnum int = -10
					var msBnum int = -11
					/*itemMerlion, err := merlionReq.GetItemsByItemIdBatch([]string{item.Service})
					if err != nil || len(*itemMerlion) == 0 || (*itemMerlion)[0].No == "" {
						log.Printf("Ошибка при получении записи из Софт-троника (createNewPositionsMS_2) softtronikCode = %s: %s\n", item.Service, err)
						continue
					}*/
					// Вытягиваем -S0/-0000B номер из мерлиона, если есть
					substringssofttronik := strings.Fields(itemCatalog.Name)
					for _, subS := range substringssofttronik {
						softtronikSnum = merlion.ExtractNumberFromS(subS)
						if softtronikSnum >= 0 {
							break
						}
						softtronikBnum = merlion.ExtractNumberFromB(subS)
						if softtronikBnum >= 0 {
							break
						}
					}
					if softtronikSnum == -1 {
						log.Printf("Ошибка при парсинге -S0 номера из Софт-троника (createNewPositionsMS_2) softtronikCode = %s: %s\n", item.Service, err)
						continue
					}
					if softtronikBnum == -1 {
						log.Printf("Ошибка при парсинге -0000B номера из Софт-троника (createNewPositionsMS_2) softtronikCode = %s: %s\n", item.Service, err)
						continue
					}
					//TODO позиция с 3 артикулами, добавить проверку на уникальность кода мс
					for i := range msItems.Rows {
						// Если полностью совпадает
						if skladReq.СontainsSubstring(msItems.Rows[i].Name, item.Manufacturer) {
							founded = true
							// Но при этом артикул пустой
							if msItems.Rows[i].Article == "" {
								emptyArticleProblem = true
								break
							}
							// Не пустой
							article = msItems.Rows[i].Article
							break
						}
						// Если начинается на DH-/DHI-, то убираем их из сравнения
						msManufacturer := merlion.IgnoreDHManufacturer(msItems.Rows[i].Name)
						merManufacturer := merlion.IgnoreDHManufacturer(item.Manufacturer)
						if skladReq.СontainsSubstring(msManufacturer, merManufacturer) {
							dhiProblem = true
							break
						}
						// Если заканчивается на -S0
						if merlion.CheckSManufacturer(msItems.Rows[i].Name, item.Manufacturer) && softtronikSnum >= 0 {
							// Вытягиваем -S0 номер из мс, если есть
							substringsMS := strings.Fields(msItems.Rows[i].Name)
							for _, subS := range substringsMS {
								msSnum = merlion.ExtractNumberFromS(subS)
								if msSnum >= 0 {
									break
								}
							}
							// Если номера совпали
							if msSnum == softtronikSnum {
								sProblem = true
								founded = true
								// Но при этом артикул пустой
								if msItems.Rows[i].Article == "" {
									emptyArticleProblem = true
									break
								}
								// Не пустой
								article = msItems.Rows[i].Article
								break
							}
						}

						// Если заканчивается на -0000B
						if merlion.CheckBManufacturer(msItems.Rows[i].Name, item.Manufacturer) && softtronikBnum >= 0 {
							// Вытягиваем -0000B номер из мс, если есть
							substringsMS := strings.Fields(msItems.Rows[i].Name)
							for _, subS := range substringsMS {
								msBnum = merlion.ExtractNumberFromB(subS)
								if msBnum >= 0 {
									break
								}
							}
							// Если номера совпали
							if msBnum == softtronikBnum {
								bProblem = true
								founded = true
								// Но при этом артикул пустой
								if msItems.Rows[i].Article == "" {
									emptyArticleProblem = true
									break
								}
								// Не пустой
								article = msItems.Rows[i].Article
								break
							}
						}
					}
					if msSnum == -1 {
						log.Printf("Ошибка при парсинге -S0 номера из мс manufacturer = %s\n", item.Manufacturer)
						continue
					}
					if msBnum == -1 {
						log.Printf("Ошибка при парсинге -0000B номера из мс manufacturer = %s\n", item.Manufacturer)
						continue
					}
					if dhiProblem {
						log.Printf("Проблема DH-/DHI- manufacturer = %s", item.Manufacturer)
						continue
					}
					if emptyArticleProblem {
						log.Printf("Проблема пустого артикула manufacturer = %s", item.Manufacturer)
						continue
					}
					if sProblem {
						log.Printf("Проблема окончания -S0 manufacturer = %s | соответствие на мс = %s\n", item.Manufacturer, article)
					}
					if bProblem {
						log.Printf("Проблема окончания -0000B manufacturer = %s | соответствие на мс = %s\n", item.Manufacturer, article)
					}
					if !founded {
						log.Printf("Полных соответствий не найдено manufacturer = %s (создаю новую позицию)", item.Manufacturer)
						if sklad.CreateNewItemMS(&item, itemCatalog.Name+" "+itemCatalog.Article, catMeta, &counter, dbInstance, typesDB.SofttronikTable) {
							createdItems++
						}
						continue
					}
					_, exists, err := dbInstance.GetCodeRecordByMS(article, typesDB.SofttronikTable)
					if err != nil {
						log.Printf("Ошибка при получении записи из БД (createNewPositionsMS_2): %s\n", err)
						continue
					}
					if exists {
						log.Printf("Вероятное задвоение manufacturer = %s | соответствие на мс = %s\n", item.Manufacturer, article)
						continue
					}
					/*if article == "" {
						log.Printf("Соответствий не найдено или поле артикулов пустое manufacturer = %s: ", item.Manufacturer)
						continue
					}*/
					item.MoySklad = article
					//! опасный блок, может скопировать код чужой записи
					/*
						if article == "" {
							copyItem.MoySklad = msItems.Rows[0].Article
						} else {
							copyItem.MoySklad = article
						}*/
					// Извлекаем counter из артикула
					ownId, _ := merlion.ExtractCounterFromOwnID(item.MoySklad)
					if ownId > 0 {
						item.MsOwnId = ownId
					}

					if err = dbInstance.EditCodeRecord(&item, typesDB.SofttronikTable); err != nil {
						log.Printf("Ошибка при изменении записи в БД (createNewPositionsMS) softtronikCode = %s: %s\n", item.Service, err)
					}
				}
			}
			time.Sleep(time.Millisecond * 150)
		}
		fmt.Println("Закончил сверять новые позиции с мс")
		fmt.Printf("(%s) Добавлено %d записей на МС\n", sklad.GetTimeNow(), createdItems)
		return nil
	}
}

func UpdateRemainsMS(ctx context.Context) error {
	//fillItemsGlobal(token)
	fmt.Println("Начал обновлять остатки на мс")
	select {
	case <-ctx.Done():
		fmt.Println("updateRemainsMS работу закончил из-за контекста")
		return nil
	default:
		//Экземпляр БД
		dbInstance, err := db.GetDBInstance()
		if err != nil {
			log.Printf("Ошибка при получении экземпляра БД (updateRemainsMS): %s\n", err)
			return err
		}
		//Записи из БД
		items, err := dbInstance.GetCodeRecordsFilledMS(typesDB.SofttronikTable)
		if err != nil {
			log.Printf("Ошибка при получении записей из БД (updateRemainsMS): %s\n", err)
			return err
		}
		//Записи из Софт-троника
		catID, err := softtronikReq.GetAllCategoryCodes()
		if err != nil {
			return err
		}
		if len(catID) == 0 {
			return fmt.Errorf("ошибка при получении номеров каталога: len(catID) = 0")
		}
		itemsSofttronik, err := softtronikReq.GetItemsAvailsAll(catID)
		if err != nil {
			log.Printf("Ошибка при получении записей из Софт-троника (updateRemainsMS): %s\n", err)
			return err
		}
		itemStockSofttronik := getAvailsItemsRecords(itemsSofttronik)
		//Мета организации МС
		orgMeta, err := skladReq.GetOrganizationMeta(keeper.K.GetSofttronikOrg())
		if err != nil {
			log.Printf("Ошибка при получении метаданных организации (updateRemainsMS): %s\n", err)
			return err
		}
		//Мета склада МС
		storeMeta, err := skladReq.GetStoreMeta(keeper.K.GetSofttronikSklad())
		if err != nil || storeMeta.Href == "" {
			log.Printf("Ошибка при получении метаданных склада (updateRemainsMS): %s\n", err)
			return err
		}
		//ID склада МС
		storeUUID, err := skladReq.GetStoreUUID(keeper.K.GetSofttronikSklad())
		if err != nil || storeUUID == "" {
			log.Printf("Ошибка при получении UUID склада (updateRemainsMS): %s\n", err)
			return err
		}
		//Курс доллара
		/*var currency float64 = 0
		for i := range itemsSofttronik.Body.Currency {
			if itemsSofttronik.Body.Currency[i].ID == "USD" {
				currency, err = strconv.ParseFloat(itemsSofttronik.Body.Currency[i].Rate, 64)
				if err != nil {
					log.Printf("Ошибка при получении курса валют Софт-троник (updateRemainsMS): %s\n", err)
					return err
				}
				break
			}
		}
		if currency == 0 {
			log.Printf("Ошибка при получении курса валют Софт-троник (updateRemainsMS): %s\n", err)
			return err
		}*/
		//На увеличение остатков
		acceptanceReq := skladTypes.Acceptance{
			Organization: skladTypes.MetaMiddle{Meta: orgMeta},
			Agent:        skladTypes.MetaMiddle{Meta: orgMeta},
			Store:        skladTypes.MetaMiddle{Meta: storeMeta},
			Description:  "Automatically generated by the script",
			Applicable:   true,
		}
		//На снижение остатков
		woffReq := skladTypes.Acceptance{
			Organization: skladTypes.MetaMiddle{Meta: orgMeta},
			Agent:        skladTypes.MetaMiddle{Meta: orgMeta},
			Store:        skladTypes.MetaMiddle{Meta: storeMeta},
			Description:  "Automatically generated by the script",
			Applicable:   true,
		}
		//Список приемки
		acceptanceList := []skladTypes.PositionsAdd{}
		//Список списания
		woffList := []skladTypes.PositionsAdd{}
		//addition := 0
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

				itemRemainsSofttronik, ok := itemStockSofttronik[item.Manufacturer]
				if !ok {
					//itemRemainsSofttronik, err = softtronikReq.GetItemsByItemIdFormatted(item.Service, token)
					//if err != nil {
					log.Printf("Ошибка при получении остатков с Софт-троника (updateRemainsMS) softtronikCode = %s\n", item.Service)
					continue
					//}
				}
				if itemMS.IsSerialTrackable {
					log.Printf("Товар с серийным учетом в приемку/списание не попал (updateRemainsMS) msCode = %s isSerialTrackable = %t\n", item.MoySklad, itemMS.IsSerialTrackable)
					continue
				}
				itemRemainsMS, err := skladReq.GetItemsAvail(itemMS.Id, storeUUID)
				if err != nil {
					log.Printf("Ошибка при получении остатков с мс (updateRemainsMS) msCode = %s: %s\n", item.MoySklad, err)
					continue
				}

				//Если остаток меньше 0, то делаем приемку с учетом этого
				/*addition = 0
				if itemRemainsMS < 0 {
					addition = -itemRemainsMS
					itemRemainsMS = 0
				}*/
				if itemRemainsMS > itemRemainsSofttronik.Stocks {
					woffList = append(woffList, skladTypes.PositionsAdd{
						Quantity:   itemRemainsMS - itemRemainsSofttronik.Stocks,
						Assortment: skladTypes.MetaMiddle{Meta: itemMS.Meta},
					})
					//fmt.Printf("%s - %s\n", item.Service, item.MoySklad)
				} else if itemRemainsMS < itemRemainsSofttronik.Stocks {
					rub_price := itemRemainsSofttronik.Price
					half_rub_price := int(math.Ceil(rub_price * 100))
					acceptanceList = append(acceptanceList, skladTypes.PositionsAdd{
						Quantity:   itemRemainsSofttronik.Stocks - itemRemainsMS, //+ addition
						Assortment: skladTypes.MetaMiddle{Meta: itemMS.Meta},
						Price:      float32(half_rub_price),
					})
				}
			}
			time.Sleep(time.Millisecond * 150)
		}
		if len(acceptanceList) != 0 {
			acceptanceReq.Positions = acceptanceList
			fmt.Printf("Позиций в приемке (updateRemainsMS): %d\n", len(acceptanceList))
			err = skladReq.IncreaseItemsAvail(&acceptanceReq)
			if err != nil {
				log.Printf("Ошибка при создании приемки (updateRemainsMS): %s\n", err)
			}
		}
		if len(woffList) != 0 {
			woffReq.Positions = woffList
			fmt.Printf("Позиций в списании (updateRemainsMS): %d\n", len(woffList))
			err = skladReq.DecreaseItemsAvail(&woffReq)
			if err != nil {
				log.Printf("Ошибка при создании списания (updateRemainsMS): %s\n", err)
			}
		}
		fmt.Println("Закончил обновлять остатки на мс")
		return nil
	}
}
