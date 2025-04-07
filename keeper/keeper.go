package keeper

import "encoding/base64"

var K ShareData

type ShareData struct {
	skladCredentials   string
	merlionCredentials string
	skladName          string
	orgName            string
}

func (sd *ShareData) SetData(SkladCredentials string, MerlionCredentials string, SkladName string, OrgName string) {
	sd.skladCredentials = base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	sd.merlionCredentials = base64.StdEncoding.EncodeToString([]byte(MerlionCredentials))
	sd.skladName = SkladName
	sd.orgName = OrgName
}

func (sd *ShareData) GetCredentials() (string, string) {
	return sd.skladCredentials, sd.merlionCredentials
}

func (sd *ShareData) GetMSCredentials() string {
	return sd.skladCredentials
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
