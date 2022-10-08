package transport

import (
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type FieldTypeSuite struct {
	suite.Suite
}

type UnmarshalFieldsSuite struct {
	suite.Suite
}

type MarshalFieldsSuite struct {
	suite.Suite
}

func generateLastModifiedStruct() any {
	timestamp := time.Date(1998, 3, 31, 0, 0, 0, 0, time.UTC)
	value := struct {
		LastModified time.Time `transport:"Last-Modified"`
	}{LastModified: timestamp}
	return &value
}

func (suite *MarshalFieldsSuite) TestMarshalFields() {
	fields, err := MarshalFields(&struct {
		URI         *url.URL
		ResumePoint string `transport:"Resume-Point"`
	}{URI: &url.URL{Scheme: "http", Host: "example.com"}, ResumePoint: "test"})
	suite.Require().NoError(err)
	suite.Require().Contains(fields, CanonicalFieldsKey("URI"))
}

func (suite *MarshalFieldsSuite) TestMarshalFieldsWithTime() {
	timestamp := time.Date(1998, 3, 31, 0, 0, 0, 0, time.UTC)
	fields, err := MarshalFields(generateLastModifiedStruct())
	suite.Require().NoError(err)
	suite.Require().Contains(fields, CanonicalFieldsKey("Last-Modified"))
	suite.Require().Equal(fields["Last-Modified"][0], timestamp.Format(time.RFC1123))
}

func (suite *MarshalFieldsSuite) TestMarshalFieldsWithTimePtr() {
	timestamp := time.Date(1998, 3, 31, 0, 0, 0, 0, time.UTC)
	fields, err := MarshalFields(generateLastModifiedStruct())
	suite.Require().NoError(err)
	suite.Require().Contains(fields, CanonicalFieldsKey("Last-Modified"))
	suite.Require().Equal(fields["Last-Modified"][0], timestamp.Format(time.RFC1123))
}

func (suite *UnmarshalFieldsSuite) TestDynamic() {
	/* The first version release date for APT according to wikipedia */
	timestamp := time.Date(1998, 3, 31, 0, 0, 0, 0, time.UTC)
	fields := Fields{
		"Last-Modified": []string{timestamp.Format(time.RFC1123)},
		"URI":           []string{"test://testing.example.whatever"},
		"Password":      []string{"hunter2"},
		"Needs-Cleanup": []string{"true"},
	}
	dynamic := struct {
		LastModified time.Time `transport:"Last-Modified"`
		URI          *url.URL
		Password     string
		NeedsCleanup bool `transport:"Needs-Cleanup"`
	}{}
	err := UnmarshalFields(fields, &dynamic)
	suite.NoError(err)
	suite.Equal(timestamp, dynamic.LastModified)
	suite.Equal("hunter2", dynamic.Password)
	suite.True(dynamic.NeedsCleanup)
}

func (suite *FieldTypeSuite) TestGetFieldType() {
	suite.Equal(GetFieldType(reflect.ValueOf("string")), StringFieldType)
	suite.Equal(GetFieldType(reflect.ValueOf(1)), IntegerFieldType)
	suite.Equal(GetFieldType(reflect.ValueOf(uint32(1))), UnsignedFieldType)
	suite.Equal(GetFieldType(reflect.ValueOf(1.0)), FloatFieldType)
	suite.Equal(GetFieldType(reflect.ValueOf(true)), BooleanFieldType)
	suite.Equal(GetFieldType(reflect.ValueOf(time.Now())), TimeFieldType)
	suite.Equal(GetFieldType(reflect.ValueOf(&url.URL{})), URIFieldType)
	suite.Equal(GetFieldType(reflect.ValueOf([]byte{})), UnknownFieldType)
}

func (suite *FieldTypeSuite) TestGetFieldName() {
	names := reflect.TypeOf(struct {
		URI         string
		ResumePoint string `transport:"Resume-Point"`
	}{})
	suite.Equal(GetFieldName(names.Field(0)), "URI")
	suite.Equal(GetFieldName(names.Field(1)), "Resume-Point")
}

func TestFields(test *testing.T) {
	suite.Run(test, new(UnmarshalFieldsSuite))
	suite.Run(test, new(MarshalFieldsSuite))
	suite.Run(test, new(FieldTypeSuite))
}
