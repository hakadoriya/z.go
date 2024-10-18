package cliz

import (
	"context"
	"testing"

	"github.com/hakadoriya/z.go/testingz/requirez"
)

type (
	testStruct struct {
		String  string  `cli:"string-opt,alias=s,env=ENVZ_TEST_STRING,default=defaultString,required,description=string description"`
		Bool    bool    `cli:"bool-opt,alias=b,env=ENVZ_TEST_BOOL,default=true,description=bool description"`
		Int64   int64   `cli:"int64-opt,alias=i64,env=ENVZ_TEST_INT64,default=128,description=int64 description"`
		Uint64  uint64  `cli:"uint64-opt,alias=u64,env=ENVZ_TEST_UINT64,default=256,description=uint64 description"`
		Float64 float64 `cli:"float64-opt,alias=f64,env=ENVZ_TEST_FLOAT64,default=3.141592,description=float64 description"`
		Hidden  string  `cli:"hidden-opt,hidden"`

		TagNotSet string
	}
	testStructCannotSet struct {
		cannotSet string `cli:"string-opt"`
	}
	testStructInvalidTagValue struct {
		String string `cli:",invalid"`
	}
	testStructOptionStringNotFound struct {
		String string `cli:"not-found"`
	}
	testStructOptionBoolNotFound struct {
		Bool bool `cli:"not-found"`
	}
	testStructOptionInt64NotFound struct {
		Int64 int64 `cli:"not-found"`
	}
	testStructOptionUint64NotFound struct {
		Uint64 uint64 `cli:"not-found"`
	}
	testStructOptionFloat64NotFound struct {
		Float64 float64 `cli:"not-found"`
	}
	testStructFieldTypeNotSupported struct {
		NotSupported struct{} `cli:"string-opt"`
	}
)

func TestCommand_Unmarshal(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		osArgs := []string{"main-cli", "--string-opt=STRING", "--bool-opt=true", "--int64-opt=128", "--uint64-opt=256", "--float64-opt=3.141592", "--hidden-opt=HIDDEN"}
		_, err := c.parse(context.Background(), osArgs)
		requirez.NoError(t, err)
		var v testStruct
		err = UnmarshalOptions(c, &v, WithUnmarshalOptionsOptionTagKey(DefaultTagKey))
		requirez.NoError(t, err)
	})

	t.Run("error,Ptr,ErrInvalidType", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		osArgs := []string{"main-cli", "--string-opt=STRING", "--bool-opt=true", "--int64-opt=128", "--uint64-opt=256", "--float64-opt=3.141592"}
		_, err := c.parse(context.Background(), osArgs)
		requirez.NoError(t, err)
		var v int
		err = UnmarshalOptions(c, v)
		requirez.ErrorIs(t, err, ErrInvalidType)
	})

	t.Run("error,Struct,ErrInvalidType", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		osArgs := []string{"main-cli", "--string-opt=STRING", "--bool-opt=true", "--int64-opt=128", "--uint64-opt=256", "--float64-opt=3.141592"}
		_, err := c.parse(context.Background(), osArgs)
		requirez.NoError(t, err)
		var v int
		err = UnmarshalOptions(c, &v)
		requirez.ErrorIs(t, err, ErrInvalidType)
	})

	t.Run("error,ErrStructFieldCannotBeSet", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		osArgs := []string{"main-cli", "--string-opt=STRING", "--bool-opt=true", "--int64-opt=128", "--uint64-opt=256", "--float64-opt=3.141592"}
		_, err := c.parse(context.Background(), osArgs)
		requirez.NoError(t, err)
		var v testStructCannotSet
		err = UnmarshalOptions(c, &v)
		requirez.ErrorIs(t, err, ErrStructFieldCannotBeSet)
	})

	t.Run("error,ErrInvalidTagValue", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		osArgs := []string{"main-cli", "--string-opt=STRING", "--bool-opt=true", "--int64-opt=128", "--uint64-opt=256", "--float64-opt=3.141592"}
		_, err := c.parse(context.Background(), osArgs)
		requirez.NoError(t, err)
		var v testStructInvalidTagValue
		err = UnmarshalOptions(c, &v)
		requirez.ErrorIs(t, err, ErrInvalidTagValue)
	})

	t.Run("error,string,ErrUnknownOption", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		osArgs := []string{"main-cli", "--bool-opt=true", "--int64-opt=128", "--uint64-opt=256", "--float64-opt=3.141592"}
		_, err := c.parse(context.Background(), osArgs)
		requirez.NoError(t, err)
		var v testStructOptionStringNotFound
		err = UnmarshalOptions(c, &v)
		requirez.ErrorIs(t, err, ErrUnknownOption)
	})

	t.Run("error,bool,ErrUnknownOption", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		osArgs := []string{"main-cli", "--string-opt=STRING", "--int64-opt=128", "--uint64-opt=256", "--float64-opt=3.141592"}
		_, err := c.parse(context.Background(), osArgs)
		requirez.NoError(t, err)
		var v testStructOptionBoolNotFound
		err = UnmarshalOptions(c, &v)
		requirez.ErrorIs(t, err, ErrUnknownOption)
	})

	t.Run("error,int64,ErrUnknownOption", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		osArgs := []string{"main-cli", "--string-opt=STRING", "--bool-opt=true", "--uint64-opt=256", "--float64-opt=3.141592"}
		_, err := c.parse(context.Background(), osArgs)
		requirez.NoError(t, err)
		var v testStructOptionInt64NotFound
		err = UnmarshalOptions(c, &v)
		requirez.ErrorIs(t, err, ErrUnknownOption)
	})

	t.Run("error,uint64,ErrUnknownOption", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		osArgs := []string{"main-cli", "--string-opt=STRING", "--bool-opt=true", "--int64-opt=128", "--float64-opt=3.141592"}
		_, err := c.parse(context.Background(), osArgs)
		requirez.NoError(t, err)
		var v testStructOptionUint64NotFound
		err = UnmarshalOptions(c, &v)
		requirez.ErrorIs(t, err, ErrUnknownOption)
	})

	t.Run("error,float64,ErrUnknownOption", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		osArgs := []string{"main-cli", "--string-opt=STRING", "--bool-opt=true", "--int64-opt=128", "--uint64-opt=256"}
		_, err := c.parse(context.Background(), osArgs)
		requirez.NoError(t, err)
		var v testStructOptionFloat64NotFound
		err = UnmarshalOptions(c, &v)
		requirez.ErrorIs(t, err, ErrUnknownOption)
	})

	t.Run("error,ErrFieldTypeNotSupported", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		osArgs := []string{"main-cli", "--string-opt=STRING"}
		_, err := c.parse(context.Background(), osArgs)
		requirez.NoError(t, err)
		var v testStructFieldTypeNotSupported
		err = UnmarshalOptions(c, &v)
		requirez.ErrorIs(t, err, ErrStructFieldTypeNotSupported)
	})
}
