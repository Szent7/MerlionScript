package controller

import (
	"MerlionScript/keeper"
	"MerlionScript/service/merlion"
	"MerlionScript/service/sklad"
	"MerlionScript/types/restTypes"
	"MerlionScript/utils/db"
	"MerlionScript/utils/db/typesDB"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

/*const (
	maxThreads = 10
)*/

func GetTimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func StartController(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	checkMerlionNewPositions(ctx)
	createNewPositionsMS(ctx)
	updateRemainsMS(ctx)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Контроллер завершает работу из-за отмены контекста")
			return
		case <-time.After(time.Hour * 1):
			checkMerlionNewPositions(ctx)
			createNewPositionsMS(ctx)
			updateRemainsMS(ctx)
		}
	}
}

// Проверяем Мерлион на наличие новых позиций и добавляем их в БД
func checkMerlionNewPositions(ctx context.Context) error {
	select {
	case <-ctx.Done():
		fmt.Println("checkMerlionNewPositions работу закончил из-за контекста")
		return nil
	default:
		fmt.Println("Начал проверку новых позиций на Мерлионе")
		catID, err := merlion.GetAllCatalogCodes()
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
		re := regexp.MustCompile(`(?i)dahua`)

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
				items, err := merlion.GetItemsByCatId(id)
				if err != nil {
					log.Printf("Ошибка при получении товаров по каталогу (checkMerlionNewPositions) id = %s: %s\n", id, err)
					//<-sem
				}
				//цикл по товарам из каталога
				for _, item := range items {
					if re.MatchString(item.Brand) {
						//Добавляем каждую позицию, т.к. в БД могут быть только уникальные записи
						added, err := dbInstance.AddCodeRecord(&typesDB.Codes{
							MoySklad:     "",
							Manufacturer: item.Vendor_part,
							Merlion:      item.No,
						})
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
			time.Sleep(time.Second / 2)
		}
		fmt.Println("Закончил проверку новых позиций на Мерлионе")

		//wg.Wait()
		//fmt.Printf("Всего горутин запущено: %d\n", idW)
		fmt.Printf("(%s) Добавлено %d записей при проверке позиций в Мерлионе\n", GetTimeNow(), countRecords)
		return nil
	}
}

// Создаем отсутствующие позиции на МС, добавляем в БД существующие
func createNewPositionsMS(ctx context.Context) error {
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

		items, err := dbInstance.GetCodeRecordsNoMS()
		if err != nil {
			log.Printf("Ошибка при получении записей из БД (createNewPositionsMS): %s\n", err)
			return err
		}
		//TODO Необходимо добавить парсинг последнего кода из базы мс, чтобы счетчик не задвоился
		counter, err := dbInstance.GetLastOwnIdMS()
		if err != nil {
			log.Printf("Ошибка при получении counter из БД (createNewPositionsMS): %s\n", err)
			return err
		}
		var createdItems int = 0
		for _, item := range *items {
			select {
			case <-ctx.Done():
				return nil
			default:
				manufacturerReplace := strings.Replace(item.Manufacturer, " ", "+", -1)
				response, err := sklad.GetItemByManufacturer(manufacturerReplace)
				if err != nil || response.StatusCode != 200 {
					log.Printf("Ошибка при получении записи из МС (createNewPositionsMS) manufacturer = %s StatusCode = %d: %s, \n", item.Manufacturer, response.StatusCode, err)
					continue
				}
				msItems := restTypes.SearchItem{}
				if err := json.Unmarshal(response.Body, &msItems); err != nil {
					log.Printf("Ошибка при декодировании item (createNewPositionsMS) manufacturer = %s: %s", item.Manufacturer, err)
					continue
				}
				//Если этого товара не существует на МС
				copyItem := &item
				if len(msItems.Rows) == 0 {
					if createNewItemMS(&item, copyItem, &counter, dbInstance) {
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
					itemMerlion, err := merlion.GetItemsByItemId(item.Merlion)
					if err != nil || len(itemMerlion) == 0 || itemMerlion[0].No == "" {
						log.Printf("Ошибка при получении записи из Мерлиона (createNewPositionsMS_2) merlionCode = %s: %s\n", item.Merlion, err)
						continue
					}
					// Вытягиваем -S0 номер из мерлиона, если есть
					substringsMer := strings.Fields(itemMerlion[0].Name)
					for _, subS := range substringsMer {
						merlionSnum = extractNumberFromS(subS)
						if merlionSnum >= 0 {
							break
						}
					}
					if merlionSnum == -1 {
						log.Printf("Ошибка при парсинге -S0 номера из Мерлиона (createNewPositionsMS_2) merlionCode = %s: %s\n", item.Merlion, err)
						continue
					}
					//TODO позиция с 3 артикулами, добавить проверку на уникальность кода мс
					for i := range msItems.Rows {
						// Если полностью совпадает
						if sklad.СontainsSubstring(msItems.Rows[i].Name, item.Manufacturer) {
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
						msManufacturer := ignoreDHManufacturer(msItems.Rows[i].Name)
						merManufacturer := ignoreDHManufacturer(item.Manufacturer)
						if sklad.СontainsSubstring(msManufacturer, merManufacturer) {
							dhiProblem = true
							break
						}
						// Если заканчивается на -S0
						if checkSManufacturer(msItems.Rows[i].Name, item.Manufacturer) && merlionSnum >= 0 {
							// Вытягиваем -S0 номер из мс, если есть
							substringsMS := strings.Fields(msItems.Rows[i].Name)
							for _, subS := range substringsMS {
								msSnum = extractNumberFromS(subS)
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
					}
					if msSnum == -1 {
						log.Printf("Ошибка при парсинге -S0 номера из мс manufacturer = %s\n", item.Manufacturer)
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
					if !founded {
						log.Printf("Полных соответствий не найдено manufacturer = %s (создаю новую позицию)", item.Manufacturer)
						if createNewItemMS(&item, copyItem, &counter, dbInstance) {
							createdItems++
						}
						continue
					}
					_, exists, err := dbInstance.GetCodeRecordByMS(article)
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
					copyItem.MoySklad = article
					//! опасный блок, может скопировать код чужой записи
					/*
						if article == "" {
							copyItem.MoySklad = msItems.Rows[0].Article
						} else {
							copyItem.MoySklad = article
						}*/
					// Извлекаем counter из артикула
					ownId, _ := extractCounterFromOwnID(copyItem.MoySklad)
					if ownId > 0 {
						copyItem.MsOwnId = ownId
					}

					if err = dbInstance.EditCodeRecord(copyItem); err != nil {
						log.Printf("Ошибка при изменении записи в БД (createNewPositionsMS) merlionCode = %s: %s\n", item.Merlion, err)
					}
				}
			}
			time.Sleep(time.Second / 2)
		}
		fmt.Println("Закончил сверять новые позиции с мс")
		fmt.Printf("(%s) Добавлено %d записей на МС\n", GetTimeNow(), createdItems)
		return nil
	}
}

// Обновляем остатки
func updateRemainsMS(ctx context.Context) error {
	fmt.Println("Начал обновлять остатки на мс")
	select {
	case <-ctx.Done():
		fmt.Println("updateRemainsMS работу закончил из-за контекста")
		return nil
	default:
		dbInstance, err := db.GetDBInstance()
		if err != nil {
			log.Printf("Ошибка при получении экземпляра БД (createNewPositionsMS): %s\n", err)
			return err
		}

		items, err := dbInstance.GetCodeRecordsFilledMS()
		if err != nil {
			log.Printf("Ошибка при получении записей из БД (createNewPositionsMS): %s\n", err)
			return err
		}
		//orgName := strings.Replace(keeper.K.GetOrgName(), " ", "+", -1)
		orgMeta, err := sklad.GetOrganizationMeta(keeper.K.GetOrgName())
		if err != nil {
			log.Printf("Ошибка при получении метаданных организации (updateRemainsMS): %s\n", err)
			return err
		}
		storeMeta, err := sklad.GetStoreMeta(keeper.K.GetSkladName())
		if err != nil || storeMeta.Href == "" {
			log.Printf("Ошибка при получении метаданных склада (updateRemainsMS): %s\n", err)
			return err
		}
		storeUUID, err := sklad.GetStoreUUID(keeper.K.GetSkladName())
		if err != nil || storeUUID == "" {
			log.Printf("Ошибка при получении UUID склада (updateRemainsMS): %s\n", err)
			return err
		}
		//На увеличение остатков
		acceptanceReq := restTypes.Acceptance{
			Organization: restTypes.MetaMiddle{Meta: orgMeta},
			Agent:        restTypes.MetaMiddle{Meta: orgMeta},
			Store:        restTypes.MetaMiddle{Meta: storeMeta},
			Description:  "Automatically generated by the script",
			Applicable:   true,
		}
		//На снижение остатков
		woffReq := restTypes.Acceptance{
			Organization: restTypes.MetaMiddle{Meta: orgMeta},
			Agent:        restTypes.MetaMiddle{Meta: orgMeta},
			Store:        restTypes.MetaMiddle{Meta: storeMeta},
			Description:  "Automatically generated by the script",
			Applicable:   true,
		}
		acceptanceList := []restTypes.PositionsAdd{}
		woffList := []restTypes.PositionsAdd{}
		for _, item := range *items {
			select {
			case <-ctx.Done():
				return nil
			default:
				itemUUID, isSerialTrackable, err := sklad.GetItemUUID(item.MoySklad)
				if err != nil || itemUUID == "" {
					log.Printf("Ошибка при получении UUID товара (updateRemainsMS) msCode = %s: %s\n", item.MoySklad, err)
					continue
				}
				if isSerialTrackable {
					log.Printf("Товар с серийным учетом в приемку/списание не попал (updateRemainsMS) msCode = %s isSerialTrackable = %t\n", item.MoySklad, isSerialTrackable)
					continue
				}
				itemRemainsMerlion, err := merlion.GetItemsAvailByItemId(item.Merlion)
				if err != nil || len(itemRemainsMerlion) == 0 || itemRemainsMerlion[0].No == "" {
					log.Printf("Ошибка при получении остатков с мерлиона (updateRemainsMS) merlionCode = %s: %s\n", item.Merlion, err)
					continue
				}
				itemRemainsMS, err := sklad.GetItemsAvail(itemUUID, storeUUID)
				if err != nil || itemRemainsMS < 0 {
					log.Printf("Ошибка при получении остатков с мс (updateRemainsMS) msCode = %s: %s\n", item.MoySklad, err)
					continue
				}
				itemMeta, err := sklad.GetItemMeta(item.MoySklad)
				if err != nil {
					log.Printf("Ошибка при получении метаданных товара (updateRemainsMS) msCode = %s: %s\n", item.MoySklad, err)
					continue
				}
				if itemRemainsMS > itemRemainsMerlion[0].AvailableClient_MSK {
					woffList = append(woffList, restTypes.PositionsAdd{
						Quantity:   itemRemainsMS - itemRemainsMerlion[0].AvailableClient_MSK,
						Assortment: restTypes.MetaMiddle{Meta: itemMeta},
					})
				} else if itemRemainsMS < itemRemainsMerlion[0].AvailableClient_MSK {
					acceptanceList = append(acceptanceList, restTypes.PositionsAdd{
						Quantity:   itemRemainsMerlion[0].AvailableClient_MSK - itemRemainsMS,
						Assortment: restTypes.MetaMiddle{Meta: itemMeta},
						Price:      itemRemainsMerlion[0].PriceClientRUB_MSK * 100,
					})
				}
			}
			time.Sleep(time.Second / 2)
		}
		if len(acceptanceList) != 0 {
			acceptanceReq.Positions = acceptanceList
			fmt.Printf("Позиций в приемке (updateRemainsMS): %d\n", len(acceptanceList))
			err = sklad.IncreaseItemsAvail(&acceptanceReq)
			if err != nil {
				log.Printf("Ошибка при создании приемки (updateRemainsMS): %s\n", err)
			}
		}
		if len(woffList) != 0 {
			woffReq.Positions = woffList
			fmt.Printf("Позиций в списании (updateRemainsMS): %d\n", len(woffList))
			err = sklad.DecreaseItemsAvail(&woffReq)
			if err != nil {
				log.Printf("Ошибка при создании списания (updateRemainsMS): %s\n", err)
			}
		}
		fmt.Println("Закончил обновлять остатки на мс")
		return nil
	}
}

// Увеличиваем счетчик для собственных артикулов
func incrementID(counter int64) string {
	newNumPart := counter
	newID := fmt.Sprintf("I%05d", newNumPart)
	return newID
}

// Игнорируем DH- и DHI-
func ignoreDHManufacturer(manufacturer string) string {
	if strings.HasPrefix(manufacturer, "DH-") {
		return manufacturer[3:]
	} else if strings.HasPrefix(manufacturer, "DHI-") {
		return manufacturer[4:]
	}
	return manufacturer
}

// Проверка на постфикс -S0
func checkSManufacturer(s string, substr string) bool {
	re := regexp.MustCompile(regexp.QuoteMeta(substr) + `-S\d`)
	matched := re.FindString(s)
	if matched != "" {
		return true
		//fmt.Println("Подстрока найдена:", matched)
	} else {
		return false
		//fmt.Println("Подстрока не найдена.")
	}
}

// Извлечение числа из постфикса -S0
// Возвращает -1 при ошибке парсинга
// Возвращает -2 если постфикса -S0 нет
func extractNumberFromS(s string) int {
	re := regexp.MustCompile(`-S(\d+)`)
	matches := re.FindStringSubmatch(s)

	if len(matches) > 1 {
		// Преобразование найденного числа в int
		number, err := strconv.Atoi(matches[1])
		if err != nil {
			return -1
		}
		return number
	}

	return -2
}

// Парсинг counter из артикула
func extractCounterFromOwnID(ownId string) (int64, error) {
	// Регулярное выражение для поиска числовой части в формате I00000
	re := regexp.MustCompile(`I(\d+)`)

	// Находим все совпадения
	matches := re.FindStringSubmatch(ownId)
	if len(matches) != 2 {
		return -1, fmt.Errorf("invalid code format")
	}

	// Извлекаем числовую часть из совпадений
	numStr := matches[1]

	// Преобразуем строку в int64
	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("failed to convert number: %w", err)
	}

	return num, nil
}

func createNewItemMS(item *typesDB.Codes, copyItem *typesDB.Codes, counter *int64, dbInstance *db.DBInstance) bool {
	newId := incrementID(*counter)
	itemMerlion, err := merlion.GetItemsByItemId(item.Merlion)
	if err != nil || len(itemMerlion) == 0 || itemMerlion[0].No == "" {
		log.Printf("Ошибка при получении записи из Мерлиона (createNewPositionsMS) merlionCode = %s: %s\n", item.Merlion, err)
		//continue
		return false
	}
	newItem := restTypes.CreateItem{
		Name:    itemMerlion[0].Name,
		Article: newId,
		ProductFolder: restTypes.MetaMiddle{
			Meta: restTypes.Meta{
				Href:         "https://api.moysklad.ru/api/remap/1.2/entity/productfolder/db27556f-153a-11f0-0a80-1751000ec1e4",
				MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/productfolder/metadata",
				Type:         "productfolder",
				MediaType:    "application/json",
				UuidHref:     "https://online.moysklad.ru/app/#good/edit?id=db27556f-153a-11f0-0a80-1751000ec1e4",
			},
		},
	}
	response, err := sklad.CreateItem(newItem)
	if err != nil || response.StatusCode != 200 {
		log.Printf("Ошибка при создании записи на МС (createNewPositionsMS) merlionCode = %s: %s\n", item.Merlion, err)
		//continue
		return false
	}
	copyItem.MoySklad = newId
	*counter++
	copyItem.MsOwnId = *counter
	if err = dbInstance.EditCodeRecord(copyItem); err != nil {
		log.Printf("Ошибка при изменении записи в БД (createNewPositionsMS) merlionCode = %s: %s\n", item.Merlion, err)
	}
	return true
}
