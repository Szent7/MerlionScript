package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/allegro/bigcache/v3"
)

var Cache *bigcache.BigCache

func InitCache(ctx context.Context) error {
	config := bigcache.Config{
		Shards:             1024,
		LifeWindow:         30 * time.Minute,
		CleanWindow:        10 * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       1024, //в байтах
		HardMaxCacheSize:   512,  //в МБ
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	}
	var initErr error

	Cache, initErr = bigcache.New(ctx, config)
	if initErr != nil {
		return initErr
	}

	return nil
}

func Serialize[T any](data T) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Deserialize[T any](data []byte) (T, error) {
	var result T
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&result); err != nil {
		return result, err
	}
	return result, nil
}

func CacheRecord[T any](key string, record T) error {
	sRecord, err := Serialize(record)
	if err != nil {
		return err
	}

	Cache.Append(key, sRecord)
	return nil
}
