package controller

import (
	"MerlionScript/keeper"
	"MerlionScript/services/common"
	"MerlionScript/utils/backup"
	"MerlionScript/utils/db"
	"context"
	"fmt"
	"sync"
	"time"
)

func StartController(ctx context.Context, wg *sync.WaitGroup, dbInstance *db.DB) {
	defer wg.Done()
	var startHour int = 7
	var endHour int = 19
	if len(common.RegisteredServices) == 0 {
		fmt.Println("Контроллеры сервисов не обнаружены")
		return
	}
	for {
		now := time.Now()
		hour := now.Hour()

		//Если текущее время вне рабочего диапазона — ждём до следующего startHour
		if hour < startHour || hour >= endHour {
			if keeper.GetBackupToggle() {
				fmt.Println("Создание бэкапа...")
				if err := backup.CreateDefaultBackup(); err != nil {
					fmt.Printf("Ошибка при создании бэкапа: %s\n", err)
				}
			}

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
			ExecuteService(ctx, dbInstance)
		}
	}
}
