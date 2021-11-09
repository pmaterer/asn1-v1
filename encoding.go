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
	w            io.Writer
	buf          *bytes.Buffer
	encodingFunc func(reflect.Value) ([]byte, error)
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
		e.encodingFunc = encodeObjectIdentifier
		tag = TagObjectIdentifier
	case timeType:
		if options.timeType == TagUTCTime {
			e.encodingFunc = encodeUTCTime
			tag = TagUTCTime
		} else {
			e.encodingFunc = encodeGeneralizedTime
			tag = TagGeneralizedTime
		}
	}

	if e.encodingFunc == nil {
		switch v.Kind() {
		case reflect.Bool:
			e.encodingFunc = encodeBool
			tag = TagBoolean
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			e.encodingFunc = encodeInt
			tag = TagInteger
		case reflect.Float32, reflect.Float64:
			e.encodingFunc = encodeReal
			tag = TagReal
		case reflect.String:
			e.encodingFunc = encodeString
			if options.stringType != 0 {
				tag = options.stringType
			} else {
				tag = TagPrintableString
			}
		case reflect.Struct:
			e.encodingFunc = encodeStruct
			primitive = false
			tag = TagSet
		case reflect.Array, reflect.Slice:
			if v.Type().Elem().Kind() == reflect.Uint8 {
				e.encodingFunc = encodeOctetString
				tag = TagOctetString
			} else {
				e.encodingFunc = encodeSequence
				tag = TagSequence
				primitive = false
			}
		default:
			return fmt.Errorf("unsupported go type '%s'", v.Type())
		}
	}

	b, err := e.encodingFunc(v)
	if err != nil {
		return err
	}
	_, err = e.buf.Write(b)
	if err != nil {
		return err
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

	_, err = e.w.Write(e.buf.Bytes())
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
		b = append(b, encodeBase128(uint64(tag))...)
	}

	e.buf.Write(b)
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

func encodeSequence(v reflect.Value) ([]byte, error) {
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
	default:
		return nil, invalidTypeError("array/slice", v)
	}

	buf := new(bytes.Buffer)
	for i := 0; i < v.Len(); i++ {
		b, err := Marshal(v.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}
	return buf.Bytes(), nil
}

func invalidTypeError(expected string, value reflect.Value) error {
	return fmt.Errorf("invalid go type '%s', expecting '%s'", value.Type(), expected)
}
