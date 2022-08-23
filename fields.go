package transport

import (
	"bufio"
	"bytes"
	"encoding"
	"fmt"
	"io"
	"net/textproto"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type FieldType int

const (
	UnknownFieldType FieldType = iota
	UnsignedFieldType
	IntegerFieldType
	BooleanFieldType
	StringFieldType
	FloatFieldType
	TimeFieldType
	URIFieldType
)

var (
	urlType  = reflect.TypeOf((*url.URL)(nil))
	timeType = reflect.TypeOf(time.Time{})
)

// Fields represents the key-value pairs in a Message.
//
// The keys should use a canonical form, as returned by CanonicalFieldKey.
//
// This type closely matches the textproto.MIMEfields found in "net/textproto"
// and fields found in "net/http".
type Fields map[string][]string

type FieldMarshaler interface {
	//MarshalFields encodes the receiver into a Fields instance and returns the
	//result.
	MarshalFields() (Fields, error)
}

// FieldUnmarshaler is the interface implemented by any object that can
// unmarshal a Fields instance into itself.
type FieldUnmarshaler interface {
	// UnmarshalFields decodes the Fields provided into the receiver.
	UnmarshalFields(Fields) error
}

// Add adds the key, value pair to the fields.
//
// It appends to any existing values associated with key. The key is
// case-insensitive; it is canonicalized by CanonicalFieldsKey.
func (fields Fields) Add(key, value string) {
	textproto.MIMEHeader(fields).Add(key, value)
}

// Del deletes  the value associated  with key. The key is case insensitive. It
// is canonicalized by CanonicalFieldsKey.
func (fields Fields) Del(key string) {
	textproto.MIMEHeader(fields).Del(key)
}

// Set sets the field entries associated with key to the single element value.
// It replaces any  existing values associated with key. The key is case
// insensitive. It is canonicalized by CanonicalFieldsKey. To use non-canonical
// keys, assign to the Fields instance directly.
func (fields Fields) Set(key, value string) {
	textproto.MIMEHeader(fields).Set(key, value)
}

// Get gets the first value associated with  the given key. If there are no
// values associated with the key, Get returns "".
//
// The key is case insensitive; it is canonicalized by CanonicalFieldsKey. To
// use non-canonical keys, use the Fields instance directly.
func (fields Fields) Get(key string) string {
	return textproto.MIMEHeader(fields).Get(key)
}

// Values returns all values associated with the given key.
//
// It is case insensitive; it is canonicalized by CanonicalFieldsKey. To use
// non-canonical keys, access the map directly.
//
// The slice returned is NOT a copy.
func (fields Fields) Values(key string) []string {
	return textproto.MIMEHeader(fields).Values(key)
}

// Write writes the Fields as though it were a MIME fields. However, it does
// not add the trailing newline.
//
// This function performs the write all at once, but does allocate internally.
func (fields Fields) Write(writer io.Writer) error {
	data, err := fields.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

// MarshalBinary turns the Fields into the correct binary representation.
//
// This is primarily called by Fields.Write.
//
// This function will error if the fields is empty or nil.
func (fields Fields) MarshalBinary() ([]byte, error) {
	var buffer bytes.Buffer
	if len(fields) == 0 {
		return nil, ErrFieldsEmpty
	}
	for key, values := range fields {
		key := CanonicalFieldsKey(key)
		value := strings.Join(values, ",")
		fmt.Fprintf(&buffer, "%s: %s\n", key, value)
	}
	return buffer.Bytes(), nil
}

// UnmarshalBinary parses the fields fields from the provided byte slice.
//
// This function does not perform validation for the contents of the fields
// fields. It also does not validate that a fields field is not empty, as this
// is technically allowed.
//
// BUG(bruxisma): This function does not currently handle multi-line fields.
func (fields Fields) UnmarshalBinary(data []byte) error {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		items := strings.SplitN(line, ":", 2)
		if len(items) != 2 {
			return fmt.Errorf("%w %q", ErrFieldEntryInvalid, line)
		}
		key := CanonicalFieldsKey(items[0])
		values := strings.Split(items[1], ",")
		for _, value := range values {
			value = strings.TrimSpace(value)
			fields.Add(key, value)
		}
	}
	return scanner.Err()
}

// CanonicalFieldsKey returns the canonical format of the field key. The
// canonicalization conversts the first letter of each word in the string to
// upper case. All other letters are lowercased. For example, the canonical key
// for "send-config" is "Send-Config". If the key contains a space or invalid
// bytes, it is returned without modifications.
func CanonicalFieldsKey(key string) string {
	return textproto.CanonicalMIMEHeaderKey(key)
}

// GetFieldName returns either the name of the given struct field or the value
// of the struct tag "transport"
func GetFieldName(field reflect.StructField) string {
	if value, ok := field.Tag.Lookup("transport"); ok {
		return value
	}
	return field.Name
}

// GetFieldType returns the FieldType for the given value.
//
// If the FieldType returned is UnknownFieldType, the value cannot be
// automatically deserialized from a field value into a Go type.
func GetFieldType(value reflect.Value) FieldType {
	switch value.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return UnsignedFieldType
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return IntegerFieldType
	case reflect.Float32, reflect.Float64:
		return FloatFieldType
	case reflect.Bool:
		return BooleanFieldType
	case reflect.String:
		return StringFieldType
	}
	if value.Type() == timeType {
		return TimeFieldType
	} else if value.Type() == urlType {
		return URIFieldType
	}
	return UnknownFieldType
}

// MarshalFields returns the Fields representation of the value provided.
//
// MarshalFields traverses the value of the provided object recursively. If an
// encountered value implements the FieldMarshaler interface (and is not nil or
// empty), MarshalFields calls its MarshalFields method to produce a
// Header.
//
// Some message representations can be represented with a simple string or map,
// and thus these two types are permitted without being passed by pointer.
// Slices are never permitted by this function, unless they implement
// FieldMarshaler.
//
// The encoding of each field in a struct can be customized by the format
// string stored under the "transport" key in the struct field's tag. The
// format string gives the name of the field. This is intended to allow for
// aliasing field names as well as field names that conflict with variable
// naming requirements in Go.
func MarshalFields(source any) (Fields, error) {
	fields := Fields{}
	if fm, ok := source.(FieldMarshaler); ok {
		return fm.MarshalFields()
	}
	value := reflect.ValueOf(source)
	if value.IsNil() {
		return nil, ErrSourceIsNil
	}
	if value.Kind() != reflect.Ptr {
		return nil, ErrSourceNotPointer
	}
	value = value.Elem()

	members := reflect.VisibleFields(value.Type())
	for _, member := range members {
		field := GetFieldName(member)
		// TODO: get the entry, as an Interface value, then run through TextMarshaler, Stringer, etc.
		ifc := value.FieldByName(member.Name).Interface()
		var content string
		switch ifc.(type) {
		case string:
			content = ifc.(string)
		case *time.Time:
			content = ifc.(*time.Time).Format(time.RFC1123)
		case time.Time:
			content = ifc.(time.Time).Format(time.RFC1123)
		case fmt.Stringer:
			content = ifc.(fmt.Stringer).String()
		case encoding.TextMarshaler:
			text, err := ifc.(encoding.TextMarshaler).MarshalText()
			if err != nil {
				return nil, err
			}
			content = string(text)
		default:
			return nil, &FieldMarshalerError{
				Type:   reflect.TypeOf(ifc),
				Err:    fmt.Errorf("cannot marshal member %q to field %q", member.Name, field),
				source: "MarshalFields",
			}
		}
		fields.Add(field, content)
	}
	return fields, nil
}

// UnmarshalFields unmarshals the provided Fields object into the destination
// provided.
//
// MarshalFields traverses the value of the provided object recursively. If an
// encountered value implements the FieldUnmarshaler interface (and is not nil
// or empty), UnmarshalFields calls its UnmarshalFields method to deserialize
// the Fields object.
//
// Some field representations can be represented with a simple string or map,
// and thus these two types are permitted without being passed by pointer.
// Slices are never permitted by this function, unless they implement the
// FieldMarshaler interface.
//
// The encoding of each field in a struct can be cutomized by the format string
// stored under the "transport" key in the struct field's tag. The format
// string gives the name of the field. This is intended to allow for aliasing
// field names as well as field names that conflict with variable naming
// requirements in Go.
func UnmarshalFields(fields Fields, destination any) error {
	// if the destination is a FieldUnmarshaler, just use that and call it a day.
	if ifc, ok := destination.(FieldUnmarshaler); ok {
		return ifc.UnmarshalFields(fields)
	}
	value := reflect.ValueOf(destination)
	if value.IsNil() {
		return ErrDestinationIsNil
	}
	if value.Kind() != reflect.Ptr {
		return ErrDestinationNotPointer
	}
	// Get the value pointed to
	value = value.Elem()

	// Get a list of all visible fields in the destination
	members := reflect.VisibleFields(value.Type())
	for _, member := range members {
		field := GetFieldName(member)
		// skip fields that are not in the fields map
		if _, ok := fields[field]; !ok {
			continue
		}
		entry := value.FieldByName(member.Name)
		if !entry.CanSet() {
			return &FieldMarshalerError{
				Type:   entry.Type(),
				Err:    fmt.Errorf("cannot set %q", member.Name),
				source: "apt/transport.UnmarshalFields",
			}
		}
		kind := GetFieldType(entry)
		if kind == StringFieldType {
			value.FieldByName(member.Name).SetString(strings.Join(fields[field], ","))
			continue
		}
		// generate a conversion table. sadly can't be done globally because we
		// need to capture the `entry` variable.
		conversionTable := map[FieldType]func(string) error{
			UnsignedFieldType: fieldConversionFunction(parseUint, entry.SetUint),
			IntegerFieldType:  fieldConversionFunction(parseInt, entry.SetInt),
			BooleanFieldType:  fieldConversionFunction(parseBool, entry.SetBool),
			FloatFieldType:    fieldConversionFunction(parseFloat, entry.SetFloat),
			TimeFieldType:     fieldConversionFunction(parseTime, assignTime(entry)),
			URIFieldType:      fieldConversionFunction(parseURI, assignURL(entry)),
		}

		// one thing to also consider is checking if TextUnmarshaler is implemented
		// by the field and falling back to that if there is no known conversion,
		// and using the same value.FieldByName trick from above for the
		// StringFieldType
		converter, ok := conversionTable[kind]
		if !ok {
			return &FieldMarshalerError{
				Type:   entry.Type(),
				Err:    fmt.Errorf("cannot assign to field %q, no known conversion", member.Name),
				source: "apt/transport.UnmarshalFields",
			}
		}
		if err := converter(fields[field][0]); err != nil {
			return err
		}
	}
	return nil
}

func fieldConversionFunction[T any](parser func(string) (T, error), assigner func(T)) func(string) error {
	return func(text string) error {
		parsed, err := parser(text)
		if err != nil {
			return err
		}
		assigner(parsed)
		return nil
	}
}

/* parse functions provided for consistency */

func parseFloat(text string) (float64, error) {
	return strconv.ParseFloat(text, 64)
}

func parseUint(text string) (uint64, error) {
	return strconv.ParseUint(text, 10, 64)
}

func parseInt(text string) (int64, error) {
	return strconv.ParseInt(text, 10, 64)
}

func parseBool(text string) (bool, error) {
	return strconv.ParseBool(text)
}

func parseTime(text string) (time.Time, error) {
	return time.Parse(time.RFC1123, text)
}

func parseURI(text string) (*url.URL, error) {
	return url.Parse(text)
}

func assignTime(value reflect.Value) func(time.Time) {
	return func(time time.Time) {
		value.Set(reflect.ValueOf(time))
	}
}

func assignURL(value reflect.Value) func(*url.URL) {
	return func(url *url.URL) {
		value.Set(reflect.ValueOf(url))
	}
}
