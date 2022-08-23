package transport

import (
	"bufio"
	"bytes"
	"embed"
	"io"
	"log"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/stretchr/testify/suite"

	_ "embed"
)

type ScannerSuite struct {
	suite.Suite
}

//go:embed testdata/*.pass testdata/*.fail
var testdata embed.FS

func (suite *ScannerSuite) scan(filepath string) *bufio.Scanner {
	messages, err := testdata.Open(filepath)
	suite.Require().NoError(err)
	scanner := bufio.NewScanner(messages)
	scanner.Split(ScanMessages)
	return scanner
}

func (suite *ScannerSuite) TestScanMessages_Fail() {
	scanner := suite.scan("testdata/0001.capabilities.fail")
	suite.Require().Falsef(scanner.Scan(), "scanner.Scan() returned true")
	suite.Require().ErrorIs(scanner.Err(), io.ErrUnexpectedEOF)
}

// todo: use testdata of messages
func (suite *ScannerSuite) TestScanMessages() {
	scanner := suite.scan("testdata/0001.capabilities.pass")
	suite.Require().Truef(scanner.Scan(), "scanner.Scan() returned false")
	endsWith := strings.HasSuffix(scanner.Text(), "Send-Config: true")
	suite.Require().Truef(endsWith, "scanner.Text() did not end with %q", "Send-Config: true")
}

func FuzzScanMessages(fuzz *testing.F) {
	capabilities, err := testdata.ReadFile("testdata/0001.capabilities.pass")
	if err != nil {
		fuzz.Errorf("failed to open testdata/0001.capabilities.pass: %v", err)
	}
	fuzz.Add(capabilities)
	fuzz.Fuzz(func(test *testing.T, original []byte) {
		scanner := bufio.NewScanner(bytes.NewReader(original))
		scanner.Split(ScanMessages)
		scanner.Scan()
		endsWith := scanner.Text()
		if strings.HasSuffix(endsWith, "\n") {
			test.Errorf("scanner.Text() had a newline at the end")
		}
	})
}

func TestScanner(test *testing.T) {
	suite.Run(test, new(ScannerSuite))
}

func ExampleMessageScanner_Scan() {
	// note the extra newline at the end of the message
	data := heredoc.Doc(`
    100 Capabilities
    Single-Instance: true
    Needs-Cleanup: true
    Send-URI-Encoded: true
    Removable: false
    Send-Config: true

  `)
	scanner := NewMessageScanner(strings.NewReader(data))
	for scanner.Scan() {
		_, err := scanner.Message()
		if err != nil {
			log.Fatal(err)
		}
		// do something with the message
	}
}
