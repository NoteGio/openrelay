package objectstorage

import (
	"strings"
)

func GetStoredObject(uri string) StoredObject {
	if strings.HasPrefix(uri, "file://") {
		return &FileStorage{strings.TrimPrefix(uri, "file://")}
	} else if strings.HasPrefix(uri, "s3://") {
		parts := strings.SplitN(strings.TrimPrefix(uri, "s3://"), "/", 2)
		return &S3Storage{parts[0], parts[1]}
	}
	return nil
}
