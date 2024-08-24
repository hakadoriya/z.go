package buildinfoz

import (
	"io"

	"github.com/hakadoriya/z.go/bytez"
)

func Fprint(w io.Writer) error {
	var p []byte
	p = append(p, `{"version":"`...)
	p = bytez.AppendJSONEscapedString(p, buildVersion)
	p = append(p, `","revision":"`...)
	p = bytez.AppendJSONEscapedString(p, buildRevision)
	p = append(p, `","branch":"`...)
	p = bytez.AppendJSONEscapedString(p, buildBranch)
	p = append(p, `","timestamp":"`...)
	p = bytez.AppendJSONEscapedString(p, buildTimestamp)
	p = append(p, `","goVersion":"`...)
	p = bytez.AppendJSONEscapedString(p, debugBuildInfo.GoVersion)
	if cgoEnabled := CGOEnabled(); cgoEnabled != "" {
		p = append(p, `","cgoEnabled":"`...)
		p = bytez.AppendJSONEscapedString(p, buildCGOEnabled)
	}
	p = append(p, `"}`+"\n"...)

	//nolint:wrapcheck
	if _, err := w.Write(p); err != nil {
		return err
	}

	return nil
}
