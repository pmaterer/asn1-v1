package asn1

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

const (
	tagKey = "asn1"
)

func Marshal(v interface{}) ([]byte, error) {
	return MarshalWithOptions(v, "")
}

func MarshalWithOptions(v interface{}, opts string) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := NewEncoder(buf).Encode(v, opts)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type Encoder struct {
	w   io.Writer
	buf *bytes.Buffer
}

func NewEncoder(w io.Writer) *Encoder {
	buf := new(bytes.Buffer)
	return &Encoder{
		w:   w,
		buf: buf,
	}
}

func (e *Encoder) Encode(v interface{}, opts string) error {
	return e.encode(reflect.ValueOf(v), opts)
}

func (e *Encoder) encode(v reflect.Value, opts string) error {
	var tag Tag
	primitive := true
	class := TagClassUniversal

	options := parseOptions(opts)

	// check special types first
	switch v.Type() {
	case oidType:
		b, err := encodeObjectIdentifier(v)
		if err != nil {
			return err
		}
		e.buf.Write(b)
		tag = TagObjectIdentifier
	case timeType:
		if options.timeType == TagUTCTime {
			b, err := encodeUTCTime(v)
			if err != nil {
				return err
			}
			e.buf.Write(b)
			tag = TagUTCTime
		} else {
			b, err := encodeGeneralizedTime(v)
			if err != nil {
				return err
			}
			e.buf.Write(b)
			tag = TagGeneralizedTime
		}
	}

	if e.buf.Len() == 0 {
		switch v.Kind() {
		case reflect.Bool:
			b, err := encodeBool(v)
			if err != nil {
				return err
			}
			e.buf.Write(b)
			tag = TagBoolean
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			b, err := encodeInt(v)
			if err != nil {
				return err
			}
			e.buf.Write(b)
			tag = TagInteger
		case reflect.Float32, reflect.Float64:
			b, err := encodeReal(v)
			if err != nil {
				return err
			}
			e.buf.Write(b)
			tag = TagReal
		case reflect.String:
			b, err := encodeString(v)
			if err != nil {
				return err
			}
			e.buf.Write(b)
			if options.stringType != 0 {
				tag = options.stringType
			} else {
				tag = TagPrintableString
			}
		case reflect.Struct:
			b, err := encodeStruct(v)
			if err != nil {
				return err
			}
			e.buf.Write(b)
			primitive = false
			tag = TagSequence
		case reflect.Array, reflect.Slice:
			if v.Type().Elem().Kind() == reflect.Uint8 {
				b, err := encodeOctetString(v)
				if err != nil {
					return err
				}
				e.buf.Write(b)
				tag = TagOctetString
			}
		default:
			return fmt.Errorf("unsupported go type '%s'", v.Type())
		}
	}

	if options.private {
		class = TagClassPrivate
	}

	body := make([]byte, e.buf.Len())
	copy(body, e.buf.Bytes())
	e.buf.Reset()

	e.encodeIdentifier(class, tag, primitive)
	e.encodeLength(body)
	e.buf.Write(body)

	_, err := e.w.Write(e.buf.Bytes())
	if err != nil {
		return err
	}

	return nil

}

func (e *Encoder) encodeIdentifier(class TagClass, tag Tag, primitive bool) {
	b := []byte{0x00}

	b[0] |= byte(class << 6)

	if primitive {
		b[0] |= byte(0 << 5)
	} else {
		b[0] |= byte(1 << 5)
	}

	// universal tags 0-30
	if tag <= 30 {
		b[0] |= byte(tag)
	} else {
		b[0] |= byte(0x1f)
		b = append(b, encodeHighTag(tag)...)
	}

	e.buf.Write(b)
}

func encodeHighTag(tag Tag) []byte {
	buf := new(bytes.Buffer)

	for tag != 0 {
		t := tag & 0x7f
		tag >>= 7
		if len(buf.Bytes()) != 0 {
			t |= 0x80
		}
		buf.WriteByte(byte(t))
	}

	return reverseBytes(buf.Bytes())
}

func (e *Encoder) encodeLength(body []byte) {
	// only definite form supported
	// length encoded as unsigned binary integers

	length := len(body)
	b := new(bytes.Buffer)

	lengthBytes := encodeUint(uint64(length))

	// short form
	if length <= 0x7f {
		e.buf.Write(lengthBytes)
		return
	}

	// long form
	header := len(lengthBytes) | 0x80

	b.Write(encodeUint(uint64(header)))
	b.Write(lengthBytes)
	e.buf.Write(b.Bytes())
}

func invalidTypeError(expected string, value reflect.Value) error {
	return fmt.Errorf("invalid go type '%s', expecting '%s'", value.Type(), expected)
}
