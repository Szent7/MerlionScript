package softtronik

import (
	"MerlionScript/keeper"
	"MerlionScript/services/common"
	softtronikReq "MerlionScript/services/softtronik/requests"
	softtronikTypes "MerlionScript/services/softtronik/types"
	"MerlionScript/utils/db/interfaceDB"
	"MerlionScript/utils/rest"
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"time"
)

// для дебага
/*func fillItemsGlobal(token string) {
	itemsGlobal = make([]netlab.ItemNetlab, 0, 200)
	catID, _ := softtronikReq.GetAllCategoryCodes(token)
	for _, cat := range catID {
		items, _ := softtronikReq.GetItemsByCatIdFormatted(cat.ID, token)
		itemsGlobal = append(itemsGlobal, items...)
	}
}*/

type SofttronikService struct {
	itemsGlobal []softtronikTypes.ProductItem
	orgName     string
	storeName   string
}

func Init() {
	srv := new(SofttronikService)
	srv.orgName = keeper.K.GetSofttronikOrg()
	srv.storeName = keeper.K.GetSofttronikSklad()
	common.RegisteredServices[softtronikTypes.ServiceName] = common.ServiceInfo{
		DBTableName:     softtronikTypes.ServiceDBName,
		ServiceInstance: srv,
	}
	log.Printf("%s: сервис зарегистрирован\n", softtronikTypes.ServiceName)
}

func (srv *SofttronikService) GetSystemName() string {
	return softtronikTypes.ServiceName
}
func (srv *SofttronikService) GetDBTableName() string {
	return softtronikTypes.ServiceDBName
}

func (srv *SofttronikService) GetArticlesList(ctx context.Context) (*[]common.ArticleList, error) {
	select {
	case <-ctx.Done():
		fmt.Println("GetArticlesList работу закончил из-за контекста")
		return nil, nil
	default:
		catID, err := softtronikReq.GetAllCategoryCodes()
		if err != nil {
			return nil, err
		}

		if len(catID) == 0 {
			return nil, fmt.Errorf("ошибка при получении номеров каталога: len(catID) = 0")
		}

		articleList := make([]common.ArticleList, 0, 100)
		srv.itemsGlobal = make([]softtronikTypes.ProductItem, 0, 200)
		//цикл по каталогам
		for i, cat := range catID {
			fmt.Printf("обработка каталога %s %d/%d\n", cat.Name, i+1, len(catID))
			select {
			case <-ctx.Done():
				return &articleList, nil
			default:
				items, err := softtronikReq.GetItemsByCatId(cat.ID)
				if err != nil {
					log.Printf("ошибка при получении товаров по каталогу (GetArticlesList) id = %s: %s\n", cat.Name, err)
					continue
				}
				srv.itemsGlobal = append(srv.itemsGlobal, items...)
				//цикл по товарам из каталога
				for _, item := range items {
					articleList = append(articleList, common.ArticleList{
						Article:     item.Article,
						Brand:       "dahua",
						ServiceCode: item.Code,
					})
				}
			}
			time.Sleep(time.Millisecond * 150) //из-за лимитов на запросы
		}
		return &articleList, nil
	}
}

//! нет проверки на то, что itemsGlobal пуст
func (srv *SofttronikService) GetItemsList(ctx context.Context, dbInstance interfaceDB.DB) (*[]common.ItemList, error) {
	select {
	case <-ctx.Done():
		fmt.Println("GetItemsList работу закончил из-за контекста")
		return nil, nil
	default:
		items, err := dbInstance.GetServiceIDsNoMS(softtronikTypes.ServiceDBName)
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
			postitionName := itemCatalog.Name
			if !strings.Contains(itemCatalog.Name, itemCatalog.Article) {
				postitionName = postitionName + " " + itemCatalog.Article
			}
			itemList = append(itemList, common.ItemList{
				Article:      itemCatalog.Article,
				PositionName: postitionName,
			})
		}

		return &itemList, nil
	}
}

func (srv *SofttronikService) GetStocksList(ctx context.Context, dbInstance interfaceDB.DB) (*map[string]common.StockList, error) {
	select {
	case <-ctx.Done():
		fmt.Println("GetStocksList работу закончил из-за контекста")
		return nil, nil
	default:
		items, err := dbInstance.GetServiceIDsFilledMS(softtronikTypes.ServiceDBName)
		if err != nil {
			return nil, err
		}
		//Записи из Софт-троника
		catID, err := softtronikReq.GetAllCategoryCodes()
		if err != nil {
			return nil, err
		}
		if len(catID) == 0 {
			return nil, fmt.Errorf("ошибка при получении номеров каталога: len(catID) = 0")
		}
		itemsSofttronik, err := softtronikReq.GetItemsAvailsAll(catID)
		if err != nil {
			log.Printf("Ошибка при получении записей из Софт-троника (updateRemainsMS): %s\n", err)
			return nil, err
		}
		itemStockSofttronik := getAvailsItemsRecords(itemsSofttronik)

		stockList := make(map[string]common.StockList, len(*items))
		for i := range *items {
			itemCatalog, ok := itemStockSofttronik[(*items)[i].Article]
			if !ok {
				//itemRemainsSofttronik, err = softtronikReq.GetItemsByItemIdFormatted(item.Service, token)
				//if err != nil {
				log.Printf("Ошибка при получении записи с Софт-троника (updateRemainsMS) softtronikCode = %s\n", (*items)[i].ServiceCode)
				continue
				//}
			}
			rub_price := itemCatalog.Price
			half_rub_price := float32(math.Ceil(rub_price * 100))
			stockList[(*items)[i].ServiceCode] = common.StockList{
				Stock: itemCatalog.Stocks,
				Price: half_rub_price,
			}
		}

		return &stockList, nil
	}
}

func (srv *SofttronikService) GetImagesList(ctx context.Context, code string) (*[]common.ImageList, error) {
	select {
	case <-ctx.Done():
		fmt.Println("GetImagesList работу закончил из-за контекста")
		return nil, nil
	default:
		listImages, err := softtronikReq.GetImagesByItemIdFormatted(code)
		if err != nil {
			log.Printf("Ошибка при получении записей из Софт-троника (UploadAllImages) serviceCode = %s: %s\n", code, err)
			return nil, err
		}
		imageList := make([]common.ImageList, 0, len(listImages))

		for k, v := range listImages {
			response, contentType, err := rest.CreateRequestImageHeader("GET", v, nil, "")
			if err != nil || response.StatusCode != 200 {
				log.Printf("Ошибка при получении изображений из Софт-троника (UploadAllImages) url = %s: %s\n", v, err)
				continue
			}
			ext := getExtensionFromContentType(contentType)
			if ext == ".jpg" || ext == ".png" {
				imageList = append(imageList, common.ImageList{
					DownloadUrl: v,
					Filename:    k + ext,
				})
			}

		}

		return &imageList, nil
	}
}

func (srv *SofttronikService) GetOrgName() string {
	return srv.orgName
}

func (srv *SofttronikService) GetStoreName() string {
	return srv.storeName
}
