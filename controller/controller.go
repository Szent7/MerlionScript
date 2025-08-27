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

// StartController управляет запуском рабочего цикла
func StartController(ctx context.Context, wg *sync.WaitGroup, dbInstance *db.DB) {
	defer wg.Done()

	var startHour int = 7 // Начальное время рабочего периода
	var endHour int = 19  // Конечное время рабочего периода

	// Проверка существования сервисов в системе
	if len(common.RegisteredServices) == 0 {
		fmt.Println("Контроллеры сервисов не обнаружены")
		return
	}

	for {
		now := time.Now()
		hour := now.Hour()

		//Если текущее время вне рабочего диапазона, выполняется резервное копирование и ожидание до следующего startHour
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

			// Если уже после 19:00
			if now.Hour() >= endHour {
				sleepUntil = sleepUntil.Add(24 * time.Hour)
			}

			sleepDuration := time.Until(sleepUntil) // Продолжительность до следующего startHour.
			timer := time.NewTimer(sleepDuration)   // Таймер ожидания
			select {
			case <-ctx.Done():
				fmt.Println("Контроллер завершает работу из-за отмены контекста")
				timer.Stop() // Остановка таймера, если контекст отменен
				return
			case <-timer.C:
				continue // Следующая итерация после периода ожидания
			}
		}

		//Рабочее время
		select {
		case <-ctx.Done():
			fmt.Println("StartController работу закончил из-за контекста")
			return
		default:
			ExecuteService(ctx, dbInstance) // Вызов рабочего цикла
		}
	}
}
