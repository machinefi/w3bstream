package patch_std

import (
	"encoding/json"

	"github.com/agiledragon/gomonkey/v2"
)

func JsonMarshal(patch *gomonkey.Patches, data []byte, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		json.Marshal,
		func(_ interface{}) ([]byte, error) { return data, err },
	)
}

func JsonUnmarshal(patch *gomonkey.Patches, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		json.Unmarshal,
		func(_ []byte, _ interface{}) error { return err },
	)
}
