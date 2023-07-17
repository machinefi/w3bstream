package filesystem

import (
	"strconv"
	"strings"
)

type ObjectMeta struct {
	ObjectID    uint64 `json:"objectID"`
	Group       string `json:"group"`
	Size        int64  `json:"size"`
	ContentType string `json:"contentType"`
	ETag        string `json:"etag"`
}

func (meta ObjectMeta) Key() string {
	return meta.Group + "/" + strconv.FormatUint(meta.ObjectID, 10)
}

func ParseObjectMetaFromKey(key string) (*ObjectMeta, error) {
	parts := strings.Split(key, "/")
	if len(parts) != 2 {
		return nil, ErrInvalidObjectKey
	}
	grp := parts[0]

	oid, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil, ErrInvalidObjectKey
	}

	om := &ObjectMeta{ObjectID: oid, Group: grp}

	return om, nil
}
