package actioninfo

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

func TestInfoContinuesAfterErrors(t *testing.T) {
	parser := new(MockDataParser)
	parser.On("Parse", "bad parse").Return(errors.New("parse failure")).Once()
	parser.On("Parse", "bad info").Return(nil).Once()
	parser.On("ActionInfo").Return("must not print", errors.New("info failure")).Once()
	parser.On("Parse", "valid").Return(nil).Once()
	parser.On("ActionInfo").Return("valid report", nil).Once()

	var logs bytes.Buffer
	previousWriter := log.Writer()
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(previousWriter) })

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	previousStdout := os.Stdout
	os.Stdout = w
	t.Cleanup(func() {
		os.Stdout = previousStdout
		r.Close()
		w.Close()
	})

	Info([]string{"bad parse", "bad info", "valid"}, parser)
	w.Close()
	output, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != "valid report\n" {
		t.Fatalf("неожиданный вывод: %q", output)
	}
	for _, message := range []string{"parse failure", "info failure"} {
		if !strings.Contains(logs.String(), message) {
			t.Errorf("ошибка %q отсутствует в журнале: %s", message, logs.String())
		}
	}
	parser.AssertExpectations(t)
}
