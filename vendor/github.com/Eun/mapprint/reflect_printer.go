package mapprint

import (
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func defaultValuePrinterFunc(w io.Writer, printer *Printer, prefix, key []rune, value reflect.Value) (int, error) {
	// Prefix for all types:
	// +AB10
	//    ^^ Total value length is 10
	//  ^^ Fill padding with AB
	// ^ + => pad left  (ABABAHello) (default) (optional)
	//   - => pad right   (HelloABABA)
	//   | => pad middle (ABHelloABA)
	//
	//
	// Prefix only for float types:
	// +AB10.2
	//      ^^ precision is 2
	//
	// Prefix only for slice types:
	// +AB10.2
	//      ^^ take item at index 2

	if len(prefix) == 0 {
		return defaultReflectPrinter.Fprint(w, value)
	}

	type PadDirectrion int
	const (
		PadDirectionLeft PadDirectrion = iota
		PadDirectionRight
		PadDirectionCenter
	)
	padDirection := PadDirectionLeft

	rp := &reflectPrinter{
		EvaluateFunctions: true,
		IntegerBase:       10,
	}

	var i int
	var padCount int
	var padRunes []rune

	switch prefix[0] {
	case '+':
		i = 1
	case '-':
		i = 1
		padDirection = PadDirectionRight
	case '|':
		i = 1
		padDirection = PadDirectionCenter
	}

	for ; i < len(prefix); i++ {
		// we need to decide when the padding is over
		// %02d
		//  ^ padding (delimiter = 2)
		// %A2d
		//  ^ padding (delimiter = 2)
		// %A.2d
		//  ^ padding (delimiter = .)
		// so the delimiter is always a number or . but not 0 (because 0 can be used as padding)

		if prefix[i] != '0' && (prefix[i] == '.' || unicode.IsNumber(prefix[i])) {
			var padding string
			// . is for float and slices
			// before . is padding, after . is precision
			parts := strings.SplitN(string(prefix[i:]), ".", 2)
			if len(parts) == 2 {
				padding = parts[0]
				// parse precision
				// if %2.d is used avoid parsing the precision
				if len(parts[1]) > 0 {
					n, err := strconv.ParseInt(parts[1], 10, 64)
					if err != nil {
						// unskipable error
						return 0, fmt.Errorf("format %c%s%s is invalid", printer.GetKeyToken(), string(prefix), string(key))
					}
					fp := int(n)
					rp.FloatPrecision = &fp
					rp.TakeSliceItem = &fp
				} else { // there is only a . specified but no value, so assume its 0
					fp := 0
					rp.FloatPrecision = &fp
					rp.TakeSliceItem = &fp
				}
			} else {
				padding = string(prefix[i:])
			}
			if len(padding) > 0 {
				n, err := strconv.ParseInt(padding, 10, 64)
				if err != nil {
					// unskipable error
					return 0, fmt.Errorf("format %c%s%s is invalid", printer.GetKeyToken(), string(prefix), string(key))
				}
				padCount = int(n)
			}
			break
		} else {
			padRunes = append(padRunes, prefix[i])
		}
	}
	valueStr, err := rp.sprint(value)
	if err != nil {
		return 0, err
	}

	if padCount <= 0 {
		return io.WriteString(w, valueStr)
	}

	padRunesSize := len(padRunes)
	getPadRune := func(index int) rune {
		if padRunesSize == 0 {
			return ' '
		}
		return padRunes[index%padRunesSize]
	}

	switch padDirection {
	case PadDirectionLeft:
		valueRunes := []rune(valueStr)
		fillSize := padCount - len(valueRunes)
		if fillSize <= 0 {
			return io.WriteString(w, valueStr)
		}
		runes := make([]rune, padCount)
		for i = 0; i < fillSize; i++ {
			runes[i] = getPadRune(i)
		}
		copy(runes[fillSize:], valueRunes)
		return writeRune(w, runes...)
	case PadDirectionCenter:
		valueRunes := []rune(valueStr)
		size := len(valueRunes)
		fillSize := padCount - size
		if fillSize <= 0 {
			return io.WriteString(w, valueStr)
		}
		left := int(math.Floor(float64(fillSize) / 2))

		runes := make([]rune, padCount)
		for i = 0; i < left; i++ {
			runes[i] = getPadRune(i)
		}

		copy(runes[left:], valueRunes)

		for i = size + left; i < padCount; i++ {
			runes[i] = getPadRune(i - size - left)
		}

		return writeRune(w, runes...)
	default: // (PadDirectionRight)
		runes := []rune(valueStr)
		size := len(runes)
		for i = 0; size < padCount; i++ {
			runes = append(runes, getPadRune(i))
			size++
		}
		return writeRune(w, runes...)
	}
}

type reflectPrinter struct {
	FloatPrecision    *int
	TakeSliceItem     *int
	IntegerBase       int
	EvaluateFunctions bool
}

var defaultReflectPrinter = reflectPrinter{
	IntegerBase:       10,
	EvaluateFunctions: true,
}

func (p *reflectPrinter) Sprint(v reflect.Value) string {
	s, err := p.sprint(v)
	if err != nil {
		panic(err)
	}
	return s
}

func (p *reflectPrinter) sprint(v reflect.Value) (string, error) {
	var sb strings.Builder
	_, err := p.Fprint(&sb, v)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

func (p *reflectPrinter) Fprint(w io.Writer, v reflect.Value) (int, error) {
	for v.IsValid() && v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() {
		return 0, internalError{"value is invalid"}
	}
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return io.WriteString(w, "true")
		}
		return io.WriteString(w, "false")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if p.FloatPrecision != nil {
			return p.Fprint(w, reflect.ValueOf(float64(v.Int())))
		}
		return io.WriteString(w, strconv.FormatInt(v.Int(), p.IntegerBase))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if p.FloatPrecision != nil {
			return p.Fprint(w, reflect.ValueOf(float64(v.Uint())))
		}
		return io.WriteString(w, strconv.FormatUint(v.Uint(), p.IntegerBase))
	case reflect.Float32:
		fp := 6
		if p.FloatPrecision != nil {
			fp = *p.FloatPrecision
		}
		return io.WriteString(w, strconv.FormatFloat(v.Float(), 'f', fp, 32))
	case reflect.Float64:
		fp := 6
		if p.FloatPrecision != nil {
			fp = *p.FloatPrecision
		}
		return io.WriteString(w, strconv.FormatFloat(v.Float(), 'f', fp, 64))
	case reflect.String:
		return io.WriteString(w, v.String())
	case reflect.Interface:
		value := v.Elem()
		if !value.IsValid() {
			return 0, internalError{"interface is invalid"}
		}
		return p.Fprint(w, value)
	case reflect.Array, reflect.Slice:
		if p.TakeSliceItem != nil {
			if *p.TakeSliceItem >= v.Len() || *p.TakeSliceItem < 0 {
				return 0, internalError{"slice out of bounds"}
			}
			return p.Fprint(w, v.Index(*p.TakeSliceItem))
		}

		var written int
		n, err := writeRune(w, '[')
		written += n
		if err != nil {
			return written, err
		}

		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				n, err = io.WriteString(w, ", ")
				written += n
				if err != nil {
					return written, err
				}
			}
			n, err = p.Fprint(w, v.Index(i))
			written += n
			if err != nil {
				return written, err
			}
		}

		n, err = writeRune(w, ']')
		written += n
		return written, err
	case reflect.Func:
		if p.EvaluateFunctions {
			returnValues := v.Call(nil)
			l := len(returnValues)
			switch l {
			case 0:
				return 0, internalError{"function does not return value"}
			case 1:
				return p.Fprint(w, returnValues[0])
			default:
				// create an slice that has the return values
				values := reflect.MakeSlice(reflect.TypeOf([]interface{}{}), l, l)
				for i := 0; i < l; i++ {
					values.Index(i).Set(returnValues[i])
				}
				return p.Fprint(w, values)
			}
		}
		return 0, nil
	default:
		return 0, internalError{v.Type().String() + " is not supported"}
	}
}
