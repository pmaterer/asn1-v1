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
	datagram     *datagram
	encodingFunc func(reflect.Value) ([]byte, error)
	options      *options
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

func (e *Encoder) parseType(v reflect.Value) (err error) {
	var tag Tag
	var isConstructed bool

	switch v.Type() {
	case oidType:
		e.encodingFunc = encodeObjectIdentifier
		tag = TagObjectIdentifier
	case timeType:
		tag = e.options.timeType
		switch e.options.timeType {
		case TagUTCTime:
			e.encodingFunc = encodeUTCTime
		default:
			e.encodingFunc = encodeGeneralizedTime
		}
	case nullType:
		tag = TagNull
		e.encodingFunc = encodeNull
	}

	if e.encodingFunc == nil {
		switch v.Kind() {
		case reflect.Bool:
			e.encodingFunc = encodeBool
			tag = TagBoolean
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			e.encodingFunc = encodeInt
			if e.options.enumerated {
				tag = TagEnumerated
			} else {
				tag = TagInteger
			}

		case reflect.Float32, reflect.Float64:
			e.encodingFunc = encodeReal
			tag = TagReal
		case reflect.String:
			tag = e.options.stringType
			e.encodingFunc = encodeString
			switch e.options.stringType {
			case TagPrintableString:
				if !isValidPrintableString(v.String()) {
					return fmt.Errorf("string not valid printablestring")
				}
			case TagIA5String:
				if !isValidIA5String(v.String()) {
					return fmt.Errorf("string not valid ia5string")
				}
			case TagNumericString:
				if !isValidNumericString(v.String()) {
					return fmt.Errorf("string not valid numeric string")
				}
			case TagBitString:
				e.encodingFunc = EncodeBitString
			}
		case reflect.Struct:
			e.encodingFunc = encodeStruct
			isConstructed = true
			if e.options.set {
				tag = TagSet
			} else {
				tag = TagSequence
			}

		case reflect.Array, reflect.Slice:
			if v.Type().Elem().Kind() == reflect.Uint8 {
				e.encodingFunc = encodeOctetString
				tag = TagOctetString
			} else {
				e.encodingFunc = encodeList
				if e.options.set {
					tag = TagSetOf
				} else {
					tag = TagSequenceOf
				}
				isConstructed = true
			}
		default:
			return fmt.Errorf("unsupported go type '%s'", v.Type())
		}
	}

	e.datagram.identifier.tag = tag
	e.datagram.isConstructed = isConstructed

	return nil
}

func (e *Encoder) encode(v reflect.Value, opts string) error {
	var err error

	e.datagram = newDatagram()
	e.options, err = parseOptions(opts)
	if err != nil {
		return err
	}

	err = e.parseType(v)
	if err != nil {
		return err
	}

	body, err := e.encodingFunc(v)
	if err != nil {
		return err
	}
	e.datagram.body = body

	if empty(v) && e.options.optional {
		return nil
	}
	// if empty(v) {
	// 	if e.options.optional {
	// 		e.datagram.body = nil
	// 	} else {
	// 		return fmt.Errorf("empty body!")
	// 	}
	// }

	if e.options.explicit {
		if e.options.tag == nil {
			return fmt.Errorf("flag 'explicit' requires flag 'tag' to be set")
		}
		body, err := e.encodingFunc(v)
		if err != nil {
			return err
		}

		innerBody := &datagram{
			identifier: identifier{
				class:         TagClassUniversal,
				tag:           e.datagram.identifier.tag,
				isConstructed: false,
			},
			body: body,
		}
		b := innerBody.encode()
		e.datagram.body = b
		e.datagram.identifier.isConstructed = true
	}

	if e.options.tag != nil {
		if e.options.application {
			e.datagram.identifier.class = TagClassApplication
		} else if e.options.private {
			e.datagram.identifier.class = TagClassPrivate
		} else {
			e.datagram.identifier.class = TagClassContextSpecific
		}
		e.datagram.identifier.tag = Tag(*e.options.tag)
	}

	if e.options.private {
		e.datagram.identifier.class = TagClassPrivate
	}

	e.buf.Write(e.datagram.encode())

	_, err = e.w.Write(e.buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func invalidTypeError(expected string, value reflect.Value) error {
	return fmt.Errorf("invalid go type '%s', expecting '%s'", value.Type(), expected)
}
