package controller

import (
	"MerlionScript/services/common"
	"MerlionScript/src"
	"MerlionScript/utils/db"
	"context"
)

func ExecuteService(ctx context.Context, dbInstance *db.DB) {
	for _, v := range common.RegisteredServices {
		src.AddNewRecords(ctx, dbInstance, v.ServiceInstance)
		src.CreateNewPositionsERP(ctx, dbInstance, v.ServiceInstance, common.MainERP)
		src.UpdateRemainsERP(ctx, dbInstance, v.ServiceInstance, common.MainERP)
		src.UploadAllImages(ctx, dbInstance, v.ServiceInstance, common.MainERP)
	}
}
