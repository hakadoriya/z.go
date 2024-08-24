package cliz

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/hakadoriya/z.go/stringz"
)

type marshalConfig struct {
	tagKey         string
	aliasKey       string
	envKey         string
	defaultKey     string
	requiredKey    string
	descriptionKey string
	hiddenKey      string
}

type MarshalOptionsOption interface {
	apply(c *marshalConfig)
}

type withMarshalOptionsOptionTagKey struct {
	tagKey string
}

func (w *withMarshalOptionsOptionTagKey) apply(c *marshalConfig) {
	c.tagKey = w.tagKey
}

func WithMarshalOptionsOptionTagKey(tagKey string) MarshalOptionsOption {
	return &withMarshalOptionsOptionTagKey{tagKey: tagKey}
}

type withMarshalOptionsOptionAliasKey struct {
	aliasKey string
}

func (w *withMarshalOptionsOptionAliasKey) apply(c *marshalConfig) {
	c.aliasKey = w.aliasKey
}

func WithMarshalOptionsOptionAliasKey(key string) MarshalOptionsOption {
	return &withMarshalOptionsOptionAliasKey{aliasKey: key}
}

type withMarshalOptionsOptionEnvKey struct {
	envKey string
}

func (w *withMarshalOptionsOptionEnvKey) apply(c *marshalConfig) {
	c.envKey = w.envKey
}

func WithMarshalOptionsOptionEnvKey(key string) MarshalOptionsOption {
	return &withMarshalOptionsOptionEnvKey{envKey: key}
}

type withMarshalOptionsOptionDefaultKey struct {
	defaultKey string
}

func (w *withMarshalOptionsOptionDefaultKey) apply(c *marshalConfig) {
	c.defaultKey = w.defaultKey
}

func WithMarshalOptionsOptionDefaultKey(key string) MarshalOptionsOption {
	return &withMarshalOptionsOptionDefaultKey{defaultKey: key}
}

type withMarshalOptionsOptionRequiredKey struct {
	requiredKey string
}

func (w *withMarshalOptionsOptionRequiredKey) apply(c *marshalConfig) {
	c.requiredKey = w.requiredKey
}

func WithMarshalOptionsOptionRequiredKey(key string) MarshalOptionsOption {
	return &withMarshalOptionsOptionRequiredKey{requiredKey: key}
}

type withMarshalOptionsOptionDescriptionKey struct {
	descriptionKey string
}

func (w *withMarshalOptionsOptionDescriptionKey) apply(c *marshalConfig) {
	c.descriptionKey = w.descriptionKey
}

func WithMarshalOptionsOptionDescriptionKey(key string) MarshalOptionsOption {
	return &withMarshalOptionsOptionDescriptionKey{descriptionKey: key}
}

type withMarshalOptionsOptionHiddenKey struct {
	hiddenKey string
}

func (w *withMarshalOptionsOptionHiddenKey) apply(c *marshalConfig) {
	c.hiddenKey = w.hiddenKey
}

func WithMarshalOptionsOptionHiddenKey(key string) MarshalOptionsOption {
	return &withMarshalOptionsOptionHiddenKey{hiddenKey: key}
}

// MarshalOptions generates the options from the structure pointer.
// This function reads the `cliz` tag set in the structure field and generates the options.
//
//nolint:funlen,cyclop,gocognit
func MarshalOptions(v interface{}, opts ...MarshalOptionsOption) (options []Option, err error) {
	cfg := &marshalConfig{
		tagKey:         DefaultTagKey,
		aliasKey:       DefaultAliasKey,
		envKey:         DefaultEnvKey,
		defaultKey:     DefaultDefaultKey,
		requiredKey:    DefaultRequiredKey,
		descriptionKey: DefaultDescriptionKey,
	}

	for _, opt := range opts {
		opt.apply(cfg)
	}

	ref := reflect.ValueOf(v)
	if ref.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("%T: %w", v, ErrInvalidType)
	}

	ref = ref.Elem()
	if ref.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%T: %w", v, ErrInvalidType)
	}

	typ := ref.Type()
	for i := range ref.NumField() {
		field := typ.Field(i)
		fieldValue := ref.Field(i)
		if !fieldValue.CanSet() {
			return nil, fmt.Errorf("type=%s: field=%s: tag=%s: %w", typ, field.Name, cfg.tagKey, ErrStructFieldCannotBeSet)
		}

		tagValue := stringz.TrimLeftSpace(field.Tag.Get(cfg.tagKey))
		Logger.Debug(fmt.Sprintf("type=%s: field=%s: tag=%s, tagValue=%s", typ, field.Name, cfg.tagKey, tagValue))
		if tagValue == "" {
			continue
		}

		optName, opts := parseTagValue(tagValue)
		Logger.Debug(fmt.Sprintf("type=%s: field=%s: tag=%s, optName=%s, opts=%v", typ, field.Name, cfg.tagKey, optName, opts))
		if optName == "" {
			return nil, fmt.Errorf("type=%s: field=%s: tag=%s: tagValue=%s: %w", typ, field.Name, cfg.tagKey, tagValue, ErrInvalidTagValue)
		}

		var (
			aliases     []string
			envKey      string
			required    bool
			description string
			hidden      bool
		)

		if key, ok := optsContainsAliasKey(cfg, opts); ok {
			aliases = append(aliases, key)
		}

		if key, ok := optsContainsEnvKey(cfg, opts); ok {
			envKey = key
		}

		defaultValue, defaultValueIsSet := optsContainsDefaultValue(cfg, opts)

		if optsContainsRequired(cfg, opts) {
			required = true
		}

		if optsContainsHidden(cfg, opts) {
			hidden = true
		}

		if key, ok := optsContainsDescription(cfg, opts); ok {
			description = key
		}

		const base, bitSize = 10, 64
		//nolint:exhaustive
		switch fieldValue.Kind() {
		case reflect.String: // string
			//nolint:exhaustruct
			options = append(options, &StringOption{
				Name:        optName,
				Aliases:     aliases,
				Env:         envKey,
				Default:     defaultValue,
				Required:    required,
				Description: description,
				Hidden:      hidden,
			})
		case reflect.Bool: // bool
			var defaultValueBool bool
			if defaultValueIsSet {
				defaultValueBool, err = strconv.ParseBool(defaultValue)
				if err != nil {
					return nil, fmt.Errorf("type=%s: field=%s: tag=%s: %s=%s: %w", typ, field.Name, cfg.tagKey, cfg.defaultKey, defaultValue, err)
				}
			}
			//nolint:exhaustruct
			options = append(options, &BoolOption{
				Name:        optName,
				Aliases:     aliases,
				Env:         envKey,
				Default:     defaultValueBool,
				Required:    required,
				Description: description,
			})
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64: // int, int8, int16, int32, int64
			var defaultValueInt64 int64
			if defaultValueIsSet {
				defaultValueInt64, err = strconv.ParseInt(defaultValue, base, bitSize)
				if err != nil {
					return nil, fmt.Errorf("type=%s: field=%s: tag=%s: %s=%s: %w", typ, field.Name, cfg.tagKey, cfg.defaultKey, defaultValue, err)
				}
			}
			//nolint:exhaustruct
			options = append(options, &Int64Option{
				Name:        optName,
				Aliases:     aliases,
				Env:         envKey,
				Default:     defaultValueInt64,
				Required:    required,
				Description: description,
			})
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64: // uint, uint8, uint16, uint32, uint64
			var defaultValueUint64 uint64
			if defaultValueIsSet {
				defaultValueUint64, err = strconv.ParseUint(defaultValue, base, bitSize)
				if err != nil {
					return nil, fmt.Errorf("type=%s: field=%s: tag=%s: %s=%s: %w", typ, field.Name, cfg.tagKey, cfg.defaultKey, defaultValue, err)
				}
			}
			//nolint:exhaustruct
			options = append(options, &Uint64Option{
				Name:        optName,
				Aliases:     aliases,
				Env:         envKey,
				Default:     defaultValueUint64,
				Required:    required,
				Description: description,
			})
		case reflect.Float32, reflect.Float64: // float32, float64
			var defaultValueFloat64 float64
			if defaultValueIsSet {
				defaultValueFloat64, err = strconv.ParseFloat(defaultValue, bitSize)
				if err != nil {
					return nil, fmt.Errorf("type=%s: field=%s: tag=%s: %s=%s: %w", typ, field.Name, cfg.tagKey, cfg.defaultKey, defaultValue, err)
				}
			}
			//nolint:exhaustruct
			options = append(options, &Float64Option{
				Name:        optName,
				Aliases:     aliases,
				Env:         envKey,
				Default:     defaultValueFloat64,
				Required:    required,
				Description: description,
			})
		default:
			return nil, fmt.Errorf("type=%s: field=%s: fieldType=%s: tag=%s: %w", typ, field.Name, field.Type, cfg.tagKey, ErrStructFieldTypeNotSupported)
		}
	}

	return options, nil
}

func optsContainsAliasKey(c *marshalConfig, opts []string) (aliasKey string, hasAlias bool) {
	for _, opt := range opts {
		if strings.HasPrefix(opt, c.aliasKey+"=") {
			return strings.CutPrefix(opt, c.aliasKey+"=")
		}
	}

	return "", false
}

func optsContainsEnvKey(c *marshalConfig, opts []string) (envKey string, hasEnv bool) {
	for _, opt := range opts {
		if strings.HasPrefix(opt, c.envKey+"=") {
			return strings.CutPrefix(opt, c.envKey+"=")
		}
	}

	return "", false
}

func optsContainsDefaultValue(c *marshalConfig, opts []string) (defaultValue string, hasDefault bool) {
	for _, opt := range opts {
		Logger.Debug("opt=" + opt)
		if strings.HasPrefix(opt, c.defaultKey+"=") {
			return strings.CutPrefix(opt, c.defaultKey+"=")
		}
	}

	return "", false
}

func optsContainsRequired(c *marshalConfig, opts []string) bool {
	for _, opt := range opts {
		if opt == c.requiredKey {
			return true
		}
	}

	return false
}

func optsContainsHidden(c *marshalConfig, opts []string) bool {
	for _, opt := range opts {
		if opt == c.hiddenKey {
			return true
		}
	}

	return false
}

func optsContainsDescription(c *marshalConfig, opts []string) (description string, hasDescription bool) {
	for _, opt := range opts {
		Logger.Debug("opt=" + opt)
		if strings.HasPrefix(opt, c.descriptionKey+"=") {
			return strings.CutPrefix(opt, c.descriptionKey+"=")
		}
	}

	return "", false
}
