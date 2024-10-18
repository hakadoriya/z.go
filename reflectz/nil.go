package reflectz

import "reflect"

func IsNil(v interface{}) bool {
	return (v == nil) || reflect.ValueOf(v).IsNil()
}
