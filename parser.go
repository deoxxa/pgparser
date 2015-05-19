package pgparser

import (
	"fmt"
	"reflect"
	"strconv"

	"go.bmatsuo.co/go-lexer"
)

// Unmarshal takes a string representation of a complex type from Postgres and
// unpacks it into value v. Since the serialised data doesn't include field
// names or even types, values must match perfectly to be decoded. This means
// that structs must have their fields in the same order as they are in the
// database, and those fields must be the same types.
func Unmarshal(s string, v interface{}) error {
	return (&unmarshaller{
		l: lexer.New(stateBegin, s),
	}).unmarshal(reflect.ValueOf(v))
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

func (u *unmarshaller) unmarshal(v reflect.Value) error {
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("can't unmarshal into non-pointer type %s", v.Type().String())
	}

	switch v.Elem().Kind() {
	case reflect.Slice:
		return u.unmarshalSlice(v.Elem())
	case reflect.Struct:
		return u.unmarshalStruct(v.Elem())
	case reflect.String:
		return u.unmarshalString(v.Elem())
	case reflect.Int:
		return u.unmarshalInt(v.Elem())
	default:
		return fmt.Errorf("can't unmarshal into type %s", v.Type().String())
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
		if err := u.unmarshal(e); err != nil {
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
		}).unmarshal(v.Addr())
	}

	if i.Type != itemLeftParen {
		return fmt.Errorf("invalid token; expected left paren but got %q", i.String())
	}

	for j := 0; j < v.NumField(); j++ {
		f := v.Field(j)

		if err := u.unmarshal(f.Addr()); err != nil {
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

func (u *unmarshaller) unmarshalString(v reflect.Value) error {
	i := u.next()

	switch i.Type {
	case itemQuotedString:
		s, err := strconv.Unquote(i.Value)
		if err != nil {
			return err
		}

		v.SetString(s)
	case itemBareString:
		v.SetString(i.Value)
	default:
		return fmt.Errorf("invalid token; expected quoted or bare string but got %q", i.String())
	}

	return nil
}

func (u *unmarshaller) unmarshalInt(v reflect.Value) error {
	i := u.next()

	if i.Type != itemBareString {
		return fmt.Errorf("invalid token; expected bare string but got %q", i.String())
	}

	n, err := strconv.ParseInt(i.Value, 10, 64)
	if err != nil {
		return err
	}

	v.SetInt(n)

	return nil
}
