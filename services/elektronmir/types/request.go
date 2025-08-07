package types

type TokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type ItemRequest struct {
	SectionID int `json:"section_id"`
}

type StockRequest struct {
	ProductIDs []int    `json:"product_ids"`
	Page       int      `json:"page"`
	Available  []string `json:"available"`
}
