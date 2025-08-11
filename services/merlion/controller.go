package merlion

import (
	"MerlionScript/keeper"
	"MerlionScript/services/common"
	merlionReq "MerlionScript/services/merlion/requests"
	merlionTypes "MerlionScript/services/merlion/types"
	"MerlionScript/utils/db/interfaceDB"
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	batchSize = 500
)

type MerlionService struct {
	orgName   string
	storeName string
}

func Init() {
	srv := new(MerlionService)
	srv.orgName = keeper.GetMerlionOrg()
	srv.storeName = keeper.GetMerlionSklad()
	common.RegisteredServices[merlionTypes.ServiceName] = common.ServiceInfo{
		DBTableName:     merlionTypes.ServiceDBName,
		ServiceInstance: srv,
	}
	log.Printf("%s: сервис зарегистрирован\n", merlionTypes.ServiceName)
}

func (srv *MerlionService) GetSystemName() string {
	return merlionTypes.ServiceName
}
func (srv *MerlionService) GetDBTableName() string {
	return merlionTypes.ServiceDBName
}

func (srv *MerlionService) GetArticlesList(ctx context.Context) (*[]common.ArticleList, error) {
	select {
	case <-ctx.Done():
		fmt.Println("GetArticlesList работу закончил из-за контекста")
		return nil, nil
	default:
		catID, err := merlionReq.GetAllCatalogCodes()
		if err != nil {
			return nil, err
		}

		if len(catID) == 0 {
			return nil, fmt.Errorf("ошибка при получении номеров каталога: len(catID) = 0")
		}

		articleList := make([]common.ArticleList, 0, 100)

		//цикл по каталогам
		for i, id := range catID {
			fmt.Printf("обработка каталога %s %d/%d\n", id, i+1, len(catID))
			select {
			case <-ctx.Done():
				return &articleList, nil
			default:
				items, err := merlionReq.GetItemsByCatId(id)
				if err != nil {
					log.Printf("ошибка при получении товаров по каталогу (GetArticlesList) id = %s: %s\n", id, err)
					continue
				}
				//цикл по товарам из каталога
				for _, item := range items {
					lower := strings.ToLower(item.Brand)
					if strings.Contains(lower, "dahua") || strings.Contains(lower, "tenda") {
						articleList = append(articleList, common.ArticleList{
							Article:     item.Vendor_part,
							Brand:       lower,
							ServiceCode: item.No,
						})
					}
				}
			}
			time.Sleep(time.Millisecond * 150) //из-за лимитов на запросы
		}
		return &articleList, nil
	}
}

func (srv *MerlionService) GetItemsList(ctx context.Context, dbInstance interfaceDB.DB) (*[]common.ItemList, error) {
	select {
	case <-ctx.Done():
		fmt.Println("GetItemsList работу закончил из-за контекста")
		return nil, nil
	default:
		items, err := dbInstance.GetServiceIDsNoMS(merlionTypes.ServiceDBName)
		if err != nil {
			return nil, err
		}

		itemList := make([]common.ItemList, 0, len(*items))
		itemsArticle := make([]string, len(*items))
		for i := range *items {
			itemsArticle[i] = (*items)[i].ServiceCode
		}
		for i := 0; i < len(*items); i += batchSize {
			end := min(i+batchSize, len(*items))
			itemPart, err := merlionReq.GetItemsByItemIdBatch(itemsArticle[i:end])
			if err != nil || len(*itemPart) == 0 {
				return nil, fmt.Errorf("ошибка при получении позиций (GetItemsList): %s\n", err)
			}
			for _, itemRaw := range *itemPart {
				itemList = append(itemList, common.ItemList{
					Article:      itemRaw.Vendor_part,
					PositionName: itemRaw.Name,
				})
			}
			time.Sleep(time.Millisecond * 100)
		}

		return &itemList, nil
	}
}

func (srv *MerlionService) GetStocksList(ctx context.Context, dbInstance interfaceDB.DB) (*map[string]common.StockList, error) {
	select {
	case <-ctx.Done():
		fmt.Println("GetStocksList работу закончил из-за контекста")
		return nil, nil
	default:
		items, err := dbInstance.GetServiceIDsFilledMS(merlionTypes.ServiceDBName)
		if err != nil {
			return nil, err
		}

		stockList := make(map[string]common.StockList, len(*items))
		itemsArticle := make([]string, len(*items))
		for i := range *items {
			itemsArticle[i] = (*items)[i].ServiceCode
		}
		for i := 0; i < len(*items); i += batchSize {
			end := min(i+batchSize, len(*items))
			itemPart, err := merlionReq.GetItemsAvailByItemIdBatch(itemsArticle[i:end])
			if err != nil || len(*itemPart) == 0 {
				return nil, fmt.Errorf("ошибка при получении остатков (GetStocksList): %s\n", err)
			}
			for _, itemRaw := range *itemPart {
				stockList[itemRaw.No] = common.StockList{
					Stock: itemRaw.AvailableClient_MSK,
					Price: itemRaw.PriceClientRUB_MSK * 100,
				}
			}
			time.Sleep(time.Millisecond * 100)
		}

		return &stockList, nil
	}
}

func (srv *MerlionService) GetImagesList(ctx context.Context, code string) (*[]common.ImageList, error) {
	select {
	case <-ctx.Done():
		fmt.Println("GetImagesList работу закончил из-за контекста")
		return nil, nil
	default:
		imagesURL, err := merlionReq.GetImagesByItemIdFormatted(code)
		if err != nil {
			return nil, fmt.Errorf("ошибка при получении записей из мерлиона (GetImagesList) serviceCode = %s: %s\n", code, err)
		}

		imageList := make([]common.ImageList, 0, len(imagesURL))

		for i := range imagesURL {
			imageList = append(imageList, common.ImageList{
				DownloadUrl: merlionTypes.DownloadImageUrl + "/" + imagesURL[i],
				Filename:    imagesURL[i],
			})
		}

		return &imageList, nil
	}
}

func (srv *MerlionService) GetOrgName() string {
	return srv.orgName
}

func (srv *MerlionService) GetStoreName() string {
	return srv.storeName
}
