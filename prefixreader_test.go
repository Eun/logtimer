package logtimer

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPrefixReader(t *testing.T) {
	t.Run("Normal Usage", func(t *testing.T) {
		var in TestBuffer
		var index int
		reader := &PrefixReader{
			Reader: &in,
			Format: func() string {
				defer func() {
					index++
				}()
				return fmt.Sprintf("%d ", index)
			},
		}

		copyChan := make(chan error)

		var out bytes.Buffer
		go func() {
			_, err := io.Copy(&out, reader)
			copyChan <- err
		}()

		in.WriteString("Hello World\n")
		time.Sleep(time.Millisecond * 100)
		in.WriteString("Hello ")
		time.Sleep(time.Millisecond * 100)
		in.WriteString("Universe\n")
		time.Sleep(time.Millisecond * 100)
		in.WriteString("Hello\nRest")
		time.Sleep(time.Millisecond * 100)
		in.WriteString("\nTe\nst")
		time.Sleep(time.Millisecond * 100)
		in.WriteString("\n\n\n")
		time.Sleep(time.Millisecond * 100)
		in.WriteString("Foo")
		time.Sleep(time.Millisecond * 100)
		in.WriteString("Bar")
		in.Close()

		<-copyChan

		require.Equal(t, "0 Hello World\n1 Hello Universe\n2 Hello\n3 Rest\n4 Te\n5 st\n6 \n7 \n8 FooBar", out.String())
	})

	t.Run("Small buffer", func(t *testing.T) {
		var in TestBuffer
		var index int
		reader := &PrefixReader{
			Reader: &in,
			Format: func() string {
				defer func() {
					index++
				}()
				return fmt.Sprintf("%d ", index)
			},
		}

		copyChan := make(chan error)

		var out bytes.Buffer
		go func() {
			var p [6]byte
			for {
				n, err := reader.Read(p[:])
				if err != nil {
					copyChan <- err
					return
				}
				out.Write(p[:n])
			}
		}()

		in.WriteString("Hello World\n")
		time.Sleep(time.Millisecond * 100)
		in.WriteString("Hello Universe\n")
		in.Close()

		<-copyChan

		require.Equal(t, "0 Hello World\n1 Hello Universe\n", out.String())
	})
}

type TestBuffer struct {
	bytes.Buffer
	closed bool
}

func (b *TestBuffer) Read(p []byte) (int, error) {
	if b.Buffer.Len() == 0 {
		if b.closed {
			return 0, io.EOF
		}
		for b.Buffer.Len() == 0 {
			time.Sleep(time.Millisecond * 100)
		}
	}
	return b.Buffer.Read(p)
}
func (b *TestBuffer) Close() {
	b.closed = true
}
