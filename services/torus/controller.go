package torus

import (
	"MerlionScript/keeper"
	"MerlionScript/services/common"
	torusReq "MerlionScript/services/torus/requests"
	torusTypes "MerlionScript/services/torus/types"
	"MerlionScript/utils/db/interfaceDB"
	"MerlionScript/utils/dir"
	"MerlionScript/utils/excel"
	"context"
	"fmt"
	"log"
)

type TorusService struct {
	table       *excel.Workbook
	itemsGlobal map[string]int

	orgName   string
	storeName string
}

func Init() {
	srv := new(TorusService)
	srv.orgName = keeper.GetTorusOrg()
	srv.storeName = keeper.GetTorusSklad()
	newTable, err := excel.Open(torusTypes.TablePath)
	if err != nil {
		fmt.Printf("Ошибка при открытии таблицы torus: %s\n", err)
	}
	srv.table = newTable
	common.RegisteredServices[torusTypes.ServiceName] = common.ServiceInfo{
		DBTableName:     torusTypes.ServiceDBName,
		ServiceInstance: srv,
	}
	log.Printf("%s: сервис зарегистрирован\n", torusTypes.ServiceName)
}

func (srv *TorusService) GetSystemName() string {
	return torusTypes.ServiceName
}
func (srv *TorusService) GetDBTableName() string {
	return torusTypes.ServiceDBName
}

func (srv *TorusService) GetArticlesList(ctx context.Context) (*[]common.ArticleList, error) {
	select {
	case <-ctx.Done():
		fmt.Println("GetArticlesList работу закончил из-за контекста")
		return nil, nil
	default:
		if len(srv.itemsGlobal) == 0 {
			if exist, err := dir.CheckFileExists(torusTypes.TablePath); err != nil || !exist {
				return nil, nil
			}
			newTable, err := excel.Open(torusTypes.TablePath)
			if err != nil {
				fmt.Printf("Ошибка при открытии таблицы torus: %s\n", err)
				return nil, nil
			}
			srv.table = newTable
		}
		items, err := torusReq.GetItemsFormatted(srv.table)
		if err != nil {
			return nil, err
		}

		if len(*items) == 0 {
			return nil, fmt.Errorf("ошибка при получении номеров каталога: len(items) = 0")
		}

		articleList := make([]common.ArticleList, 0, 200)
		srv.itemsGlobal = *items
		//цикл по позициям
		for article := range *items {
			select {
			case <-ctx.Done():
				return &articleList, nil
			default:
				articleList = append(articleList, common.ArticleList{
					Article:     article,
					Brand:       "dahua",
					ServiceCode: article,
				})
			}
		}
		return &articleList, nil
	}

}

func (srv *TorusService) GetItemsList(ctx context.Context, dbInstance interfaceDB.DB) (*[]common.ItemList, error) {
	select {
	case <-ctx.Done():
		fmt.Println("GetItemsList работу закончил из-за контекста")
		return nil, nil
	default:
		if len(srv.itemsGlobal) == 0 {
			return nil, nil
		}
		items, err := dbInstance.GetServiceIDsNoMS(torusTypes.ServiceDBName)
		if err != nil {
			return nil, err
		}

		itemList := make([]common.ItemList, 0, len(*items))
		for i := range *items {
			itemList = append(itemList, common.ItemList{
				Article:      (*items)[i].Article,
				PositionName: (*items)[i].Article,
			})
		}

		return &itemList, nil
	}
}

func (srv *TorusService) GetStocksList(ctx context.Context, dbInstance interfaceDB.DB) (*map[string]common.StockList, error) {
	select {
	case <-ctx.Done():
		fmt.Println("GetStocksList работу закончил из-за контекста")
		return nil, nil
	default:
		if len(srv.itemsGlobal) == 0 {
			return nil, nil
		}

		items, err := dbInstance.GetServiceIDsFilledMS(torusTypes.ServiceDBName)
		if err != nil {
			return nil, err
		}

		stockList := make(map[string]common.StockList, len(*items))

		for i := range *items {
			quantity, ok := srv.itemsGlobal[(*items)[i].Article]
			if !ok {
				// нет данных по остаткам
				continue
			}
			stockList[(*items)[i].Article] = common.StockList{
				Stock: quantity,
				Price: 0,
			}
		}

		return &stockList, nil
	}
}

func (srv *TorusService) GetImagesList(ctx context.Context, code string) (*[]common.ImageList, error) {
	return &[]common.ImageList{}, nil
}

func (srv *TorusService) GetOrgName() string {
	return srv.orgName
}

func (srv *TorusService) GetStoreName() string {
	return srv.storeName
}

func (srv *TorusService) Finalize() {
	if srv.table != nil {
		srv.table.File.Close()
		srv.table = nil
		srv.itemsGlobal = nil
		dir.RemovePath(torusTypes.TablePath)
	}
}
