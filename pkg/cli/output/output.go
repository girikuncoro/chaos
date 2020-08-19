package output

import (
	"io"

	"github.com/gosuri/uitable"
	"github.com/pkg/errors"
)

// Format is output format type
type Format string

const (
	Table Format = "table"
)

// String returns the string representation of format
func (o Format) String() string {
	return string(o)
}

// Write output in the given format to io.Writer.
func (o Format) Write(out io.Writer, w Writer) error {
	switch o {
	case Table:
		return w.WriteTable(out)
	}
	return errors.New("invalid format type")
}

// Writer is an interface that any type can implement to write supported formats
type Writer interface {
	// WriteTable will write tabular output into io.Writer
	WriteTable(out io.Writer) error
}

// EncodeTable is helper function for DRY-ing the printer functions.
func EncodeTable(out io.Writer, table *uitable.Table) error {
	raw := table.Bytes()
	raw = append(raw, []byte("\n")...)
	_, err := out.Write(raw)
	if err != nil {
		return errors.Wrap(err, "unable to write table output")
	}
	return nil
}

// ParseFormat takes raw string and returns matching Format.
func ParseFormat(s string) (out Format, err error) {
	switch s {
	case Table.String():
		out, err = Table, nil
	default:
		out, err = "", errors.New("invalid format type")
	}
	return
}
