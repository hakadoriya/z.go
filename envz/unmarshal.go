package envz

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/hakadoriya/z.go/envz/envcorez"
)

// pkgInterface is a entry point for mocking.
type pkgInterface interface {
	Getenv(key string) string
}

type pkg struct {
	GetenvFunc func(key string) string
}

func (s *pkg) Getenv(key string) string {
	return s.GetenvFunc(key)
}

type unmarshalConfig struct {
	tagKey      string
	requiredKey string
	defaultKey  string
}

type UnmarshalOption interface {
	apply(c *unmarshalConfig)
}

type withUnmarshalOptionTagKey struct {
	tagKey string
}

func (w *withUnmarshalOptionTagKey) apply(c *unmarshalConfig) {
	c.tagKey = w.tagKey
}

func WithUnmarshalOptionTagKey(tagKey string) UnmarshalOption {
	return &withUnmarshalOptionTagKey{tagKey: tagKey}
}

type withUnmarshalOptionRequiredKey struct {
	requiredKey string
}

func (w *withUnmarshalOptionRequiredKey) apply(c *unmarshalConfig) {
	c.requiredKey = w.requiredKey
}

func WithUnmarshalOptionRequiredKey(key string) UnmarshalOption {
	return &withUnmarshalOptionRequiredKey{requiredKey: key}
}

type withUnmarshalOptionDefaultKey struct {
	defaultKey string
}

func (w *withUnmarshalOptionDefaultKey) apply(c *unmarshalConfig) {
	c.defaultKey = w.defaultKey
}

func WithUnmarshalOptionDefaultKey(key string) UnmarshalOption {
	return &withUnmarshalOptionDefaultKey{defaultKey: key}
}

// Unmarshal sets the value read from the environment variable to the field of the passed structure pointer.
// This function reads the value from the environment variable according to the `env` tag set in the structure field.
// The value of the `env` tag specifies the key of the environment variable.
// If the value of the tag ends with `,required`, an error is returned if the environment variable is not found.
func Unmarshal(v interface{}, opts ...UnmarshalOption) error {
	return unmarshal(&pkg{GetenvFunc: os.Getenv}, v, opts...)
}

//nolint:funlen,gocognit,cyclop
func unmarshal(iface pkgInterface, v interface{}, opts ...UnmarshalOption) error {
	c := &unmarshalConfig{
		tagKey:      envcorez.TagKey,
		requiredKey: envcorez.RequiredKey,
		defaultKey:  envcorez.DefaultKey,
	}

	for _, opt := range opts {
		opt.apply(c)
	}

	ref := reflect.ValueOf(v)
	if ref.Kind() != reflect.Ptr {
		return fmt.Errorf("%T: %w", v, envcorez.ErrInvalidType)
	}

	ref = ref.Elem()
	if ref.Kind() != reflect.Struct {
		return fmt.Errorf("%T: %w", v, envcorez.ErrInvalidType)
	}

	typ := ref.Type()
	for i := range ref.NumField() {
		field := typ.Field(i)
		fieldValue := ref.Field(i)
		if !fieldValue.CanSet() {
			return fmt.Errorf("field=%s: tag=%s: %w", field.Name, c.tagKey, envcorez.ErrStructFieldCannotBeSet)
		}

		tagValue := field.Tag.Get(c.tagKey)
		envcorez.Logger.Debug(fmt.Sprintf("tagKey=%s, tagValue=%s", c.tagKey, tagValue))
		if tagValue == "" {
			continue
		}

		envKey, opts := parseTagValue(tagValue)
		envcorez.Logger.Debug(fmt.Sprintf("tagKey=%s, envKey=%s, opts=%v", c.tagKey, envKey, opts))
		if envKey == "" {
			return fmt.Errorf("field=%s: tag=%s: tagValue=%s: %w", field.Name, c.tagKey, tagValue, envcorez.ErrInvalidTagValue)
		}

		required := false
		if optsContainsRequired(c, opts) {
			required = true
		}

		envValue := iface.Getenv(envKey)
		if envValue == "" {
			if required {
				return fmt.Errorf("field=%s: tag=%s: %s: %w", field.Name, c.tagKey, envKey, envcorez.ErrRequiredEnvironmentVariableNotFound)
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
				return fmt.Errorf("field=%s: tag=%s: %T: %w", field.Name, c.tagKey, v, envcorez.ErrStructFieldTypeNotSupported)
			}
		default:
			return fmt.Errorf("field=%s: tag=%s: %T: %w", field.Name, c.tagKey, v, envcorez.ErrStructFieldTypeNotSupported)
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

func optsContainsRequired(c *unmarshalConfig, opts []string) bool {
	for _, opt := range opts {
		if opt == c.requiredKey {
			return true
		}
	}

	return false
}

func optsContainsDefault(c *unmarshalConfig, opts []string) (defaultValue string, hasDefault bool) {
	for _, opt := range opts {
		envcorez.Logger.Debug("opt=" + opt)
		if strings.HasPrefix(opt, c.defaultKey+"=") {
			return strings.CutPrefix(opt, c.defaultKey+"=")
		}
	}

	return "", false
}
