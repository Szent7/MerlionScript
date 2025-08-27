package infocom

import (
	"MerlionScript/keeper"
	"MerlionScript/services/common"
	infocomReq "MerlionScript/services/infocom/requests"
	infocomTypes "MerlionScript/services/infocom/types"
	"MerlionScript/utils/db/interfaceDB"
	"MerlionScript/utils/dir"
	"MerlionScript/utils/excel"
	"context"
	"fmt"
	"log"
)

type InfocomService struct {
	table       *excel.Workbook
	itemsGlobal map[string]int

	orgName   string
	storeName string
}

func Init() {
	srv := new(InfocomService)
	srv.orgName = keeper.GetInfocomOrg()
	srv.storeName = keeper.GetInfocomSklad()
	newTable, err := excel.Open(infocomTypes.TablePath)
	if err != nil {
		fmt.Printf("Ошибка при открытии таблицы infocom: %s\n", err)
	}
	srv.table = newTable
	common.RegisteredServices[infocomTypes.ServiceName] = common.ServiceInfo{
		DBTableName:     infocomTypes.ServiceDBName,
		ServiceInstance: srv,
	}
	log.Printf("%s: сервис зарегистрирован\n", infocomTypes.ServiceName)
}

func (srv *InfocomService) GetSystemName() string {
	return infocomTypes.ServiceName
}
func (srv *InfocomService) GetDBTableName() string {
	return infocomTypes.ServiceDBName
}

func (srv *InfocomService) GetArticlesList(ctx context.Context) (*[]common.ArticleList, error) {
	select {
	case <-ctx.Done():
		fmt.Println("GetArticlesList работу закончил из-за контекста")
		return nil, nil
	default:
		if len(srv.itemsGlobal) == 0 {
			if exist, err := dir.CheckFileExists(infocomTypes.TablePath); err != nil || !exist {
				return nil, nil
			}
			newTable, err := excel.Open(infocomTypes.TablePath)
			if err != nil {
				fmt.Printf("Ошибка при открытии таблицы infocom: %s\n", err)
				return nil, nil
			}
			srv.table = newTable
		}
		items, err := infocomReq.GetItemsFormatted(srv.table)
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

func (srv *InfocomService) GetItemsList(ctx context.Context, dbInstance interfaceDB.DB) (*[]common.ItemList, error) {
	select {
	case <-ctx.Done():
		fmt.Println("GetItemsList работу закончил из-за контекста")
		return nil, nil
	default:
		if len(srv.itemsGlobal) == 0 {
			return nil, nil
		}
		items, err := dbInstance.GetServiceIDsNoMS(infocomTypes.ServiceDBName)
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

func (srv *InfocomService) GetStocksList(ctx context.Context, dbInstance interfaceDB.DB) (*map[string]common.StockList, error) {
	select {
	case <-ctx.Done():
		fmt.Println("GetStocksList работу закончил из-за контекста")
		return nil, nil
	default:
		if len(srv.itemsGlobal) == 0 {
			return nil, nil
		}

		items, err := dbInstance.GetServiceIDsFilledMS(infocomTypes.ServiceDBName)
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

func (srv *InfocomService) GetImagesList(ctx context.Context, code string) (*[]common.ImageList, error) {
	return &[]common.ImageList{}, nil
}

func (srv *InfocomService) GetOrgName() string {
	return srv.orgName
}

func (srv *InfocomService) GetStoreName() string {
	return srv.storeName
}

func (srv *InfocomService) Finalize() {
	if srv.table != nil {
		srv.table.File.Close()
		srv.table = nil
		srv.itemsGlobal = nil
		dir.RemovePath(infocomTypes.TablePath)
	}
}
