package envz

import (
	"errors"
	"strings"
	"testing"
)

var testPkg = &pkg{GetenvFunc: func(key string) string {
	switch key {
	case "ENVZ_TEST_STRING":
		return "hello"
	case "ENVZ_TEST_BYTES":
		return "world"
	case "ENVZ_TEST_STRING_SLICE":
		return "hello,world"
	case "ENVZ_TEST_BOOL":
		return "true"
	case "ENVZ_TEST_INT64":
		return "128"
	case "ENVZ_TEST_UINT64":
		return "256"
	case "ENVZ_TEST_FLOAT64":
		return "3.141592"
	default:
		return ""
	}
}}

type testStruct struct {
	Default      string   `env:"ENVZ_TEST_DEFAULT,default=defaultValue"`
	Default2     string   `env:"ENVZ_TEST_DEFAULT2,default2=default2Value"`
	String       string   `env:"ENVZ_TEST_STRING,default=world"`
	Bytes        []byte   `env:"ENVZ_TEST_BYTES"`
	StringSlice  []string `env:"ENVZ_TEST_STRING_SLICE"`
	StringSlice2 []string `env:"ENVZ_TEST_STRING_SLICE2,default=\"hello2,world2\""`
	Bool         bool     `env:"ENVZ_TEST_BOOL"`
	Int64        int64    `env:"ENVZ_TEST_INT64"`
	Uint64       uint64   `env:"ENVZ_TEST_UINT64"`
	Float64      float64  `env:"ENVZ_TEST_FLOAT64"`

	TagNotSet string
	TagEnv2   string `env2:"ENVZ_TEST_STRING"`
	EnvNotSet string `env:"ENVZ_TEST_ENV_NOT_SET"`
}

type testStructRequired struct {
	String string `env:"ENVZ_TEST_STRING,required"`
}

type testStructRequired2 struct {
	String string `env:"ENVZ_TEST_STRING,required2"`
}

type testStructCannotSet struct {
	cannotSet string `env:"ENVZ_TEST_STRING"`
}

type testStructInvalidTagValue struct {
	String string `env:",invalid"`
}

type testStructFieldTypeNotSupported struct {
	NotSupported struct{} `env:"ENVZ_TEST_STRING"`
}

type testStructFieldSliceTypeNotSupported struct {
	NotSupported []struct{} `env:"ENVZ_TEST_STRING"`
}

func TestMarshal(t *testing.T) {
	t.Parallel()

	t.Run("success,default", func(t *testing.T) {
		t.Parallel()

		var v testStruct
		err := Unmarshal(&v)
		if err != nil {
			t.Errorf("❌: err != nil: %+v", err)
		}
		const expected = "defaultValue"
		actual := v.Default
		if expected != actual {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, actual)
		}
	})

	t.Run("success,default2", func(t *testing.T) {
		t.Parallel()

		var v testStruct
		err := Unmarshal(&v, WithUnmarshalOptionDefaultKey("default2"))
		if err != nil {
			t.Errorf("❌: err != nil: %+v", err)
		}
		const expected = "default2Value"
		actual := v.Default2
		if expected != actual {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, actual)
		}
	})

	t.Run("error,ErrInvalidType,int", func(t *testing.T) {
		t.Parallel()

		var v int
		err := Unmarshal(v)
		if !errors.Is(err, ErrInvalidType) {
			t.Errorf("❌: !errors.Is(err, ErrInvalidType): %+v", err)
		}
	})

	t.Run("error,ErrInvalidType,intptr", func(t *testing.T) {
		t.Parallel()

		var v string
		err := Unmarshal(&v)
		if !errors.Is(err, ErrInvalidType) {
			t.Errorf("❌: !errors.Is(err, ErrInvalidType): %+v", err)
		}
	})

	t.Run("error,ErrStructFieldCannotBeSet", func(t *testing.T) {
		t.Parallel()

		var v testStructCannotSet
		err := Unmarshal(&v)
		if !errors.Is(err, ErrStructFieldCannotBeSet) {
			t.Errorf("❌: !errors.Is(err, ErrStructFieldCannotBeSet): %+v", err)
		}
	})

	t.Run("error,ErrInvalidTagValue", func(t *testing.T) {
		t.Parallel()

		var v testStructInvalidTagValue
		err := Unmarshal(&v)
		if !errors.Is(err, ErrInvalidTagValue) {
			t.Errorf("❌: !errors.Is(err, ErrInvalidTagValue): %+v", err)
		}
	})

	t.Run("error,required,ErrRequiredEnvironmentVariableNotFound", func(t *testing.T) {
		t.Parallel()

		var v testStructRequired
		err := Unmarshal(&v)
		if !errors.Is(err, ErrRequiredEnvironmentVariableNotFound) {
			t.Errorf("❌: !errors.Is(err, ErrRequiredEnvironmentVariableNotFound): %+v", err)
		}
	})

	t.Run("error,required2,ErrRequiredEnvironmentVariableNotFound", func(t *testing.T) {
		t.Parallel()

		var v testStructRequired2
		err := Unmarshal(&v, WithUnmarshalOptionRequiredKey("required2"))
		if !errors.Is(err, ErrRequiredEnvironmentVariableNotFound) {
			t.Errorf("❌: !errors.Is(err, ErrRequiredEnvironmentVariableNotFound): %+v", err)
		}
	})
}

func Test_marshal(t *testing.T) {
	t.Parallel()

	t.Run("success,string", func(t *testing.T) {
		t.Parallel()

		var v testStruct
		err := unmarshal(testPkg, &v)
		if err != nil {
			t.Errorf("❌: err != nil: %+v", err)
		}
		const expected = "hello"
		actual := v.String
		if expected != actual {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, actual)
		}
	})

	t.Run("success,[]byte", func(t *testing.T) {
		t.Parallel()

		var v testStruct
		err := unmarshal(testPkg, &v)
		if err != nil {
			t.Errorf("❌: err != nil: %+v", err)
		}
		const expected = "world"
		actual := string(v.Bytes)
		if expected != actual {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, actual)
		}
	})

	t.Run("success,[]string", func(t *testing.T) {
		t.Parallel()

		var v testStruct
		err := unmarshal(testPkg, &v)
		if err != nil {
			t.Errorf("❌: err != nil: %+v", err)
		}
		const expected1 = "hello"
		actual1 := v.StringSlice[0]
		if expected1 != actual1 {
			t.Errorf("❌: expected(%s) != actual(%s)", expected1, actual1)
		}
		const expected2 = "world"
		actual2 := v.StringSlice[1]
		if expected2 != actual2 {
			t.Errorf("❌: expected(%s) != actual(%s)", expected2, actual2)
		}
		const expected3 = "hello2"
		actual3 := v.StringSlice2[0]
		if expected3 != actual3 {
			t.Errorf("❌: expected(%s) != actual(%s)", expected3, actual3)
		}
		const expected4 = "world2"
		actual4 := v.StringSlice2[1]
		if expected4 != actual4 {
			t.Errorf("❌: expected(%s) != actual(%s)", expected4, actual4)
		}
	})

	t.Run("success,bool", func(t *testing.T) {
		t.Parallel()

		var v testStruct
		err := unmarshal(testPkg, &v)
		if err != nil {
			t.Errorf("❌: err != nil: %+v", err)
		}
		const expected = true
		actual := v.Bool
		if expected != actual {
			t.Errorf("❌: expected(%t) != actual(%t)", expected, actual)
		}
	})

	t.Run("success,int64", func(t *testing.T) {
		t.Parallel()

		var v testStruct
		err := unmarshal(testPkg, &v)
		if err != nil {
			t.Errorf("❌: err != nil: %+v", err)
		}
		const expected = int64(128)
		actual := v.Int64
		if expected != actual {
			t.Errorf("❌: expected(%d) != actual(%d)", expected, actual)
		}
	})

	t.Run("success,uint64", func(t *testing.T) {
		t.Parallel()

		var v testStruct
		err := unmarshal(testPkg, &v)
		if err != nil {
			t.Errorf("❌: err != nil: %+v", err)
		}
		const expected = uint64(256)
		actual := v.Uint64
		if expected != actual {
			t.Errorf("❌: expected(%d) != actual(%d)", expected, actual)
		}
	})

	t.Run("success,float64", func(t *testing.T) {
		t.Parallel()

		var v testStruct
		err := unmarshal(testPkg, &v)
		if err != nil {
			t.Errorf("❌: err != nil: %+v", err)
		}
		const expected = 3.141592
		actual := v.Float64
		if expected != actual {
			t.Errorf("❌: expected(%f) != actual(%f)", expected, actual)
		}
	})

	t.Run("success,env2", func(t *testing.T) {
		t.Parallel()

		var v testStruct
		err := unmarshal(testPkg, &v, WithUnmarshalOptionTagKey("env2"))
		if err != nil {
			t.Errorf("❌: err != nil: %+v", err)
		}

		const expected = "hello"
		actual := v.TagEnv2
		if expected != actual {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, actual)
		}
	})

	t.Run("error,reflect.Bool", func(t *testing.T) {
		t.Parallel()

		type testStruct struct {
			Bool bool `env:"ENVZ_TEST_STRING"`
		}
		var v testStruct
		err := unmarshal(testPkg, &v)
		const expected = `field=Bool: tag=env: strconv.ParseBool: strconv.ParseBool: parsing "hello": invalid syntax`
		if err == nil || !strings.Contains(err.Error(), expected) {
			t.Errorf("❌: !strings.Contains(err.Error(), `%s`): %+v", expected, err)
		}
	})

	t.Run("error,reflect.Int64", func(t *testing.T) {
		t.Parallel()

		type testStruct struct {
			Int64 int64 `env:"ENVZ_TEST_STRING"`
		}
		var v testStruct
		err := unmarshal(testPkg, &v)
		const expected = `field=Int64: tag=env: strconv.ParseInt: strconv.ParseInt: parsing "hello": invalid syntax`
		if err == nil || !strings.Contains(err.Error(), expected) {
			t.Errorf("❌: !strings.Contains(err.Error(), `%s`): %+v", expected, err)
		}
	})

	t.Run("error,reflect.Uint64", func(t *testing.T) {
		t.Parallel()

		type testStruct struct {
			Uint64 uint64 `env:"ENVZ_TEST_STRING"`
		}
		var v testStruct
		err := unmarshal(testPkg, &v)
		const expected = `field=Uint64: tag=env: strconv.ParseUint: strconv.ParseUint: parsing "hello": invalid syntax`
		if err == nil || !strings.Contains(err.Error(), expected) {
			t.Errorf("❌: !strings.Contains(err.Error(), `%s`): %+v", expected, err)
		}
	})

	t.Run("error,reflect.Float64", func(t *testing.T) {
		t.Parallel()

		type testStruct struct {
			Float64 float64 `env:"ENVZ_TEST_STRING"`
		}
		var v testStruct
		err := unmarshal(testPkg, &v)
		const expected = `field=Float64: tag=env: strconv.ParseFloat: strconv.ParseFloat: parsing "hello": invalid syntax`
		if err == nil || !strings.Contains(err.Error(), expected) {
			t.Errorf("❌: !strings.Contains(err.Error(), `%s`): %+v", expected, err)
		}
	})

	t.Run("error,ErrStructFieldTypeNotSupported,struct", func(t *testing.T) {
		t.Parallel()

		var v testStructFieldTypeNotSupported
		err := unmarshal(testPkg, &v)
		if !errors.Is(err, ErrStructFieldTypeNotSupported) {
			t.Errorf("❌: !errors.Is(err, ErrStructFieldTypeNotSupported): %+v", err)
		}
	})

	t.Run("error,ErrStructFieldTypeNotSupported,slice", func(t *testing.T) {
		t.Parallel()

		var v testStructFieldSliceTypeNotSupported
		err := unmarshal(testPkg, &v)
		if !errors.Is(err, ErrStructFieldTypeNotSupported) {
			t.Errorf("❌: !errors.Is(err, ErrStructFieldTypeNotSupported): %+v", err)
		}
	})
}
