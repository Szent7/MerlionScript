package restTypes

type Meta struct {
	Href         string `json:"href,omitempty"`         //URL
	MetadataHref string `json:"metadataHref,omitempty"` //URL
	Type         string `json:"type,omitempty"`
	MediaType    string `json:"mediaType,omitempty"`
	UuidHref     string `json:"uuidHref,omitempty"`     //URL
	DownloadHref string `json:"downloadHref,omitempty"` //URL
}
