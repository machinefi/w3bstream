package confs3

import "strconv"

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
