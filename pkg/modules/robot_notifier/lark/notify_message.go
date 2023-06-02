package lark

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	"github.com/machinefi/w3bstream/pkg/types"
)

type Content struct {
	Post map[string]ContentBody `json:"post"`
}

type ContentBody struct {
	Title           string                     `json:"title"`
	ContentElements [][]map[string]interface{} `json:"content"`
}

func TagElementToMap(e TagElement) (map[string]interface{}, error) {
	ret := map[string]interface{}{"tag": e.Tag()}

	err := mapstructure.Decode(e, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func TagElementsToMap(elem ...TagElement) ([][]map[string]interface{}, error) {
	ms := make([][]map[string]interface{}, 0, len(elem))
	for _, e := range elem {
		m, err := TagElementToMap(e)
		if err != nil {
			return nil, err
		}
		ms = append(ms, []map[string]interface{}{m})
	}
	return ms, nil
}

type TagElement interface {
	Tag() string
}

type ElementImg struct {
	ImageKey string `mapstructure:"image_key"`
	W        int    `mapstructure:"width"`
	H        int    `mapstructure:"height"`
}

func (v *ElementImg) Tag() string { return "img" }

type ElementText struct {
	Text string `mapstructure:"text"`
}

func (v *ElementText) Tag() string { return "text" }

type ElementHref struct {
	Text string `mapstructure:"text"`
	Ref  string `mapstructure:"href"`
}

func (v *ElementHref) Tag() string { return "a" }

type ElementPin struct {
	UserID string `mapstructure:"user_id"`
}

func (v *ElementPin) Tag() string { return "at" }

func Build(ctx context.Context, title, level, content string) ([]byte, error) {
	notifier, _ := types.RobotNotifierConfigFromContext(ctx)
	if notifier == nil {
		return nil, nil
	}
	ts := time.Now().UTC().Unix()
	sign := ""
	if notifier.SignFn != nil {
		signature, err := notifier.SignFn(ts)
		if err != nil {
			return nil, errors.Wrap(err, "signature")
		}
		sign = signature
	}

	ms, err := TagElementsToMap(
		&ElementText{
			Text: fmt.Sprintf("env: %s", notifier.Env),
		},
		&ElementText{
			Text: fmt.Sprintf("project: %s", os.Getenv(consts.EnvProjectName)),
		},
		&ElementText{
			Text: fmt.Sprintf("version: %s", os.Getenv(consts.EnvProjectVersion)),
		},
		&ElementText{
			Text: content,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "mapstructure")
	}

	msg := &struct {
		Timestamp   int64   `json:"timestamp,omitempty"`
		Sign        string  `json:"sign,omitempty"`
		MessageType string  `json:"msg_type"`
		Content     Content `json:"content"`
	}{
		Timestamp:   ts,
		Sign:        sign,
		MessageType: "post",
		Content: Content{
			Post: map[string]ContentBody{
				"en_ch": {
					Title:           fmt.Sprintf("%s [%s]", strings.ToUpper(level), title),
					ContentElements: ms,
				},
			},
		},
	}

	return json.Marshal(msg)
}

func ResponseHook(body []byte) error {
	if len(body) == 0 {
		return nil
	}

	v := &struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{}

	err := json.Unmarshal(body, v)
	if err != nil {
		return err
	}

	switch v.Code {
	case 19021:
		return ErrInvalidSignature
	case 9499:
		return ErrInvalidMessage
	default:
		if v.Code == 0 {
			return nil
		} else {
			return errors.Errorf("code: %d, msg: %s", v.Code, v.Msg)
		}
	}
}

var (
	ErrInvalidMessage   = errors.New("invalid message")
	ErrInvalidSignature = errors.New("invalid signature")
)
