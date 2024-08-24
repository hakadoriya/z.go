package importz

import (
	"errors"
	"go/format"
	"io"
	"strings"
	"testing"

	"golang.org/x/tools/imports"
)

func TestNewImportsWriter(t *testing.T) {
	t.Parallel()

	t.Run("success,format", func(t *testing.T) {
		t.Parallel()

		formatted, err := Format([]byte(`
package main

import (
	"context"
	"fmt"
)

func main(){
fmt.Println("Hello, World!")
}
`))

		if err != nil {
			t.Errorf("❌: err != nil: %+v", err)
		}

		const expected = `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
}
`

		actual := string(formatted)
		if expected != actual {
			t.Errorf("❌: expected(%q) != actual(%q)", expected, actual)
		}
	})

	t.Run("error,format", func(t *testing.T) {
		t.Parallel()

		_, err := Format([]byte(`
package main

import (
	"context"
	"fmt"
)

func main(){
fmt.Println("Hello, World!")
`))

		const expected = `expected '}', found 'EOF'`
		if err == nil || !strings.Contains(err.Error(), expected) {
			t.Errorf("❌: expected(%q) != err(%q)", expected, err)
		}
	})

	t.Run("error,imports", func(t *testing.T) {
		t.Parallel()

		_, err := internalFormat(&pkg{
			FormatSourceFunc: format.Source,
			ImportsProcessFunc: func(filename string, src []byte, opt *imports.Options) ([]byte, error) {
				return nil, io.ErrUnexpectedEOF
			},
		}, []byte(`
		package main
		
		import (
			"context"
			"fmt"
		)
		
		func main(){
		fmt.Println("Hello, World!")
		}
		`))

		const expected = `expected '}', found 'EOF'`
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Errorf("❌: expected(%q) != err(%q)", io.ErrUnexpectedEOF, err)
		}
	})
}
