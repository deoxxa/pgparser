package pgparser

import (
	"fmt"
	"reflect"
	"strconv"

	"go.bmatsuo.co/go-lexer"
)

type postgresTextUnmarshaller interface {
	UnmarshalPostgresText(b []byte) error
}

// Unmarshal takes a string representation of a complex type from Postgres and
// unpacks it into value v. Since the serialised data doesn't include field
// names or even types, values must match perfectly to be decoded. This means
// that structs must have their fields in the same order as they are in the
// database, and those fields must be the same types.
func Unmarshal(s string, v interface{}) error {
	return (&unmarshaller{
		l: lexer.New(stateBegin, s),
	}).unmarshal(v)
}

type unmarshaller struct {
	l     *lexer.Lexer
	items []*lexer.Item
}

func (u *unmarshaller) next() *lexer.Item {
	if l := len(u.items); l > 0 {
		i := u.items[l-1]
		u.items[l-1] = nil
		u.items = u.items[:l-1]
		return i
	}

	return u.l.Next()
}

func (u *unmarshaller) putback(i *lexer.Item) {
	u.items = append(u.items, i)
}

func (u *unmarshaller) unmarshal(v interface{}) error {
	switch p := v.(type) {
	case *int:
		return u.unmarshalInt(p)
	case *uint:
		return u.unmarshalUint(p)
	case *int8:
		return u.unmarshalInt8(p)
	case *uint8:
		return u.unmarshalUint8(p)
	case *int16:
		return u.unmarshalInt16(p)
	case *uint16:
		return u.unmarshalUint16(p)
	case *int32:
		return u.unmarshalInt32(p)
	case *uint32:
		return u.unmarshalUint32(p)
	case *int64:
		return u.unmarshalInt64(p)
	case *uint64:
		return u.unmarshalUint64(p)
	case *string:
		return u.unmarshalString(p)
	case *[]byte:
		return u.unmarshalByteSlice(p)
	}

	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("can't unmarshal into non-pointer type %s", rv.Type().String())
	}

	switch rv.Elem().Kind() {
	case reflect.Slice:
		return u.unmarshalSlice(rv.Elem())
	case reflect.Struct:
		return u.unmarshalStruct(rv.Elem())
	default:
		return fmt.Errorf("can't unmarshal into type %s", rv.Type().String())
	}
}

func (u *unmarshaller) unmarshalSlice(v reflect.Value) error {
	if i := u.next(); i.Type != itemLeftBrace {
		return fmt.Errorf("invalid token; expected left brace but got %q", i.String())
	}

	i := u.next()
	if i.Type == itemRightBrace {
		return nil
	}
	u.putback(i)

	for {
		e := reflect.New(v.Type().Elem())
		if err := u.unmarshal(e.Interface()); err != nil {
			return err
		}
		v.Set(reflect.Append(v, e.Elem()))

		if i := u.next(); i.Type != itemComma && i.Type != itemRightBrace {
			return fmt.Errorf("invalid token; expected right brace or comma but got %q", i.String())
		} else if i.Type == itemRightBrace {
			break
		}
	}

	return nil
}

func (u *unmarshaller) unmarshalStruct(v reflect.Value) error {
	i := u.next()

	if i.Type == itemQuotedString {
		s, err := strconv.Unquote(i.Value)
		if err != nil {
			return err
		}

		return (&unmarshaller{
			l: lexer.New(stateBegin, s),
		}).unmarshal(v.Addr().Interface())
	}

	if i.Type != itemLeftParen {
		return fmt.Errorf("invalid token; expected left paren but got %q", i.String())
	}

	for j := 0; j < v.NumField(); j++ {
		f := v.Field(j)

		if err := u.unmarshal(f.Addr().Interface()); err != nil {
			return err
		}

		if j == v.NumField()-1 {
			break
		}

		if i := u.next(); i.Type != itemComma {
			return fmt.Errorf("invalid token; expected comma but got %q", i.String())
		}
	}

	if i := u.next(); i.Type != itemRightParen {
		return fmt.Errorf("invalid token; expected right paren but got %q", i.String())
	}

	return nil
}

func (u *unmarshaller) unmarshalInt(v *int) error {
	i := u.next()

	if i.Type != itemBareString {
		return fmt.Errorf("invalid token; expected bare string but got %q", i.String())
	}

	n, err := strconv.ParseInt(i.Value, 10, 32)
	if err != nil {
		return err
	}

	*v = int(n)

	return nil
}

func (u *unmarshaller) unmarshalUint(v *uint) error {
	i := u.next()

	if i.Type != itemBareString {
		return fmt.Errorf("invalid token; expected bare string but got %q", i.String())
	}

	n, err := strconv.ParseUint(i.Value, 10, 32)
	if err != nil {
		return err
	}

	*v = uint(n)

	return nil
}

func (u *unmarshaller) unmarshalInt8(v *int8) error {
	i := u.next()

	if i.Type != itemBareString {
		return fmt.Errorf("invalid token; expected bare string but got %q", i.String())
	}

	n, err := strconv.ParseInt(i.Value, 10, 8)
	if err != nil {
		return err
	}

	*v = int8(n)

	return nil
}

func (u *unmarshaller) unmarshalUint8(v *uint8) error {
	i := u.next()

	if i.Type != itemBareString {
		return fmt.Errorf("invalid token; expected bare string but got %q", i.String())
	}

	n, err := strconv.ParseUint(i.Value, 10, 8)
	if err != nil {
		return err
	}

	*v = uint8(n)

	return nil
}

func (u *unmarshaller) unmarshalInt16(v *int16) error {
	i := u.next()

	if i.Type != itemBareString {
		return fmt.Errorf("invalid token; expected bare string but got %q", i.String())
	}

	n, err := strconv.ParseInt(i.Value, 10, 16)
	if err != nil {
		return err
	}

	*v = int16(n)

	return nil
}

func (u *unmarshaller) unmarshalUint16(v *uint16) error {
	i := u.next()

	if i.Type != itemBareString {
		return fmt.Errorf("invalid token; expected bare string but got %q", i.String())
	}

	n, err := strconv.ParseUint(i.Value, 10, 16)
	if err != nil {
		return err
	}

	*v = uint16(n)

	return nil
}

func (u *unmarshaller) unmarshalInt32(v *int32) error {
	i := u.next()

	if i.Type != itemBareString {
		return fmt.Errorf("invalid token; expected bare string but got %q", i.String())
	}

	n, err := strconv.ParseInt(i.Value, 10, 32)
	if err != nil {
		return err
	}

	*v = int32(n)

	return nil
}

func (u *unmarshaller) unmarshalUint32(v *uint32) error {
	i := u.next()

	if i.Type != itemBareString {
		return fmt.Errorf("invalid token; expected bare string but got %q", i.String())
	}

	n, err := strconv.ParseUint(i.Value, 10, 32)
	if err != nil {
		return err
	}

	*v = uint32(n)

	return nil
}

func (u *unmarshaller) unmarshalInt64(v *int64) error {
	i := u.next()

	if i.Type != itemBareString {
		return fmt.Errorf("invalid token; expected bare string but got %q", i.String())
	}

	n, err := strconv.ParseInt(i.Value, 10, 64)
	if err != nil {
		return err
	}

	*v = int64(n)

	return nil
}

func (u *unmarshaller) unmarshalUint64(v *uint64) error {
	i := u.next()

	if i.Type != itemBareString {
		return fmt.Errorf("invalid token; expected bare string but got %q", i.String())
	}

	n, err := strconv.ParseUint(i.Value, 10, 64)
	if err != nil {
		return err
	}

	*v = uint64(n)

	return nil
}

func (u *unmarshaller) unmarshalString(v *string) error {
	i := u.next()

	switch i.Type {
	case itemQuotedString:
		s, err := strconv.Unquote(i.Value)
		if err != nil {
			return err
		}

		*v = s
	case itemBareString:
		*v = i.Value
	default:
		return fmt.Errorf("invalid token; expected quoted or bare string but got %q", i.String())
	}

	return nil
}

func (u *unmarshaller) unmarshalByteSlice(v *[]byte) error {
	i := u.next()

	switch i.Type {
	case itemQuotedString:
		s, err := strconv.Unquote(i.Value)
		if err != nil {
			return err
		}

		*v = []byte(s)
	case itemBareString:
		*v = []byte(i.Value)
	default:
		return fmt.Errorf("invalid token; expected quoted or bare string but got %q", i.String())
	}

	return nil
}
