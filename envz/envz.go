package envz

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/hakadoriya/z.go/envz/envcorez"
)

var (
	ErrInvalidType                         = errors.New("invalid type; must be a pointer to a struct")
	ErrStructFieldCannotBeSet              = errors.New("struct field cannot be set; unexported field or field is not settable")
	ErrInvalidTagValue                     = errors.New("invalid tag value")
	ErrStructFieldTypeNotSupported         = errors.New("struct field type not supported")
	ErrRequiredEnvironmentVariableNotFound = errors.New("required environment variable not found")
)

// pkg is a entry point for mocking.
type pkg struct {
	GetenvFunc func(key string) string
}

func (s *pkg) Getenv(key string) string {
	return s.GetenvFunc(key)
}

type marshalConfig struct {
	tagKey            string
	requiredOptionKey string
	defaultOptionKey  string
}

type MarshalOption interface {
	apply(c *marshalConfig)
}

type withMarshalOptionTagKey struct {
	tagKey string
}

func (w *withMarshalOptionTagKey) apply(c *marshalConfig) {
	c.tagKey = w.tagKey
}

func WithMarshalOptionTagKey(tagKey string) MarshalOption {
	return &withMarshalOptionTagKey{tagKey: tagKey}
}

type withMarshalOptionRequiredKey struct {
	requiredOptionKey string
}

func (w *withMarshalOptionRequiredKey) apply(c *marshalConfig) {
	c.requiredOptionKey = w.requiredOptionKey
}

func WithMarshalOptionRequiredOptionKey(key string) MarshalOption {
	return &withMarshalOptionRequiredKey{requiredOptionKey: key}
}

type withMarshalOptionDefaultKey struct {
	defaultOptionKey string
}

func (w *withMarshalOptionDefaultKey) apply(c *marshalConfig) {
	c.defaultOptionKey = w.defaultOptionKey
}

func WithMarshalOptionDefaultOptionKey(key string) MarshalOption {
	return &withMarshalOptionDefaultKey{defaultOptionKey: key}
}

// Marshal sets the value read from the environment variable to the field of the passed structure pointer.
// This function reads the value from the environment variable according to the `env` tag set in the structure field.
// The value of the `env` tag specifies the key of the environment variable.
// If the value of the tag ends with `,required`, an error is returned if the environment variable is not found.
func Marshal(v interface{}, opts ...MarshalOption) error {
	return marshal(&pkg{GetenvFunc: os.Getenv}, v, opts...)
}

//nolint:funlen,gocognit,cyclop
func marshal(
	iface interface {
		Getenv(key string) string
	},
	v interface{},
	opts ...MarshalOption,
) error {
	c := &marshalConfig{
		tagKey:            envcorez.TagKey,
		requiredOptionKey: envcorez.OptionKeyRequired,
		defaultOptionKey:  envcorez.OptionKeyDefault,
	}

	for _, opt := range opts {
		opt.apply(c)
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
			return fmt.Errorf("field=%s: tag=%s: %w", field.Name, c.tagKey, ErrStructFieldCannotBeSet)
		}

		tagValue := field.Tag.Get(c.tagKey)
		envcorez.Logger.Debug(fmt.Sprintf("tagKey=%s, tagValue=%s", c.tagKey, tagValue))
		if tagValue == "" {
			continue
		}

		envKey, opts := parseTagValue(tagValue)
		envcorez.Logger.Debug(fmt.Sprintf("tagKey=%s, envKey=%s, opts=%v", c.tagKey, envKey, opts))
		if envKey == "" {
			return fmt.Errorf("field=%s: tag=%s: tagValue=%s: %w", field.Name, c.tagKey, tagValue, ErrInvalidTagValue)
		}

		required := false
		if optsContainsRequired(c, opts) {
			required = true
		}

		envValue := iface.Getenv(envKey)
		if envValue == "" {
			if required {
				return fmt.Errorf("field=%s: tag=%s: %s: %w", field.Name, c.tagKey, envKey, ErrRequiredEnvironmentVariableNotFound)
			}

			defaultValue, hasDefault := optsContainsDefault(c, opts)
			if !hasDefault {
				// If the environment variable is not found and there is no default value, skip setting the field.
				continue

			}

			envValue = defaultValue
		}

		const base, bitSize = 10, 64
		//nolint:exhaustive
		switch fieldValue.Kind() {
		case reflect.String: // string
			fieldValue.SetString(envValue)
		case reflect.Bool: // bool
			envBool, err := strconv.ParseBool(envValue)
			if err != nil {
				return fmt.Errorf("field=%s: tag=%s: strconv.ParseBool: %w", field.Name, c.tagKey, err)
			}
			fieldValue.SetBool(envBool)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64: // int, int8, int16, int32, int64
			envInt64, err := strconv.ParseInt(envValue, base, bitSize)
			if err != nil {
				return fmt.Errorf("field=%s: tag=%s: strconv.ParseInt: %w", field.Name, c.tagKey, err)
			}
			fieldValue.SetInt(envInt64)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64: // uint, uint8, uint16, uint32, uint64
			envUint, err := strconv.ParseUint(envValue, base, bitSize)
			if err != nil {
				return fmt.Errorf("field=%s: tag=%s: strconv.ParseUint: %w", field.Name, c.tagKey, err)
			}
			fieldValue.SetUint(envUint)
		case reflect.Float32, reflect.Float64: // float32, float64
			envFloat, err := strconv.ParseFloat(envValue, bitSize)
			if err != nil {
				return fmt.Errorf("field=%s: tag=%s: strconv.ParseFloat: %w", field.Name, c.tagKey, err)
			}
			fieldValue.SetFloat(envFloat)
		case reflect.Slice:
			//nolint:exhaustive
			switch fieldValue.Type().Elem().Kind() {
			case reflect.Uint8: // []byte
				fieldValue.SetBytes([]byte(envValue))
			case reflect.String: // []string
				fieldValue.Set(reflect.ValueOf(strings.Split(envValue, ",")))
			default:
				return fmt.Errorf("field=%s: tag=%s: %T: %w", field.Name, c.tagKey, v, ErrStructFieldTypeNotSupported)
			}
		default:
			return fmt.Errorf("field=%s: tag=%s: %T: %w", field.Name, c.tagKey, v, ErrStructFieldTypeNotSupported)
		}
	}

	return nil
}

func parseTagValue(tagValue string) (envKey string, opts []string) {
	if i := strings.Index(tagValue, ","); i != -1 {
		envKey = tagValue[:i]
		opts = strings.Split(tagValue[i+1:], ",")
	} else {
		envKey = tagValue
	}

	return
}

func optsContainsRequired(c *marshalConfig, opts []string) bool {
	for _, opt := range opts {
		if opt == c.requiredOptionKey {
			return true
		}
	}

	return false
}

func optsContainsDefault(c *marshalConfig, opts []string) (defaultValue string, hasDefault bool) {
	for _, opt := range opts {
		envcorez.Logger.Debug(fmt.Sprintf("opt=%s", opt))
		if strings.HasPrefix(opt, c.defaultOptionKey+"=") {
			return strings.CutPrefix(opt, c.defaultOptionKey+"=")
		}
	}

	return "", false
}
