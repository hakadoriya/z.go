package bytez

import "testing"

func TestAppendJSONEscapedString(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		var b []byte
		b = AppendJSONEscapedString(b, `{"string":"json","ctrl":"`+"\b"+"\f"+"\n"+"\r"+"\t"+"\u0000"+`","bool":true,"number":1}`)
		const expected = `{\"string\":\"json\",\"ctrl\":\"\b\f\n\r\t\u0000\",\"bool\":true,\"number\":1}`
		actual := string(b)
		if expected != actual {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, actual)
		}
	})
}
