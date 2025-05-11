package netlab

import (
	netlabTypesSoap "MerlionScript/types/soapTypes/netlab"
)

func getGlobalItemsRecord(id string, GlobalItems []netlabTypesSoap.ItemNetlab) (record netlabTypesSoap.ItemNetlab, found bool) {
	for i := range GlobalItems {
		if GlobalItems[i].Id == id {
			return GlobalItems[i], true
		}
	}
	return netlabTypesSoap.ItemNetlab{}, false
}

func getExtensionFromContentType(contentType string) string {
	switch contentType {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}
