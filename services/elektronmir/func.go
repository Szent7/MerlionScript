package elektronmir

import (
	elektronmirTypes "MerlionScript/services/elektronmir/types"
)

func getGlobalItemsRecord(id int, GlobalItems []elektronmirTypes.Data) (record elektronmirTypes.Data, found bool) {
	for i := range GlobalItems {
		if GlobalItems[i].ID == id {
			return GlobalItems[i], true
		}
	}
	return elektronmirTypes.Data{}, false
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
