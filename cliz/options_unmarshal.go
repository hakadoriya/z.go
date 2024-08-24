package cliz

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hakadoriya/z.go/slicez"
	"github.com/hakadoriya/z.go/stringz"
)

type unmarshalConfig struct {
	tagKey string
}

type UnmarshalOptionsOption interface {
	apply(c *unmarshalConfig)
}

type withUnmarshalOptionsOptionTagKey struct {
	tagKey string
}

func (w *withUnmarshalOptionsOptionTagKey) apply(c *unmarshalConfig) {
	c.tagKey = w.tagKey
}

func WithUnmarshalOptionsOptionTagKey(tagKey string) UnmarshalOptionsOption {
	return &withUnmarshalOptionsOptionTagKey{tagKey: tagKey}
}

// UnmarshalOptions sets the value read from (*Command).Options to the field of the passed structure pointer.
// This function reads the value from (*Command).Options according to the `cli` tag set in the structure field.
// The value of the `cli` tag specifies the key of the cliz.Option name.
//
//nolint:funlen,cyclop
func UnmarshalOptions(c *Command, v interface{}, opts ...UnmarshalOptionsOption) error {
	cfg := &unmarshalConfig{
		tagKey: DefaultTagKey,
	}

	for _, opt := range opts {
		opt.apply(cfg)
	}

	ref := reflect.ValueOf(v)
	if ref.Kind() != reflect.Ptr {
		return fmt.Errorf("%T: %w", v, ErrInvalidType)
	}

	ref = ref.Elem()
	if ref.Kind() != reflect.Struct {
		return fmt.Errorf("%T: %w", v, ErrInvalidType)
	}

	typ := ref.Type()
	for i := range ref.NumField() {
		field := typ.Field(i)
		fieldValue := ref.Field(i)
		if !fieldValue.CanSet() {
			return fmt.Errorf("field=%s: tag=%s: %w", field.Name, cfg.tagKey, ErrStructFieldCannotBeSet)
		}

		tagValue := stringz.TrimLeftSpace(field.Tag.Get(cfg.tagKey))
		Logger.Debug(fmt.Sprintf("tagKey=%s, tagValue=%s", cfg.tagKey, tagValue))
		if tagValue == "" {
			continue
		}

		optName, opts := parseTagValue(tagValue)
		Logger.Debug(fmt.Sprintf("tagKey=%s, envKey=%s, opts=%v", cfg.tagKey, optName, opts))
		if optName == "" {
			return fmt.Errorf("field=%s: tag=%s: tagValue=%s: %w", field.Name, cfg.tagKey, tagValue, ErrInvalidTagValue)
		}

		//nolint:exhaustive
		switch fieldValue.Kind() {
		case reflect.String: // string
			optValue, err := c.GetOptionString(optName)
			if err != nil {
				return fmt.Errorf("field=%s: tag=%s: cmd.GetOptionString: %w", field.Name, cfg.tagKey, err)
			}
			fieldValue.SetString(optValue)
		case reflect.Bool: // bool
			optValue, err := c.GetOptionBool(optName)
			if err != nil {
				return fmt.Errorf("field=%s: tag=%s: cmd.GetOptionBool: %w", field.Name, cfg.tagKey, err)
			}
			fieldValue.SetBool(optValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64: // int, int8, int16, int32, int64
			optValue, err := c.GetOptionInt64(optName)
			if err != nil {
				return fmt.Errorf("field=%s: tag=%s: cmd.GetOptionInt64: %w", field.Name, cfg.tagKey, err)
			}
			fieldValue.SetInt(optValue)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64: // uint, uint8, uint16, uint32, uint64
			optValue, err := c.getOptionUint64(optName)
			if err != nil {
				return fmt.Errorf("field=%s: tag=%s: cmd.GetOptionUint64: %w", field.Name, cfg.tagKey, err)
			}
			fieldValue.SetUint(optValue)
		case reflect.Float32, reflect.Float64: // float32, float64
			optValue, err := c.GetOptionFloat64(optName)
			if err != nil {
				return fmt.Errorf("field=%s: tag=%s: cmd.GetOptionFloat64: %w", field.Name, cfg.tagKey, err)
			}
			fieldValue.SetFloat(optValue)
		default:
			return fmt.Errorf("field=%s: tag=%s: %T: %w", field.Name, cfg.tagKey, v, ErrStructFieldTypeNotSupported)
		}
	}

	return nil
}

func parseTagValue(tagValue string) (envKey string, opts []string) {
	if i := strings.Index(tagValue, ","); i != -1 {
		envKey = tagValue[:i]
		opts = slicez.Map(strings.Split(tagValue[i+1:], ","), func(input string) string {
			return stringz.TrimLeftSpace(input)
		})
	} else {
		envKey = tagValue
	}

	return
}
