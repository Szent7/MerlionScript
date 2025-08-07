package merlion

/*
func GetItemAndCache(code string) (skladTypes.Rows, error) {
	itemMS, err := skladReq.GetItem(code)
	if err != nil {
		return skladTypes.Rows{}, err
	}

	if itemMS.Id == "" {
		return skladTypes.Rows{}, fmt.Errorf("пустой ID для кода: %s", code)
	}

	if err := cache.CacheRecord(code, itemMS); err != nil {
		return skladTypes.Rows{}, err
	}

	return itemMS, nil
}
*/
