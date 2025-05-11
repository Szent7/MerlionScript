package merlion

import (
	"MerlionScript/cache"
	"MerlionScript/keeper"
	merlionReq "MerlionScript/services/merlion/requests"
	"MerlionScript/services/sklad"
	skladReq "MerlionScript/services/sklad/requests"
	skladTypes "MerlionScript/types/restTypes/sklad"
	merlionTypes "MerlionScript/types/soapTypes/merlion"
	"MerlionScript/utils/db"
	"MerlionScript/utils/db/typesDB"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"time"
)

const (
	//maxThreads = 10
	batchSize = 500
)

// Проверяем Мерлион на наличие новых позиций и добавляем их в БД
func CheckMerlionNewPositions(ctx context.Context) error {
	select {
	case <-ctx.Done():
		fmt.Println("checkMerlionNewPositions работу закончил из-за контекста")
		return nil
	default:
		fmt.Println("Начал проверку новых позиций на Мерлионе")
		catID, err := merlionReq.GetAllCatalogCodes()
		if err != nil {
			return err
		}

		if len(catID) == 0 {
			return fmt.Errorf("ошибка при получении номеров каталога: len(catID) = 0")
		}

		dbInstance, err := db.GetDBInstance()
		if err != nil {
			log.Printf("Ошибка при получении экземпляра БД (checkMerlionNewPositions): %s\n", err)
			return err
		}

		var countRecords int32 = 0
		//var idW = 0
		//var wg sync.WaitGroup
		//sem := make(chan struct{}, maxThreads)
		//re := regexp.MustCompile(`(?i)dahua`)

		//цикл по каталогам
		for i, id := range catID {
			fmt.Printf("Обработка каталога %s %d/%d\n", id, i+1, len(catID))
			//wg.Add(1)
			//go func(id string) {
			//defer wg.Done()
			select {
			case <-ctx.Done():
				//wg.Done()
				return nil
			default:
				//sem <- struct{}{}
				items, err := merlionReq.GetItemsByCatId(id)
				if err != nil {
					log.Printf("Ошибка при получении товаров по каталогу (checkMerlionNewPositions) id = %s: %s\n", id, err)
					//<-sem
				}
				//цикл по товарам из каталога
				for _, item := range items {
					lower := strings.ToLower(item.Brand)
					//if re.MatchString(item.Brand) {
					if strings.Contains(lower, "dahua") || strings.Contains(lower, "tenda") {
						newRecord := typesDB.Codes{
							MoySklad:         "",
							Manufacturer:     item.Vendor_part,
							ManufacturerName: item.Brand,
							Service:          item.No,
							MsOwnId:          0,
						}
						record, exists, err := dbInstance.GetCodeRecordByManufacturerAll(item.Vendor_part)
						if err != nil {
							log.Printf("Ошибка при получении записи из БД (checkMerlionNewPositions) manufacturer = %s: %s\n", item.Vendor_part, err)
							continue
						}
						if exists {
							newRecord.MoySklad = record.MoySklad
							newRecord.MsOwnId = record.MsOwnId
						}
						//Добавляем каждую позицию, т.к. в БД могут быть только уникальные записи
						added, err := dbInstance.AddCodeRecord(&newRecord, typesDB.MerlionTable)
						if err != nil {
							log.Printf("Ошибка при добавлении новой записи из Мерлиона (checkMerlionNewPositions) manufacturer = %s: %s\n", item.Vendor_part, err)
							continue
						}
						if added {
							atomic.AddInt32(&countRecords, 1)
						}
					}
				}
				//wg.Done()
				//<-sem
			}
			//}(id)
			//idW++
			time.Sleep(time.Millisecond * 150)
		}
		fmt.Println("Закончил проверку новых позиций на Мерлионе")

		//wg.Wait()
		//fmt.Printf("Всего горутин запущено: %d\n", idW)
		fmt.Printf("(%s) Добавлено %d записей при проверке позиций в Мерлионе\n", sklad.GetTimeNow(), countRecords)
		return nil
	}
}

// Создаем отсутствующие позиции на МС, добавляем в БД существующие
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

		items, err := dbInstance.GetCodeRecordsNoMS(typesDB.MerlionTable)
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
		//Остатки с Мерлиона
		itemsManufacturer := make([]string, len(*items))
		for i := range *items {
			itemsManufacturer[i] = (*items)[i].Service
		}
		MerlionItemsRaw := make([]merlionTypes.ItemCatalog, 0, len(*items))
		for i := 0; i < len(*items); i += batchSize {
			end := min(i+batchSize, len(*items))
			itemPart, err := merlionReq.GetItemsByItemIdBatch(itemsManufacturer[i:end])
			if err != nil || len(*itemPart) == 0 {
				log.Printf("Ошибка при получении остатков с Мерлиона (updateRemainsMS): %s\n", err)
				return err
			}
			MerlionItemsRaw = append(MerlionItemsRaw, *itemPart...)
		}
		MerlionItems := make(map[string]merlionTypes.ItemCatalog, len(MerlionItemsRaw))
		for i := range MerlionItemsRaw {
			if MerlionItemsRaw[i].No != "" {
				MerlionItems[MerlionItemsRaw[i].No] = MerlionItemsRaw[i]
			}
		}
		var createdItems int = 0
		for _, item := range *items {
			if createdItems >= 25 {
				break
			}
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
				itemCatalog, exist := MerlionItems[item.Service]
				if !exist {
					log.Printf("Ошибка при получении записи из Мерлиона (createNewPositionsMS) merlionCode = %s: %s\n", item.Service, err)
				}
				//Если этого товара не существует на МС
				if len(msItems.Rows) == 0 {
					if sklad.CreateNewItemMS(&item, itemCatalog.Name, catMeta, &counter, dbInstance, typesDB.MerlionTable) {
						createdItems++
					}
				} else { //Если найдены совпадения/похожие
					article := ""
					var founded bool = false
					var dhiProblem bool = false
					var emptyArticleProblem bool = false
					var sProblem bool = false
					var merlionSnum int = -10
					var msSnum int = -11
					var bProblem bool = false
					var merlionBnum int = -10
					var msBnum int = -11
					itemMerlion, err := merlionReq.GetItemsByItemIdBatch([]string{item.Service})
					if err != nil || len(*itemMerlion) == 0 || (*itemMerlion)[0].No == "" {
						log.Printf("Ошибка при получении записи из Мерлиона (createNewPositionsMS_2) merlionCode = %s: %s\n", item.Service, err)
						continue
					}
					// Вытягиваем -S0/-0000B номер из мерлиона, если есть
					substringsMer := strings.Fields((*itemMerlion)[0].Name)
					for _, subS := range substringsMer {
						merlionSnum = ExtractNumberFromS(subS)
						if merlionSnum >= 0 {
							break
						}
						merlionBnum = ExtractNumberFromB(subS)
						if merlionBnum >= 0 {
							break
						}
					}
					if merlionSnum == -1 {
						log.Printf("Ошибка при парсинге -S0 номера из Мерлиона (createNewPositionsMS_2) merlionCode = %s: %s\n", item.Service, err)
						continue
					}
					if merlionBnum == -1 {
						log.Printf("Ошибка при парсинге -0000B номера из Мерлиона (createNewPositionsMS_2) merlionCode = %s: %s\n", item.Service, err)
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
						msManufacturer := IgnoreDHManufacturer(msItems.Rows[i].Name)
						merManufacturer := IgnoreDHManufacturer(item.Manufacturer)
						if skladReq.СontainsSubstring(msManufacturer, merManufacturer) {
							dhiProblem = true
							break
						}
						// Если заканчивается на -S0
						if CheckSManufacturer(msItems.Rows[i].Name, item.Manufacturer) && merlionSnum >= 0 {
							// Вытягиваем -S0 номер из мс, если есть
							substringsMS := strings.Fields(msItems.Rows[i].Name)
							for _, subS := range substringsMS {
								msSnum = ExtractNumberFromS(subS)
								if msSnum >= 0 {
									break
								}
							}
							// Если номера совпали
							if msSnum == merlionSnum {
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
						if CheckBManufacturer(msItems.Rows[i].Name, item.Manufacturer) && merlionBnum >= 0 {
							// Вытягиваем -0000B номер из мс, если есть
							substringsMS := strings.Fields(msItems.Rows[i].Name)
							for _, subS := range substringsMS {
								msBnum = ExtractNumberFromB(subS)
								if msBnum >= 0 {
									break
								}
							}
							// Если номера совпали
							if msBnum == merlionBnum {
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
						if sklad.CreateNewItemMS(&item, itemCatalog.Name, catMeta, &counter, dbInstance, typesDB.MerlionTable) {
							createdItems++
						}
						continue
					}
					_, exists, err := dbInstance.GetCodeRecordByMS(article, typesDB.MerlionTable)
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
					ownId, _ := ExtractCounterFromOwnID(item.MoySklad)
					if ownId > 0 {
						item.MsOwnId = ownId
					}

					if err = dbInstance.EditCodeRecord(&item, typesDB.MerlionTable); err != nil {
						log.Printf("Ошибка при изменении записи в БД (createNewPositionsMS) merlionCode = %s: %s\n", item.Service, err)
					}
				}
			}
			time.Sleep(time.Millisecond * 100)
		}
		fmt.Println("Закончил сверять новые позиции с мс")
		fmt.Printf("(%s) Добавлено %d записей на МС\n", sklad.GetTimeNow(), createdItems)
		return nil
	}
}

// Обновляем остатки
func UpdateRemainsMS(ctx context.Context) error {
	fmt.Println("Начал обновлять остатки на мс")
	select {
	case <-ctx.Done():
		fmt.Println("updateRemainsMS работу закончил из-за контекста")
		return nil
	default:
		//Экземпляр БД
		dbInstance, err := db.GetDBInstance()
		if err != nil {
			log.Printf("Ошибка при получении экземпляра БД (createNewPositionsMS): %s\n", err)
			return err
		}
		//Записи из БД
		items, err := dbInstance.GetCodeRecordsFilledMS(typesDB.MerlionTable)
		if err != nil {
			log.Printf("Ошибка при получении записей из БД (createNewPositionsMS): %s\n", err)
			return err
		}
		//Мета организации МС
		orgMeta, err := skladReq.GetOrganizationMeta(keeper.K.GetMerlionOrg())
		if err != nil {
			log.Printf("Ошибка при получении метаданных организации (updateRemainsMS): %s\n", err)
			return err
		}
		//Мета склада МС
		storeMeta, err := skladReq.GetStoreMeta(keeper.K.GetMerlionSklad())
		if err != nil || storeMeta.Href == "" {
			log.Printf("Ошибка при получении метаданных склада (updateRemainsMS): %s\n", err)
			return err
		}
		//ID склада МС
		storeUUID, err := skladReq.GetStoreUUID(keeper.K.GetMerlionSklad())
		if err != nil || storeUUID == "" {
			log.Printf("Ошибка при получении UUID склада (updateRemainsMS): %s\n", err)
			return err
		}
		//Остатки с Мерлиона
		itemsManufacturer := make([]string, len(*items))
		for i := range *items {
			itemsManufacturer[i] = (*items)[i].Service
		}
		MerlionItemsAvailRaw := make([]merlionTypes.ItemAvail, 0, len(*items))
		for i := 0; i < len(*items); i += batchSize {
			end := min(i+batchSize, len(*items))
			availPart, err := merlionReq.GetItemsAvailByItemIdBatch(itemsManufacturer[i:end])
			if err != nil || len(*availPart) == 0 {
				log.Printf("Ошибка при получении остатков с Мерлиона (updateRemainsMS): %s\n", err)
				return err
			}
			MerlionItemsAvailRaw = append(MerlionItemsAvailRaw, *availPart...)
		}
		MerlionItemsAvail := make(map[string]merlionTypes.ItemAvailPrice, len(MerlionItemsAvailRaw))
		for i := range MerlionItemsAvailRaw {
			if MerlionItemsAvailRaw[i].No != "" {
				MerlionItemsAvail[MerlionItemsAvailRaw[i].No] = merlionTypes.ItemAvailPrice{
					AvailableClient_MSK: MerlionItemsAvailRaw[i].AvailableClient_MSK,
					PriceClientRUB_MSK:  MerlionItemsAvailRaw[i].PriceClientRUB_MSK,
				}
			}
		}
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
		for _, item := range *items {
			select {
			case <-ctx.Done():
				return nil
			default:
				//Данные МС
				var itemMS skladTypes.Rows
				itemMSRaw, err := cache.Cache.Get(item.MoySklad)
				//! нужно отрефакторить
				/*if err != nil {
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
				}*/

				if err != nil {
					itemMS, err = GetItemAndCache(item.MoySklad)
					if err != nil {
						log.Printf("Ошибка при получении товара МС (updateRemainsMS) msCode = %s: %s\n", item.MoySklad, err)
						continue
					}
				} else {
					itemMS, err = cache.Deserialize[skladTypes.Rows](itemMSRaw)
					if err != nil {
						cache.Cache.Delete(item.MoySklad)
						itemMS, err = GetItemAndCache(item.MoySklad)
						if err != nil {
							log.Printf("Ошибка при получении товара МС (updateRemainsMS) msCode = %s: %s\n", item.MoySklad, err)
							continue
						}
					}
				}

				itemRemainsMerlion, exists := MerlionItemsAvail[item.Service]
				if !exists {
					log.Printf("Ошибка при получении остатков с мерлиона (updateRemainsMS) merlionCode = %s\n", item.Service)
					continue
				}

				if itemMS.IsSerialTrackable {
					log.Printf("Товар с серийным учетом в приемку/списание не попал (updateRemainsMS) msCode = %s isSerialTrackable = %t\n", item.MoySklad, itemMS.IsSerialTrackable)
					continue
				}
				itemRemainsMS, err := skladReq.GetItemsAvail(itemMS.Id, storeUUID)
				if err != nil || itemRemainsMS < 0 {
					log.Printf("Ошибка при получении остатков с мс (updateRemainsMS) msCode = %s: %s\n", item.MoySklad, err)
					continue
				}
				//Данные Мерлиона
				/*itemRemainsMerlion, err := merlion.GetItemsAvailByItemId(item.Merlion)
				if err != nil || len(itemRemainsMerlion) == 0 || itemRemainsMerlion[0].No == "" {
					log.Printf("Ошибка при получении остатков с мерлиона (updateRemainsMS) merlionCode = %s: %s\n", item.Merlion, err)
					continue
				}*/
				if itemRemainsMS > itemRemainsMerlion.AvailableClient_MSK {
					woffList = append(woffList, skladTypes.PositionsAdd{
						Quantity:   itemRemainsMS - itemRemainsMerlion.AvailableClient_MSK,
						Assortment: skladTypes.MetaMiddle{Meta: itemMS.Meta},
					})
				} else if itemRemainsMS < itemRemainsMerlion.AvailableClient_MSK {
					acceptanceList = append(acceptanceList, skladTypes.PositionsAdd{
						Quantity:   itemRemainsMerlion.AvailableClient_MSK - itemRemainsMS,
						Assortment: skladTypes.MetaMiddle{Meta: itemMS.Meta},
						Price:      itemRemainsMerlion.PriceClientRUB_MSK * 100,
					})
				}
			}
			time.Sleep(time.Millisecond * 250)
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
