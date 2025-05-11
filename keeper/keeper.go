package keeper

import (
	"encoding/base64"
)

var K ShareData
var MerlionMainURL = "https://api.merlion.com/rl/mlservice3"

const (
	MerlionCredentialsEnv = "MERLION_CREDENTIALS"
	MerlionOrgEnv         = "MERLION_ORGANIZATION"
	MerlionSkladEnv       = "MERLION_SKLAD"

	NetlabLoginEnv    = "NETLAB_LOGIN"
	NetlabPasswordEnv = "NETLAB_PASSWORD"
	NetlabOrgEnv      = "NETLAB_ORGANIZATION"
	NetlabSkladEnv    = "NETLAB_SKLAD"

	SkladTokenEnv   = "MOY_SKLAD_TOKEN"
	CatSkladNameEnv = "CATALOG"
)

type ShareData struct {
	merlionCredentials string
	merlionOrg         string
	merlionSklad       string

	netlabLogin    string
	netlabPassword string
	netlabOrg      string
	netlabSklad    string

	skladToken   string
	catSkladName string
}

func (sd *ShareData) SetData(data map[string]string) {
	sd.merlionCredentials = base64.StdEncoding.EncodeToString([]byte(data[MerlionCredentialsEnv]))
	sd.merlionOrg = data[MerlionOrgEnv]
	sd.merlionSklad = data[MerlionSkladEnv]

	sd.netlabLogin = data[NetlabLoginEnv]
	sd.netlabPassword = data[NetlabPasswordEnv]
	sd.netlabOrg = data[NetlabOrgEnv]
	sd.netlabSklad = data[NetlabSkladEnv]

	sd.skladToken = data[SkladTokenEnv]
	sd.catSkladName = data[CatSkladNameEnv]
}

// Merlion
func (sd *ShareData) GetMerlionCredentials() string {
	return sd.merlionCredentials
}

func (sd *ShareData) GetMerlionOrg() string {
	return sd.merlionOrg
}

func (sd *ShareData) GetMerlionSklad() string {
	return sd.merlionSklad
}

// Netlab
func (sd *ShareData) GetCredentialsNetlab() (string, string) {
	return sd.netlabLogin, sd.netlabPassword
}

func (sd *ShareData) GetNetlabOrg() string {
	return sd.netlabOrg
}

func (sd *ShareData) GetNetlabSklad() string {
	return sd.netlabSklad
}

// MS
func (sd *ShareData) GetMSCredentials() string {
	return sd.skladToken
}

func (sd *ShareData) GetSkladCat() string {
	return sd.catSkladName
}
