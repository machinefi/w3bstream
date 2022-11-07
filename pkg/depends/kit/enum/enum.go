package enum

import "strconv"

type IntStringerEnum interface {
	TypeName() string
	Int() int
	String() string
	Label() string
	ConstValues() []IntStringerEnum
}

type Enum = IntStringerEnum

type ValueOffset interface {
	Offset() int
}

func ScanEnum(src interface{}, offset int) (int, error) {
	return ScanIntEnumStringer(src, offset)
}

func ScanIntEnumStringer(src interface{}, offset int) (int, error) {
	n, err := toInteger(src, offset)
	if err != nil {
		return offset, err
	}
	return n - offset, nil
}

func toInteger(src interface{}, dft int) (int, error) {
	switch v := src.(type) {
	case []byte:
		if len(v) > 0 {
			i, err := strconv.ParseInt(string(v), 10, 64)
			if err != nil {
				return dft, err
			}
			return int(i), err
		}
		return dft, nil
	case string:
		if v != "" {
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return dft, err
			}
			return int(i), err
		}
		return dft, nil
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case nil:
		return dft, nil
	default:
		return dft, nil
	}
}
