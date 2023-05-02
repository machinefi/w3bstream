package utils

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

func PrintResponse(cmd *cobra.Command, resp []byte) error {
	if !gjson.ValidBytes(resp) {
		return errors.New("failed to decode instance response")
	}
	ret := gjson.ParseBytes(resp)
	if code := ret.Get("code"); code.Exists() && code.Uint() != 0 {
		return errors.Errorf("error code: %d, error message: %s", code.Uint(), ret.Get("desc").String())
	}
	cmd.Println(ret.Get("desc").String())
	return nil
}
