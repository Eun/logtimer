package mapprint

import (
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

const defaultToken = '%'

var defaultPrinter = Printer{
	KeyToken:       defaultToken,
	KeyNotFound:    KeepKey(),
	PrintValue:     defaultValuePrinterFunc,
	SuppressErrors: true,
}

// Printer is a object that can be used to initialize the printer with custom settings
// A typical example could be:
// p := Printer{
//    KeyToken:    '$',
//    KeyNotFound: DefaultValue("Unknown"),
//    DefaultBindings: map[string]interface{}{
//        "Key1": "Value1"
//    },
// }
// p.Sprintf("Key1 is $Key1")
// p.Sprintf("Key2 is $Key2")
//
type Printer struct {
	// KeyToken specifies how keys start
	// the Default value is % (percent sign)
	KeyToken        rune
	actualKeyToken  rune // the actual KeyToken to use (we don't want modify the public set printer instance fields)
	DefaultBindings interface{}
	KeyNotFound     KeyNotFoundFunc
	PrintValue      PrintValueFunc
	// SuppressErrors if possible
	SuppressErrors bool
}

// KeyNotFoundFunc describes the custom function that will be called if a Key was not found
type KeyNotFoundFunc func(w io.Writer, printer *Printer, prefix, key []rune, defaultPrinter PrintValueFunc) (int, error)

// PrintValueFunc describes the custom function that will be called to print a reflect.Value
type PrintValueFunc func(w io.Writer, printer *Printer, prefix, key []rune, value reflect.Value) (int, error)

// Some default functions for KeyNotFound

// KeepKey returns the requested key
func KeepKey() KeyNotFoundFunc {
	return func(w io.Writer, printer *Printer, prefix, key []rune, _ PrintValueFunc) (int, error) {
		n, err := writeRune(w, printer.actualKeyToken)
		written := n
		if err != nil {
			return n, err
		}
		if len(prefix) > 0 {
			n, err := writeRune(w, prefix...)
			written += n
			if err != nil {
				return n, err
			}
		}

		if len(key) > 0 {
			n, err := writeRune(w, key...)
			written += n
			if err != nil {
				return n, err
			}
		}

		return written, nil
	}
}

// DefaultValue returns a default value
func DefaultValue(defaultValue interface{}) KeyNotFoundFunc {
	return func(w io.Writer, printer *Printer, prefix, key []rune, defaultPrinter PrintValueFunc) (int, error) {
		return defaultPrinter(w, printer, prefix, key, reflect.ValueOf(defaultValue))
	}
}

// ClearKey returns an empty string
func ClearKey() KeyNotFoundFunc {
	return func(io.Writer, *Printer, []rune, []rune, PrintValueFunc) (int, error) {
		return 0, nil
	}
}

// GetKeyToken returns the rune that is actually used for the key identification
func (printer *Printer) GetKeyToken() rune {
	if printer.KeyToken == 0 {
		return defaultToken
	}
	return printer.KeyToken
}

// Fprintf formats a map/struct according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
// example:
// Fprintf(w, "%Key1", map[string]interface{}{
//    "Key1": "Value1",
// })
func (printer *Printer) Fprintf(w io.Writer, format string, bindings ...interface{}) (int, error) {
	// set key token
	// unfortunately it is not possible to use \0 as a token
	printer.actualKeyToken = printer.GetKeyToken()

	// we convert the map into our binding scheme because we need to order the keys descending (largest key first)
	// so we dont insert values to early, example:
	// Fprintf(w, "%textbye", map[string]interface{}{
	//   "text": "Hello",
	//   "textbye": "Goodbye",
	// })
	// should print "Goodbye" and not "Hellobye"

	// the provided data should overwrite the default Bindings

	binds, err := printer.makeBindings(bindings...)
	if err != nil {
		return 0, err
	}

	written := 0
	f := []rune(format)
	size := len(f)

	for i := 0; i < size; i++ {
		if f[i] != printer.actualKeyToken {
			n, err := writeRune(w, f[i])
			if err != nil {
				return written, err
			}
			written += n
			continue
		}

		n, jump, err := printer.placeValue(w, f, size, i, binds)
		if err != nil {
			return written, err
		}
		i += jump
		written += n
	}

	return written, nil
}

// Fprintf formats a map/struct according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
// example:
// Fprintf(w, "%Key1", map[string]interface{}{
//    "Key1": "Value1",
// })
func Fprintf(w io.Writer, format string, bindings ...interface{}) (int, error) {
	return defaultPrinter.Fprintf(w, format, bindings...)
}

// Sprintf formats a map/struct according according to a format specifier and returns the resulting string.
func (printer *Printer) Sprintf(format string, bindings ...interface{}) string {
	var buf strings.Builder
	_, err := printer.Fprintf(&buf, format, bindings...)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

// Sprintf formats a map/struct according according to a format specifier and returns the resulting string.
func Sprintf(format string, bindings ...interface{}) string {
	return defaultPrinter.Sprintf(format, bindings...)
}

// Printf formats a map/struct according according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func (printer *Printer) Printf(format string, bindings ...interface{}) (int, error) {
	return printer.Fprintf(os.Stdout, format, bindings...)
}

// Printf formats a map/struct according according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func Printf(format string, bindings ...interface{}) (int, error) {
	return defaultPrinter.Printf(format, bindings...)
}

type binding struct {
	Key   []rune
	Value reflect.Value
}

type bindings []binding

func runesAreEqual(a, b []rune) bool {
	s := len(a)
	if s != len(b) {
		return false
	}
	for i := s - 1; i >= 0; i-- {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func runesLower(a, b []rune) bool {
	al := len(a)
	bl := len(b)

	if al < bl {
		return true
	} else if al > bl {
		return false
	}

	for i := 0; i < al; i++ {
		if a[i] < b[i] {
			return true
		}
	}
	return false
}

func (v bindings) Get(key []rune) *binding {
	for i := len(v) - 1; i >= 0; i-- {
		if runesAreEqual(v[i].Key, key) {
			return &v[i]
		}
	}
	return nil
}

func (printer *Printer) makeBindings(values ...interface{}) (bindings, error) {
	var binds bindings
	if printer.DefaultBindings != nil {
		var err error
		binds, err = makeBindings(reflect.ValueOf(printer.DefaultBindings))
		if err != nil && !printer.SuppressErrors {
			return nil, err
		}
	}

	for i := 0; i < len(values); i++ {
		if values[i] == nil {
			continue
		}
		additionalBinds, err := makeBindings(reflect.ValueOf(values[i]))
		if err != nil && !printer.SuppressErrors {
			return nil, err
		}
		// merge bindings
		for _, b := range additionalBinds {
			// if key already exists => override
			if k := binds.Get(b.Key); k != nil {
				k.Value = b.Value
			} else {
				// if not append
				binds = append(binds, b)
			}
		}
	}

	// sort bindings
	sort.Slice(binds, func(i, j int) bool {
		return runesLower(binds[i].Key, binds[j].Key)
	})

	return binds, nil
}

func makeBindings(v reflect.Value) (bindings, error) {
	for v.IsValid() && v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() {
		return nil, internalError{"value must be map or struct"}
	}
	switch v.Kind() {
	case reflect.Map:
		mapKeys := v.MapKeys()
		binds := make([]binding, len(mapKeys))
		for i, key := range mapKeys {
			keyValue, err := defaultReflectPrinter.sprint(key)
			if err != nil {
				return nil, err
			}
			binds[i].Key = []rune(keyValue)
			binds[i].Value = v.MapIndex(key)
		}
		return binds, nil
	case reflect.Struct:
		t := v.Type()
		binds := make([]binding, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			binds[i].Key = []rune(t.Field(i).Name)
			binds[i].Value = v.Field(i)
		}
		return binds, nil
	default:
		return nil, internalError{"value must be map or struct, was " + v.Type().String()}
	}
}

func (printer *Printer) placeValue(w io.Writer, f []rune, size int, index int, bindings bindings) (int, int, error) {
	index++
	lastRune := utf8.RuneError
	keyPos := -1
	for i := index; i < size; i++ {
		if unicode.IsSpace(f[i]) {
			size = i
			break
		}

		if unicode.IsLetter(f[i]) && (unicode.IsNumber(lastRune) || lastRune == '.') {
			keyPos = i
			break
		}
		lastRune = f[i]
	}

	if keyPos == -1 {
		keyPos = index
	}

	var prefix []rune
	if keyPos > index {
		prefix = f[index:keyPos]
	}

	// find the end of the key
	keyEnd := size
	for i := keyPos; i < size; i++ {
		if !unicode.IsLetter(f[i]) {
			if i > keyPos && unicode.IsNumber(f[i]) {
				continue
			}
			keyEnd = i
			break
		}
	}

	// fmt.Printf("PREFIX: `%s', KEY: `%s'\n", string(prefix), string(f[keyPos:keyEnd]))
	if keyPos == keyEnd {
		n, err := writeRune(w, printer.actualKeyToken)
		if f[keyEnd] == printer.actualKeyToken { // print %% to %
			return n, 1, err
		}
		return n, 0, err
	}

	var printValue PrintValueFunc
	if printer.PrintValue == nil {
		printValue = defaultValuePrinterFunc
	} else {
		printValue = printer.PrintValue
	}

	for i := keyEnd; i >= keyPos+1; i-- {
		binding := bindings.Get(f[keyPos:i])
		if binding != nil {
			n, err := printValue(w, printer, prefix, binding.Key, binding.Value)
			if err != nil {
				if printer.SuppressErrors && isInternalError(err) {
					return 0, i - index, nil
				}
				return 0, 0, err
			}
			return n, i - index, nil
		}
	}

	var fallBack KeyNotFoundFunc
	if printer.KeyNotFound == nil {
		fallBack = KeepKey()
	} else {
		fallBack = printer.KeyNotFound
	}
	n, err := fallBack(w, printer, prefix, f[keyPos:keyEnd], printValue)
	if err != nil {
		if printer.SuppressErrors && isInternalError(err) {
			return 0, keyEnd - index, nil
		}
		return 0, 0, err
	}
	return n, keyEnd - index, nil
}

func writeRune(w io.Writer, r ...rune) (int, error) {
	return io.WriteString(w, string(r))
}

type internalError struct {
	Message string
}

func (e internalError) Error() string {
	return e.Message
}

func isInternalError(err error) bool {
	_, ok := err.(internalError)
	return ok
}
