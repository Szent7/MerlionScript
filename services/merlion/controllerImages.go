package merlion

import (
	"context"
	"fmt"
)

func UploadAllImages(ctx context.Context) error {
	select {
	case <-ctx.Done():
		fmt.Println("UploadAllImages работу закончил из-за контекста")
		return nil
	default:
		return nil
	}
}
