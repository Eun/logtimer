package logtimer

import (
	"bytes"
	"fmt"
	"io"

	"github.com/pborman/ansi"
)

type ColorCorrection int

const (
	Disabled ColorCorrection = iota
	Enabled
	Alternate
)

type FormatFunc func() string

type PrefixReader struct {
	Format          FormatFunc
	skipNextPrint   bool
	buffer          bytes.Buffer
	ColorCorrection ColorCorrection
	io.Reader
}

func (lt *PrefixReader) writeFormat(w io.Writer) (int, error) { //nolint: unparam // allow unused int return
	if lt.ColorCorrection == Disabled {
		return io.WriteString(w, lt.Format())
	}

	saveCursor := "\x1b[s"
	restoreCursor := "\x1b[u"

	if lt.ColorCorrection == Alternate {
		saveCursor = "\x1b7"
		restoreCursor = "\x1b8"
	}

	f := lt.Format()
	written, err := fmt.Fprintf(w, "%s\x1b[0m", saveCursor)
	if err != nil {
		return 0, err
	}
	n, err := io.WriteString(w, f)
	written += n
	if err != nil {
		return written, err
	}

	// determinate size
	r, err := ansi.Strip([]byte(f))
	if err != nil {
		return written, err
	}
	n, err = fmt.Fprintf(w, "%s\x1b[%dC", restoreCursor, len(r))
	written += n
	return written, err
}

func (lt *PrefixReader) Read(p []byte) (int, error) {
	if lt.buffer.Len() > 0 {
		return lt.buffer.Read(p)
	}
	n, err := lt.Reader.Read(p)
	if err != nil {
		return n, err
	}

	if !lt.skipNextPrint {
		_, _ = lt.writeFormat(&lt.buffer)
		lt.skipNextPrint = true
	}

	for i := 0; i < n; i++ {
		_ = lt.buffer.WriteByte(p[i])
		if p[i] == '\n' {
			if i < n-1 {
				_, _ = lt.writeFormat(&lt.buffer)
			} else {
				lt.skipNextPrint = false
			}
		}
	}

	return lt.buffer.Read(p)
}
