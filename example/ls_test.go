package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/cucumber/gherkin-go"
)

type lsFeature struct {
	dir string
	buf *bytes.Buffer
}

func lsFeatureContext(s godog.Suite) {
	c := &lsFeature{buf: bytes.NewBuffer(make([]byte, 1024))}

	s.Step(`^I am in a directory "([^"]*)"$`, c.iAmInDirectory)
	s.Step(`^I have a (file|directory) named "([^"]*)"$`, c.iHaveFileOrDirectoryNamed)
	s.Step(`^I run ls$`, c.iRunLs)
	s.Step(`^I should get output:$`, c.iShouldGetOutput)
}

func (f *lsFeature) iAmInDirectory(name string) error {
	f.dir = os.TempDir() + "/" + name
	if err := os.RemoveAll(f.dir); err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.Mkdir(f.dir, 0775)
}

func (f *lsFeature) iHaveFileOrDirectoryNamed(typ, name string) (err error) {
	if len(f.dir) == 0 {
		return fmt.Errorf("the directory was not chosen yet")
	}
	switch typ {
	case "file":
		err = ioutil.WriteFile(f.dir+"/"+name, []byte{}, 0664)
	case "directory":
		err = os.Mkdir(f.dir+"/"+name, 0775)
	}
	return err
}

func (f *lsFeature) iShouldGetOutput(names *gherkin.DocString) error {
	expected := strings.Split(names.Content, "\n")
	actual := strings.Split(strings.TrimSpace(f.buf.String()), "\n")
	if len(expected) != len(actual) {
		return fmt.Errorf("number of expected output lines %d, does not match actual: %d", len(expected), len(actual))
	}
	for i, line := range actual {
		if line != expected[i] {
			return fmt.Errorf(`expected line "%s" at position: %d to match "%s", but it did not`, expected[i], i, line)
		}
	}
	return nil
}

func (f *lsFeature) iRunLs() error {
	f.buf.Reset()
	return ls(f.dir, f.buf)
}
