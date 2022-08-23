package transport

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MessageWriterSuite struct {
	suite.Suite
}

func (suite *MessageWriterSuite) TestLog() {
	buffer := strings.Builder{}
	writer := NewMessageWriter(&buffer)
	writer.Log("Hello, World!")
	suite.Equal(buffer.String(), "101 Log\nMessage: Hello, World!\n\n")
}

func (suite *MessageWriterSuite) TestLogf() {
	buffer := strings.Builder{}
	writer := NewMessageWriter(&buffer)
	writer.Logf("Hello, %s", "World!")
	suite.Equal(buffer.String(), "101 Log\nMessage: Hello, World!\n\n")
}

func (suite *MessageWriterSuite) TestStatus() {
	buffer := strings.Builder{}
	writer := NewMessageWriter(&buffer)
	writer.Status("Hello, World!")
	suite.Equal(buffer.String(), "102 Status\nMessage: Hello, World!\n\n")
}

func (suite *MessageWriterSuite) TestWarning() {
	buffer := strings.Builder{}
	writer := NewMessageWriter(&buffer)
	writer.Warning("Hello, World!")
	suite.Equal(buffer.String(), "104 Warning\nMessage: Hello, World!\n\n")
}

func TestMessageWriter(test *testing.T) {
	suite.Run(test, new(MessageWriterSuite))
}
