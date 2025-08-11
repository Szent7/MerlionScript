package keeper

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var k *ShareData

func init() {
	k = new(ShareData)
}

type ShareData struct {
	MerlionCredentials string `mapstructure:"MERLION_CREDENTIALS"`
	MerlionOrg         string `mapstructure:"MERLION_ORGANIZATION"`
	MerlionSklad       string `mapstructure:"MERLION_SKLAD"`

	NetlabLogin    string `mapstructure:"NETLAB_LOGIN"`
	NetlabPassword string `mapstructure:"NETLAB_PASSWORD"`
	NetlabOrg      string `mapstructure:"NETLAB_ORGANIZATION"`
	NetlabSklad    string `mapstructure:"NETLAB_SKLAD"`

	SofttronikContractor  string `mapstructure:"SOFTTRONIK_CONTRACTOR_KEY"`
	SofttronikContractKey string `mapstructure:"SOFTTRONIK_CONTRACT_KEY"`
	SofttronikOrg         string `mapstructure:"SOFTTRONIK_ORGANIZATION"`
	SofttronikSklad       string `mapstructure:"SOFTTRONIK_SKLAD"`

	ElektronmirID       string `mapstructure:"ELEKTRONMIR_ID"`
	ElektronmirSecret   string `mapstructure:"ELEKTRONMIR_SECRET"`
	ElektronmirOrg      string `mapstructure:"ELEKTRONMIR_ORGANIZATION"`
	ElektronmirSkladOne string `mapstructure:"ELEKTRONMIR_SKLAD_ONE"`

	SkladToken   string `mapstructure:"MOY_SKLAD_TOKEN"`
	CatSkladName string `mapstructure:"CATALOG"`
}

func InitKeeper() { k.initKeeper() }
func (k *ShareData) initKeeper() {
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	//viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	secret := viper.GetString("MERLION_CREDENTIALS")
	if secret == "" {
		log.Fatal("Переменная 'MERLION_CREDENTIALS' не найдена в переменных окружения")
	}
	encodedSecret := base64.StdEncoding.EncodeToString([]byte(secret))
	viper.Set("MERLION_CREDENTIALS", encodedSecret)

	if err := viper.Unmarshal(k); err != nil {
		log.Fatalf("Ошибка разбора конфигурации: %v", err)
	}
}

// Merlion
func GetMerlionCredentials() string {
	return k.MerlionCredentials
}

func GetMerlionOrg() string {
	return k.MerlionOrg
}

func GetMerlionSklad() string {
	return k.MerlionSklad
}

// Netlab
func GetCredentialsNetlab() (string, string) {
	return k.NetlabLogin, k.NetlabPassword
}

func GetNetlabOrg() string {
	return k.NetlabOrg
}

func GetNetlabSklad() string {
	return k.NetlabSklad
}

// Softtronik
func GetSofttronikContractor() string {
	return k.SofttronikContractor
}

func GetSofttronikContractKey() string {
	return k.SofttronikContractKey
}

func GetSofttronikOrg() string {
	return k.SofttronikOrg
}

func GetSofttronikSklad() string {
	return k.SofttronikSklad
}

// Elektronmir
func GetElektronmirID() string {
	return k.ElektronmirID
}

func GetElektronmirSecret() string {
	return k.ElektronmirSecret
}

func GetElektronmirOrg() string {
	return k.ElektronmirOrg
}

func GetElektronmirSkladOne() string {
	return k.ElektronmirSkladOne
}

// MS
func GetMSCredentials() string {
	return k.SkladToken
}

func GetSkladCat() string {
	return k.CatSkladName
}
