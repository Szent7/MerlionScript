package types

const (
	TokenUrl = "http://services.netlab.ru/rest/authentication/token.json?username=%s&password=%s"
	//CatalogUrl  = "http://services.netlab.ru/rest/catalogsZip/list.xml?oauth_token=%s"
	CategoryUrl = "http://services.netlab.ru/rest/catalogsZip/%s.xml?oauth_token=%s"
	ItemUrl     = "http://services.netlab.ru/rest/catalogsZip/%s/%s.xml?oauth_token=%s"
	ItemIdUrl   = "http://services.netlab.ru/rest/catalogsZip/goodsByUid/%s.xml?oauth_token=%s"
	CurrencyUrl = "http://services.netlab.ru/rest/catalogsZip/info.xml?oauth_token=%s"
	ImageUrl    = "http://services.netlab.ru/rest/catalogsZip/goodsImages/%s.xml?oauth_token=%s"

	CatalogName = "В наличии"
)
