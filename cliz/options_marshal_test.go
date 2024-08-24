package cliz

import (
	"testing"

	"github.com/hakadoriya/z.go/testingz/assertz"
	"github.com/hakadoriya/z.go/testingz/requirez"
)

type (
	testStructParseBool struct {
		Bool bool `cli:"bool-opt,alias=b,env=ENVZ_TEST_BOOL,default=INVALID"`
	}
	testStructParseInt struct {
		Int64 int64 `cli:"int64-opt,alias=i64,env=ENVZ_TEST_INT64,default=INVALID"`
	}
	testStructParseUint struct {
		Uint64 uint64 `cli:"uint64-opt,alias=u64,env=ENVZ_TEST_UINT64,default=INVALID"`
	}
	testStructParseFloat struct {
		Float64 float64 `cli:"float64-opt,alias=f64,env=ENVZ_TEST_FLOAT64,default=INVALID"`
	}
)

func TestMarshalOptions(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		actual, err := MarshalOptions(
			&testStruct{},
			WithMarshalOptionsOptionTagKey("cli"),
			WithMarshalOptionsOptionAliasKey("alias"),
			WithMarshalOptionsOptionEnvKey("env"),
			WithMarshalOptionsOptionDefaultKey("default"),
			WithMarshalOptionsOptionRequiredKey("required"),
			WithMarshalOptionsOptionDescriptionKey("description"),
			WithMarshalOptionsOptionHiddenKey("hidden"),
		)
		requirez.NoError(t, err)
		assertz.Equal(t, 6, len(actual))
		// string
		assertz.Equal(t, "string-opt", actual[0].(*StringOption).Name)
		assertz.Equal(t, "s", actual[0].(*StringOption).Aliases[0])
		assertz.Equal(t, "ENVZ_TEST_STRING", actual[0].(*StringOption).Env)
		assertz.Equal(t, "defaultString", actual[0].(*StringOption).Default)
		assertz.Equal(t, true, actual[0].(*StringOption).Required)
		assertz.Equal(t, "string description", actual[0].(*StringOption).Description)
		// bool
		assertz.Equal(t, "bool-opt", actual[1].(*BoolOption).Name)
		assertz.Equal(t, "b", actual[1].(*BoolOption).Aliases[0])
		assertz.Equal(t, "ENVZ_TEST_BOOL", actual[1].(*BoolOption).Env)
		assertz.Equal(t, true, actual[1].(*BoolOption).Default)
		assertz.Equal(t, false, actual[1].(*BoolOption).Required)
		assertz.Equal(t, "bool description", actual[1].(*BoolOption).Description)
		// int64
		assertz.Equal(t, "int64-opt", actual[2].(*Int64Option).Name)
		assertz.Equal(t, "i64", actual[2].(*Int64Option).Aliases[0])
		assertz.Equal(t, "ENVZ_TEST_INT64", actual[2].(*Int64Option).Env)
		assertz.Equal(t, int64(128), actual[2].(*Int64Option).Default)
		assertz.Equal(t, false, actual[2].(*Int64Option).Required)
		assertz.Equal(t, "int64 description", actual[2].(*Int64Option).Description)
		// uint64
		assertz.Equal(t, "uint64-opt", actual[3].(*Uint64Option).Name)
		assertz.Equal(t, "u64", actual[3].(*Uint64Option).Aliases[0])
		assertz.Equal(t, "ENVZ_TEST_UINT64", actual[3].(*Uint64Option).Env)
		assertz.Equal(t, uint64(256), actual[3].(*Uint64Option).Default)
		assertz.Equal(t, false, actual[3].(*Uint64Option).Required)
		assertz.Equal(t, "uint64 description", actual[3].(*Uint64Option).Description)
		// float64
		assertz.Equal(t, "float64-opt", actual[4].(*Float64Option).Name)
		assertz.Equal(t, "f64", actual[4].(*Float64Option).Aliases[0])
		assertz.Equal(t, "ENVZ_TEST_FLOAT64", actual[4].(*Float64Option).Env)
		assertz.Equal(t, 3.141592, actual[4].(*Float64Option).Default)
		assertz.Equal(t, false, actual[4].(*Float64Option).Required)
		assertz.Equal(t, "float64 description", actual[4].(*Float64Option).Description)
		// hidden
		assertz.Equal(t, "hidden-opt", actual[5].(*StringOption).Name)
		assertz.True(t, actual[5].(*StringOption).Hidden)
	})

	t.Run("error,Ptr,ErrInvalidType", func(t *testing.T) {
		t.Parallel()

		var v int
		_, err := MarshalOptions(v)
		requirez.ErrorIs(t, err, ErrInvalidType)
	})

	t.Run("error,Struct,ErrInvalidType", func(t *testing.T) {
		t.Parallel()

		var v int
		_, err := MarshalOptions(&v)
		requirez.ErrorIs(t, err, ErrInvalidType)
	})

	t.Run("error,ErrStructFieldCannotBeSet", func(t *testing.T) {
		t.Parallel()

		var v testStructCannotSet
		_, err := MarshalOptions(&v)
		requirez.ErrorIs(t, err, ErrStructFieldCannotBeSet)
	})

	t.Run("error,ErrInvalidTagValue", func(t *testing.T) {
		t.Parallel()

		var v testStructInvalidTagValue
		_, err := MarshalOptions(&v)
		requirez.ErrorIs(t, err, ErrInvalidTagValue)
	})

	t.Run("error,default,strconv.ParseBool", func(t *testing.T) {
		t.Parallel()

		var v testStructParseBool
		_, err := MarshalOptions(&v)
		requirez.ErrorContains(t, err, `strconv.ParseBool: parsing "INVALID": invalid syntax`)
	})

	t.Run("error,default,strconv.ParseInt", func(t *testing.T) {
		t.Parallel()

		var v testStructParseInt
		_, err := MarshalOptions(&v)
		requirez.ErrorContains(t, err, `strconv.ParseInt: parsing "INVALID": invalid syntax`)
	})

	t.Run("error,default,strconv.ParseUint", func(t *testing.T) {
		t.Parallel()

		var v testStructParseUint
		_, err := MarshalOptions(&v)
		requirez.ErrorContains(t, err, `strconv.ParseUint: parsing "INVALID": invalid syntax`)
	})

	t.Run("error,default,strconv.ParseFloat", func(t *testing.T) {
		t.Parallel()

		var v testStructParseFloat
		_, err := MarshalOptions(&v)
		requirez.ErrorContains(t, err, `strconv.ParseFloat: parsing "INVALID": invalid syntax`)
	})

	t.Run("error,ErrFieldTypeNotSupported", func(t *testing.T) {
		t.Parallel()

		var v testStructFieldTypeNotSupported
		_, err := MarshalOptions(&v)
		requirez.ErrorIs(t, err, ErrStructFieldTypeNotSupported)
	})
}
