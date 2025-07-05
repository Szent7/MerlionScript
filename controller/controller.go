package controller

import (
	"MerlionScript/services/merlion"
	"MerlionScript/services/netlab"
	netlabReq "MerlionScript/services/netlab/requests"
	"MerlionScript/services/softtronik"
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

func StartController(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	var startHour int = 7
	var endHour int = 19
	for {
		now := time.Now()
		hour := now.Hour()

		//Если текущее время вне рабочего диапазона — ждём до следующей "startHour"
		if hour < startHour || hour >= endHour {
			fmt.Println("Вне рабочего времени. Ожидаем следующего старта.")
			sleepUntil := time.Date(
				now.Year(), now.Month(), now.Day(),
				startHour, 0, 0, 0, now.Location(),
			)

			//Если уже позже 19:00, ждём до завтрашнего 7:00
			if now.Hour() >= endHour {
				sleepUntil = sleepUntil.Add(24 * time.Hour)
			}

			sleepDuration := time.Until(sleepUntil)
			timer := time.NewTimer(sleepDuration)
			select {
			case <-ctx.Done():
				fmt.Println("Контроллер завершает работу из-за отмены контекста")
				timer.Stop()
				return
			case <-timer.C:
				continue //Начинаем новый рабочий период
			}
		}
		//Рабочее время
		select {
		case <-ctx.Done():
			fmt.Println("StartController работу закончил из-за контекста")
			return
		default:
			merlionController(ctx)
			netlabController(ctx)
			softtronikController(ctx)
		}
	}
}

func merlionController(ctx context.Context) {
	merlion.CheckMerlionNewPositions(ctx)
	merlion.CreateNewPositionsMS(ctx)
	merlion.UpdateRemainsMS(ctx)
	merlion.UploadAllImages(ctx)
}

func netlabController(ctx context.Context) {
	token, err := netlabReq.GetTempToken()
	if err != nil {
		log.Printf("Ошибка при получении токена Netlab (netlabController): %s\n", err)
		return
	}
	netlab.CheckNetlabNewPositions(ctx, token)
	netlab.CreateNewPositionsMS(ctx, token)
	netlab.UpdateRemainsMS(ctx, token)
	netlab.UploadAllImages(ctx, token)
}

func softtronikController(ctx context.Context) {
	softtronik.CheckSofttronikNewPositions(ctx)
	softtronik.CreateNewPositionsMS(ctx)
	softtronik.UpdateRemainsMS(ctx)
	softtronik.UploadAllImages(ctx)
}

/*
func fixNames() {
	var itemsGlobal []softtronikTypes.ProductItem
	dbInstance, err := db.GetDBInstance()
	if err != nil {
		log.Printf("Ошибка при получении экземпляра БД (CheckSofttronikNewPositions): %s\n", err)
		return
	}
	records, err := dbInstance.GetCodeRecords(typesDB.SofttronikTable)
	if err != nil {
		log.Printf("Ошибка при получении записей из БД (CheckSofttronikNewPositions): %s\n", err)
		return
	}
	sort.Slice(*records, func(i, j int) bool {
		return (*records)[i].MoySklad < (*records)[j].MoySklad
	})
	re := regexp.MustCompile(`^I\d{5}$`)
	catID, err := softtronikReq.GetAllCategoryCodes()
	if err != nil {
		panic(err)
	}
	if len(catID) == 0 {
		panic(fmt.Errorf("ошибка при получении номеров каталога: len(catID) = 0"))
	}
	for _, cat := range catID {
		items, err := softtronikReq.GetItemsByCatId(cat.ID)
		if err != nil {
			log.Printf("Ошибка при получении товаров по каталогу (CheckSofttronikNewPositions) id = %s: %s\n", cat.Name, err)
			return
		}
		itemsGlobal = append(itemsGlobal, items...)
	}
	itemsMap := make(map[string]softtronikTypes.ProductItem, len(itemsGlobal))
	for i := range itemsGlobal {
		itemsMap[itemsGlobal[i].Article] = itemsGlobal[i]
	}
	for _, record := range *records {
		if re.MatchString(record.MoySklad) {
			// Извлечение числа из строки
			numStr := record.MoySklad[1:]
			num, err := strconv.Atoi(numStr)
			if err != nil {
				panic(err)
			}
			if num >= 382 {
				oldName, ok := itemsMap[record.Manufacturer]
				if !ok {
					panic(fmt.Errorf("Ключ не найден"))
				}
				validName := oldName.Name + " " + record.Manufacturer

				itemMS, err := skladReq.GetItem(record.MoySklad)
				if err != nil || itemMS.Id == "" {
					log.Printf("Ошибка при получении товара МС (updateRemainsMS) msCode = %s: %s\n", record.MoySklad, err)
					return
				}

				err1 := UpdateItem(NewItem{
					Name: validName,
				}, itemMS.Id)
				if err1 == nil {
					fmt.Printf("Название изменено %s -> %s (MScode: %s)\n", oldName.Name, validName, record.MoySklad)
				}
			}
		}
	}
}

type NewItem struct {
	Name string `json:"name"`
}

func UpdateItem(product NewItem, uuid string) error {
	//jsonBody, err := json.Marshal(reqBody)
	jsonBody, err := json.MarshalIndent(product, "", "  ")
	//fmt.Println("тело запроса в JSON:", string(jsonBody))
	if err != nil {
		log.Println("ошибка при преобразовании структуры в JSON:", err)
		return err
	}
	//authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
	//url := "https://api.moysklad.ru/api/remap/1.2/entity/product"
	body, err := rest.CreateRequestMS("PUT", "https://api.moysklad.ru/api/remap/1.2/entity/product/"+uuid, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	if body.StatusCode != 200 {
		fmt.Println(body.StatusCode)
		fmt.Println(string(body.Body))
		return fmt.Errorf("bad code")
	}
	return nil
}
*/
