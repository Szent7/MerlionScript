package types

import (
	"time"

	"github.com/google/uuid"
)

type TestProductGroup struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Code        string `json:"code,omitempty"`
}

type TestProductMeta struct {
	Meta Meta `json:"meta"`
}

type TestProduct struct {
	Name          string          `json:"name"`
	Code          string          `json:"code"`
	Vat           int             `json:"vat"`
	Weight        int             `json:"weight"`
	ProductFolder TestProductMeta `json:"productFolder"`
}

type ProductGroup struct {
	AccountId           uuid.UUID         `json:"accountId,omitempty"`
	Archived            bool              `json:"archived,omitempty"`
	Code                string            `json:"code,omitempty"`
	Description         string            `json:"description,omitempty"`
	EffectiveVat        int               `json:"effectiveVat,omitempty"`
	EffectiveVatEnabled bool              `json:"effectiveVatEnabled,omitempty"`
	ExternalCode        string            `json:"externalCode,omitempty"`
	Group               Meta              `json:"group,omitempty"`
	Id                  uuid.UUID         `json:"id,omitempty"`
	Meta                Meta              `json:"meta,omitempty"`
	Name                string            `json:"name"`
	Owner               Meta              `json:"owner,omitempty"`
	PathName            string            `json:"pathName,omitempty"`
	ProductFolder       Meta              `json:"productFolder,omitempty"`
	Shared              bool              `json:"shared,omitempty"`
	TaxSystem           map[string]string `json:"taxSystem,omitempty"` //enum
	Updated             time.Time         `json:"updated,omitempty"`   //DateTime
	UseParentVat        bool              `json:"useParentVat,omitempty"`
	Vat                 int               `json:"vat,omitempty"`
	VatEnabled          bool              `json:"vatEnabled,omitempty"`
}

type Meta struct {
	Href         string `json:"href,omitempty"`         //URL
	MetadataHref string `json:"metadataHref,omitempty"` //URL
	Type         string `json:"type,omitempty"`
	MediaType    string `json:"mediaType,omitempty"`
	UuidHref     string `json:"uuidHref,omitempty"`     //URL
	DownloadHref string `json:"downloadHref,omitempty"` //URL
}
