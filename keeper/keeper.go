package keeper

import "encoding/base64"

var K ShareData
var MerlionMainURL = "https://api.merlion.com/rl/mlservice3"

type ShareData struct {
	skladCredentials   string
	skladToken         string
	merlionCredentials string
	skladName          string
	orgName            string
	catSkladName       string
}

func (sd *ShareData) SetData(SkladToken string, MerlionCredentials string, SkladName string, OrgName string, catSkladName string) {
	//sd.skladCredentials = base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	sd.skladToken = SkladToken
	sd.merlionCredentials = base64.StdEncoding.EncodeToString([]byte(MerlionCredentials))
	sd.skladName = SkladName
	sd.orgName = OrgName
	sd.catSkladName = catSkladName
}

func (sd *ShareData) GetCredentials() (string, string) {
	return sd.skladCredentials, sd.merlionCredentials
}

func (sd *ShareData) GetMSCredentials() string {
	return sd.skladToken
}

func (sd *ShareData) GetMerlionCredentials() string {
	return sd.merlionCredentials
}

func (sd *ShareData) GetSkladName() string {
	return sd.skladName
}

func (sd *ShareData) GetOrgName() string {
	return sd.orgName
}

func (sd *ShareData) GetSkladCat() string {
	return sd.catSkladName
}
