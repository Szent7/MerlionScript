package netlab

import (
	"MerlionScript/keeper"
	"MerlionScript/services/common"
	netlabReq "MerlionScript/services/netlab/requests"
	netlabTypes "MerlionScript/services/netlab/types"
	"MerlionScript/utils/db/interfaceDB"
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"time"
)

// для дебага
/*
func (srv *NetlabService) fillItemsGlobal(token string) {
	srv.itemsGlobal = make([]netlabTypes.ItemNetlab, 0, 200)
	catID, _ := netlabReq.GetAllCategoryCodes(token)
	for _, cat := range catID {
		items, _ := netlabReq.GetItemsByCatIdFormatted(cat.ID, token)
		srv.itemsGlobal = append(srv.itemsGlobal, items...)
	}
}
*/
type NetlabService struct {
	itemsGlobal []netlabTypes.ItemNetlab

	orgName   string
	storeName string
}

func Init() {
	srv := new(NetlabService)
	srv.orgName = keeper.K.GetNetlabOrg()
	srv.storeName = keeper.K.GetNetlabSklad()
	common.RegisteredServices[netlabTypes.ServiceName] = common.ServiceInfo{
		DBTableName:     netlabTypes.ServiceDBName,
		ServiceInstance: srv,
	}
	log.Printf("%s: сервис зарегистрирован\n", netlabTypes.ServiceName)
}

func (srv *NetlabService) GetSystemName() string {
	return netlabTypes.ServiceName
}
func (srv *NetlabService) GetDBTableName() string {
	return netlabTypes.ServiceDBName
}

func (srv *NetlabService) GetArticlesList(ctx context.Context) (*[]common.ArticleList, error) {
	token, err := netlabReq.GetNewToken()
	if err != nil || token == "" {
		return nil, err
	}
	select {
	case <-ctx.Done():
		fmt.Println("GetArticlesList работу закончил из-за контекста")
		return nil, nil
	default:
		// srv.mu.RLock()
		catID, err := netlabReq.GetAllCategoryCodes(token)
		// srv.mu.RUnlock()
		if err != nil {
			return nil, err
		}

		if len(catID) == 0 {
			return nil, fmt.Errorf("ошибка при получении номеров каталога: len(catID) = 0")
		}

		articleList := make([]common.ArticleList, 0, 100)
		srv.itemsGlobal = make([]netlabTypes.ItemNetlab, 0, 200)
		//цикл по каталогам
		for i, cat := range catID {
			fmt.Printf("обработка каталога %s %d/%d\n", cat.Name, i+1, len(catID))
			select {
			case <-ctx.Done():
				return &articleList, nil
			default:
				// srv.mu.RLock()
				items, err := netlabReq.GetItemsByCatIdFormatted(cat.ID, token)
				// srv.mu.RUnlock()
				if err != nil {
					log.Printf("ошибка при получении товаров по каталогу (GetArticlesList) id = %s: %s\n", cat.Name, err)
					continue
				}
				srv.itemsGlobal = append(srv.itemsGlobal, items...)
				//цикл по товарам из каталога
				for _, item := range items {
					lower := strings.ToLower(item.ManufacturerName)
					if strings.Contains(lower, "dahua") || strings.Contains(lower, "tenda") {
						articleList = append(articleList, common.ArticleList{
							Article:     item.Manufacturer,
							Brand:       lower,
							ServiceCode: item.Id,
						})
					}
				}
			}
			time.Sleep(time.Millisecond * 150) //из-за лимитов на запросы
		}
		return &articleList, nil
	}
}

//! нет проверки на то, что itemsGlobal пуст
func (srv *NetlabService) GetItemsList(ctx context.Context, dbInstance interfaceDB.DB) (*[]common.ItemList, error) {
	select {
	case <-ctx.Done():
		fmt.Println("GetItemsList работу закончил из-за контекста")
		return nil, nil
	default:
		items, err := dbInstance.GetServiceIDsNoMS(netlabTypes.ServiceDBName)
		if err != nil {
			return nil, err
		}

		itemList := make([]common.ItemList, 0, len(*items))
		for i := range *items {
			itemCatalog, found := getGlobalItemsRecord((*items)[i].ServiceCode, srv.itemsGlobal)
			if !found {
				log.Printf("ошибка при получении записи из Нетлаба (GetItemsList) netlabCode = %s: %s\n", (*items)[i].ServiceCode, err)
				continue
			}
			itemList = append(itemList, common.ItemList{
				Article:      itemCatalog.Manufacturer,
				PositionName: itemCatalog.Name,
			})
		}

		return &itemList, nil
	}
}

//! нет проверки на то, что itemsGlobal пуст
func (srv *NetlabService) GetStocksList(ctx context.Context, dbInstance interfaceDB.DB) (*map[string]common.StockList, error) {
	token, err := netlabReq.GetNewToken()
	if err != nil || token == "" {
		return nil, err
	}
	// if len(srv.itemsGlobal) == 0 {
	// 	srv.fillItemsGlobal(token)
	// }
	select {
	case <-ctx.Done():
		fmt.Println("GetStocksList работу закончил из-за контекста")
		return nil, nil
	default:
		items, err := dbInstance.GetServiceIDsFilledMS(netlabTypes.ServiceDBName)
		if err != nil {
			return nil, err
		}
		//Курс доллара
		// srv.mu.RLock()
		currency, err := netlabReq.GetCurrency(token)
		// srv.mu.RUnlock()
		if err != nil || currency == 0 {
			log.Printf("ошибка при получении курса валют Нетлаб (GetStocksList): %s\n", err)
			return nil, err
		}

		stockList := make(map[string]common.StockList, len(*items))
		for i := range *items {
			itemCatalog, found := getGlobalItemsRecord((*items)[i].ServiceCode, srv.itemsGlobal)
			if !found {
				// itemRemainsNetlab, err = netlabReq.GetItemsByItemIdFormatted(item.Service, token)
				// if err != nil {
				// 	log.Printf("Ошибка при получении остатков с Нетлаба (updateRemainsMS) netlabCode = %s\n", item.Service)
				// 	continue
				// }
				log.Printf("ошибка при получении записи из Нетлаба (GetItemsList) netlabCode = %s: %s\n", (*items)[i].ServiceCode, err)
				continue
			}
			rub_price := itemCatalog.Price * currency
			half_rub_price := float32(math.Ceil(rub_price * 100))
			stockList[(*items)[i].ServiceCode] = common.StockList{
				Stock: itemCatalog.Remains,
				Price: half_rub_price,
			}
		}

		return &stockList, nil
	}
}

func (srv *NetlabService) GetImagesList(ctx context.Context, code string) (*[]common.ImageList, error) {
	token, err := netlabReq.GetNewToken()
	if err != nil || token == "" {
		return nil, err
	}
	select {
	case <-ctx.Done():
		fmt.Println("GetImagesList работу закончил из-за контекста")
		return nil, nil
	default:
		// srv.mu.RLock()
		imagesURL, err := netlabReq.GetImagesByItemIdFormatted(code, token)
		// srv.mu.RUnlock()
		if err != nil {
			return nil, fmt.Errorf("ошибка при получении записей из Нетлаба (GetImagesList) serviceCode = %s: %s\n", code, err)
		}

		imageList := make([]common.ImageList, 0, len(imagesURL))

		for k, v := range imagesURL {
			imageList = append(imageList, common.ImageList{
				DownloadUrl: v,
				Filename:    k,
			})
		}

		return &imageList, nil
	}
}

func (srv *NetlabService) GetOrgName() string {
	return srv.orgName
}

func (srv *NetlabService) GetStoreName() string {
	return srv.storeName
}
