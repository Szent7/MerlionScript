package controller

import (
	"MerlionScript/services/merlion"
	"MerlionScript/services/netlab"
	netlabReq "MerlionScript/services/netlab/requests"
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
		if (hour < startHour || hour >= endHour) && false {
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
			//merlionController(ctx)
			netlabController(ctx)
		}
	}
}

func merlionController(ctx context.Context) {
	merlion.CheckMerlionNewPositions(ctx)
	merlion.CreateNewPositionsMS(ctx)
	merlion.UpdateRemainsMS(ctx)
	//merlion.UploadAllImages(ctx)
	//merlion.UploadAllProperties(ctx)
}

func netlabController(ctx context.Context) {
	token, err := netlabReq.GetTempToken()
	if err != nil {
		log.Printf("Ошибка при получении токена Netlab (netlabController): %s\n", err)
		return
	}
	//netlab.CheckNetlabNewPositions(ctx, token)
	//netlab.CreateNewPositionsMS(ctx, token)
	//netlab.UpdateRemainsMS(ctx, token)
	netlab.UploadAllImages(ctx, token)
	//netlab.UploadAllProperties(ctx)
}
