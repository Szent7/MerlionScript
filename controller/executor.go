package controller

import (
	"MerlionScript/services/common"
	"MerlionScript/src"
	"MerlionScript/utils/db"
	"context"
)

// ExecuteService выполняет логику рабочего цикла для каждого сервиса
func ExecuteService(ctx context.Context, dbInstance *db.DB) {
	// Цикл по зарегистрированным сервисам
	for _, v := range common.RegisteredServices {
		src.AddNewRecords(ctx, dbInstance, v.ServiceInstance)                         // Добавление новых позиций в БД
		src.CreateNewPositionsERP(ctx, dbInstance, v.ServiceInstance, common.MainERP) // Создание новых позиций/привязок в ERP системе
		src.UpdateRemainsERP(ctx, dbInstance, v.ServiceInstance, common.MainERP)      // Обновление остатков в ERP системе
		src.UploadAllImages(ctx, dbInstance, v.ServiceInstance, common.MainERP)       // Загрузка изображений в ERP систему
		v.ServiceInstance.Finalize()                                                  // Завершение цикла сервиса (финализация)
	}
}
