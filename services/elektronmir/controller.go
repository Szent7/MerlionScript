package elektronmir

import (
	"MerlionScript/keeper"
	"MerlionScript/services/common"
	elektronmirReq "MerlionScript/services/elektronmir/requests"
	elektronmirTypes "MerlionScript/services/elektronmir/types"
	"MerlionScript/utils/db/interfaceDB"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
)

// для дебага
/*func fillItemsGlobal(token string) {
	itemsGlobal = make([]netlab.ItemNetlab, 0, 200)
	catID, _ := netlabReq.GetAllCategoryCodes(token)
	for _, cat := range catID {
		items, _ := netlabReq.GetItemsByCatIdFormatted(cat.ID, token)
		itemsGlobal = append(itemsGlobal, items...)
	}
}*/

type ElektronmirService struct {
	itemsGlobal []elektronmirTypes.Data

	orgName   string
	storeName string
}

func Init() {
	srv := new(ElektronmirService)
	srv.orgName = keeper.K.GetElektronmirOrg()
	srv.storeName = keeper.K.GetElektronmirSkladOne()
	common.RegisteredServices[elektronmirTypes.ServiceName] = common.ServiceInfo{
		DBTableName:     elektronmirTypes.ServiceDBName,
		ServiceInstance: srv,
	}
	log.Printf("%s: сервис зарегистрирован\n", elektronmirTypes.ServiceName)
}

func (srv *ElektronmirService) GetSystemName() string {
	return elektronmirTypes.ServiceName
}
func (srv *ElektronmirService) GetDBTableName() string {
	return elektronmirTypes.ServiceDBName
}

func (srv *ElektronmirService) GetArticlesList(ctx context.Context) (*[]common.ArticleList, error) {
	token, err := elektronmirReq.GetNewToken()
	if err != nil || token == "" {
		return nil, err
	}
	select {
	case <-ctx.Done():
		fmt.Println("GetArticlesList работу закончил из-за контекста")
		return nil, nil
	default:
		catID, err := elektronmirReq.GetAllCategoryCodesFormatted(token)
		if err != nil {
			return nil, err
		}

		if len(catID.Data) == 0 {
			return nil, fmt.Errorf("ошибка при получении номеров каталога: len(catID) = 0")
		}

		articleList := make([]common.ArticleList, 0, 200)
		srv.itemsGlobal = make([]elektronmirTypes.Data, 0, 200)
		//цикл по каталогам
		for i, cat := range catID.Data {
			fmt.Printf("обработка каталога %s %d/%d\n", cat.Name, i+1, len(catID.Data))
			select {
			case <-ctx.Done():
				return &articleList, nil
			default:
				items, err := elektronmirReq.GetItemsByCatIdFormatted(cat.ID, token)
				if err != nil {
					log.Printf("ошибка при получении товаров по каталогу (GetArticlesList) id = %s: %s\n", cat.Name, err)
					continue
				}
				srv.itemsGlobal = append(srv.itemsGlobal, items.Data...)
				//цикл по товарам из каталога
				for _, item := range items.Data {
					articleList = append(articleList, common.ArticleList{
						Article:     item.Article,
						Brand:       item.Vendor,
						ServiceCode: strconv.Itoa(item.ID),
					})
				}
			}
			time.Sleep(time.Millisecond * 150) //из-за лимитов на запросы
		}
		return &articleList, nil
	}

}

//! нет проверки на то, что itemsGlobal пуст
func (srv *ElektronmirService) GetItemsList(ctx context.Context, dbInstance interfaceDB.DB) (*[]common.ItemList, error) {
	token, err := elektronmirReq.GetNewToken()
	if err != nil || token == "" {
		return nil, err
	}
	select {
	case <-ctx.Done():
		fmt.Println("GetItemsList работу закончил из-за контекста")
		return nil, nil
	default:
		items, err := dbInstance.GetServiceIDsNoMS(elektronmirTypes.ServiceDBName)
		if err != nil {
			return nil, err
		}

		itemList := make([]common.ItemList, 0, len(*items))
		for i := range *items {
			intCode, err := strconv.Atoi((*items)[i].ServiceCode)
			if err != nil {
				log.Printf("ошибка при конвертации id записи из Электронмира (GetItemsList) elektronmirCode = %s: %s\n", (*items)[i].ServiceCode, err)
				continue
			}
			itemCatalog, found := getGlobalItemsRecord(intCode, srv.itemsGlobal)
			if !found {
				log.Printf("ошибка при получении записи из Электронмира (GetItemsList) elektronmirCode = %s: %s\n", (*items)[i].ServiceCode, err)
				continue
			}
			itemList = append(itemList, common.ItemList{
				Article:      itemCatalog.Article,
				PositionName: itemCatalog.Name,
			})
		}

		return &itemList, nil
	}
}

func (srv *ElektronmirService) GetStocksList(ctx context.Context, dbInstance interfaceDB.DB) (*map[string]common.StockList, error) {
	token, err := elektronmirReq.GetNewToken()
	if err != nil || token == "" {
		return nil, err
	}
	select {
	case <-ctx.Done():
		fmt.Println("GetStocksList работу закончил из-за контекста")
		return nil, nil
	default:
		items, err := dbInstance.GetServiceIDsFilledMS(elektronmirTypes.ServiceDBName)
		if err != nil {
			return nil, err
		}

		stockList := make(map[string]common.StockList, len(*items))
		itemsArticle := make([]int, len(*items))
		for i := range *items {
			code, err := strconv.Atoi((*items)[i].ServiceCode)
			if err != nil {
				log.Printf("ошибка при конвертации id записи Elektronmir (GetStocksList): %s\n", err)
				continue
			}
			itemsArticle[i] = code
		}
		itemPart, err := elektronmirReq.GetItemsAvailByItemIdBatch(itemsArticle, token)
		if err != nil || len((*itemPart).Data) == 0 {
			return nil, fmt.Errorf("ошибка при получении остатков (GetStocksList): %s\n", err)
		}

		for i := range (*itemPart).Data {
			code := strconv.Itoa((*itemPart).Data[i].ID)
			for _, offer := range (*itemPart).Data[i].Offers {
				if offer.WarehouseID == 1 {
					stockList[code] = common.StockList{
						Stock: offer.Quantity,
						Price: float32(offer.PriceRRC * 100),
					}
					break
				}
			}

		}

		return &stockList, nil
	}
}

func (srv *ElektronmirService) GetImagesList(ctx context.Context, code string) (*[]common.ImageList, error) {
	// if srv.token == "" {
	// 	srv.once.Do(func() {
	// 		srv.StartTokenRefresh(ctx, time.Hour)
	// 	})
	// }
	// select {
	// case <-ctx.Done():
	// 	fmt.Println("GetImagesList работу закончил из-за контекста")
	// 	return nil, nil
	// default:
	// 	imagesURL, err := netlabReq.GetImagesByItemIdFormatted(code, srv.token)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("ошибка при получении записей из Нетлаба (GetImagesList) serviceCode = %s: %s\n", code, err)
	// 	}

	// 	imageList := make([]common.ImageList, 0, len(imagesURL))

	// 	for k, v := range imagesURL {
	// 		imageList = append(imageList, common.ImageList{
	// 			DownloadUrl: v,
	// 			Filename:    k,
	// 		})
	// 	}

	// 	return &imageList, nil
	// }
	return nil, nil
}

func (srv *ElektronmirService) GetOrgName() string {
	return srv.orgName
}

func (srv *ElektronmirService) GetStoreName() string {
	return srv.storeName
}
