// This is a generated source file. DO NOT EDIT
// Source: jwt/signing_method__generated.go

package jwt

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidSigningMethod = errors.New("invalid SigningMethod type")

func ParseSigningMethodFromString(s string) (SigningMethod, error) {
	switch s {
	default:
		return SIGNING_METHOD_UNKNOWN, InvalidSigningMethod
	case "":
		return SIGNING_METHOD_UNKNOWN, nil
	case "ECDSA256":
		return SIGNING_METHOD__ECDSA256, nil
	case "ECDSA384":
		return SIGNING_METHOD__ECDSA384, nil
	case "ECDSA512":
		return SIGNING_METHOD__ECDSA512, nil
	case "HMAC256":
		return SIGNING_METHOD__HMAC256, nil
	case "HMAC384":
		return SIGNING_METHOD__HMAC384, nil
	case "HMAC512":
		return SIGNING_METHOD__HMAC512, nil
	case "RSA256":
		return SIGNING_METHOD__RSA256, nil
	case "RSA384":
		return SIGNING_METHOD__RSA384, nil
	case "RSA512":
		return SIGNING_METHOD__RSA512, nil
	case "RSAPSS256":
		return SIGNING_METHOD__RSAPSS256, nil
	case "RSAPSS384":
		return SIGNING_METHOD__RSAPSS384, nil
	case "RSAPSS512":
		return SIGNING_METHOD__RSAPSS512, nil
	}
}

func ParseSigningMethodFromLabel(s string) (SigningMethod, error) {
	switch s {
	default:
		return SIGNING_METHOD_UNKNOWN, InvalidSigningMethod
	case "":
		return SIGNING_METHOD_UNKNOWN, nil
	case "ECDSA256":
		return SIGNING_METHOD__ECDSA256, nil
	case "ECDSA384":
		return SIGNING_METHOD__ECDSA384, nil
	case "ECDSA512":
		return SIGNING_METHOD__ECDSA512, nil
	case "HMAC256":
		return SIGNING_METHOD__HMAC256, nil
	case "HMAC384":
		return SIGNING_METHOD__HMAC384, nil
	case "HMAC512":
		return SIGNING_METHOD__HMAC512, nil
	case "RSA256":
		return SIGNING_METHOD__RSA256, nil
	case "RSA384":
		return SIGNING_METHOD__RSA384, nil
	case "RSA512":
		return SIGNING_METHOD__RSA512, nil
	case "RSAPSS256":
		return SIGNING_METHOD__RSAPSS256, nil
	case "RSAPSS384":
		return SIGNING_METHOD__RSAPSS384, nil
	case "RSAPSS512":
		return SIGNING_METHOD__RSAPSS512, nil
	}
}

func (v SigningMethod) Int() int {
	return int(v)
}

func (v SigningMethod) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case SIGNING_METHOD_UNKNOWN:
		return ""
	case SIGNING_METHOD__ECDSA256:
		return "ECDSA256"
	case SIGNING_METHOD__ECDSA384:
		return "ECDSA384"
	case SIGNING_METHOD__ECDSA512:
		return "ECDSA512"
	case SIGNING_METHOD__HMAC256:
		return "HMAC256"
	case SIGNING_METHOD__HMAC384:
		return "HMAC384"
	case SIGNING_METHOD__HMAC512:
		return "HMAC512"
	case SIGNING_METHOD__RSA256:
		return "RSA256"
	case SIGNING_METHOD__RSA384:
		return "RSA384"
	case SIGNING_METHOD__RSA512:
		return "RSA512"
	case SIGNING_METHOD__RSAPSS256:
		return "RSAPSS256"
	case SIGNING_METHOD__RSAPSS384:
		return "RSAPSS384"
	case SIGNING_METHOD__RSAPSS512:
		return "RSAPSS512"
	}
}

func (v SigningMethod) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case SIGNING_METHOD_UNKNOWN:
		return ""
	case SIGNING_METHOD__ECDSA256:
		return "ECDSA256"
	case SIGNING_METHOD__ECDSA384:
		return "ECDSA384"
	case SIGNING_METHOD__ECDSA512:
		return "ECDSA512"
	case SIGNING_METHOD__HMAC256:
		return "HMAC256"
	case SIGNING_METHOD__HMAC384:
		return "HMAC384"
	case SIGNING_METHOD__HMAC512:
		return "HMAC512"
	case SIGNING_METHOD__RSA256:
		return "RSA256"
	case SIGNING_METHOD__RSA384:
		return "RSA384"
	case SIGNING_METHOD__RSA512:
		return "RSA512"
	case SIGNING_METHOD__RSAPSS256:
		return "RSAPSS256"
	case SIGNING_METHOD__RSAPSS384:
		return "RSAPSS384"
	case SIGNING_METHOD__RSAPSS512:
		return "RSAPSS512"
	}
}

func (v SigningMethod) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/depends/conf/jwt.SigningMethod"
}

func (v SigningMethod) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{SIGNING_METHOD__ECDSA256, SIGNING_METHOD__ECDSA384, SIGNING_METHOD__ECDSA512, SIGNING_METHOD__HMAC256, SIGNING_METHOD__HMAC384, SIGNING_METHOD__HMAC512, SIGNING_METHOD__RSA256, SIGNING_METHOD__RSA384, SIGNING_METHOD__RSA512, SIGNING_METHOD__RSAPSS256, SIGNING_METHOD__RSAPSS384, SIGNING_METHOD__RSAPSS512}
}

func (v SigningMethod) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidSigningMethod
	}
	return []byte(s), nil
}

func (v *SigningMethod) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseSigningMethodFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *SigningMethod) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = SigningMethod(i)
	return nil
}

func (v SigningMethod) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
