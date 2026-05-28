package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestHelloCmd_Default(t *testing.T) {
	cmd := HelloCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs(nil)
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if got := buf.String(); !strings.Contains(got, "Olá, mundo!") {
		t.Errorf("esperava saudação default, recebi: %q", got)
	}
}

func TestHelloCmd_ComNome(t *testing.T) {
	cmd := HelloCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"eco"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if got := buf.String(); !strings.Contains(got, "Olá, eco!") {
		t.Errorf("esperava saudação com nome, recebi: %q", got)
	}
}
