package iniparser

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestLoadFromString(t *testing.T) {
	t.Run("correct loading", func(t *testing.T) {
		input := `
[owner]
name = John Doe
organization = Acme Widgets Inc.
[database]
server = 192.0.2.62     
port = 143
file = payroll.dat`
		parser := NewIniParser()
		e := parser.LoadFromString(input)
		if e != nil {
			t.Errorf("got %q want %q", e, "nil")
		}
		m := map[string]map[string]string{
			"owner": {
				"name":         "John Doe",
				"organization": "Acme Widgets Inc.",
			},
			"database": {
				"server": "192.0.2.62",
				"port":   "143",
				"file":   "payroll.dat",
			},
		}
		if !reflect.DeepEqual(parser.GetSections(), m) {
			t.Errorf("got %q want %q", parser.GetSections(), m)
		}

	})
	t.Run("Invalid ini content", func(t *testing.T) {
		input := `
[owner]
name = John Doe
organization = Acme Widgets Inc.
database
server = 192.0.2.62     
port = 143
file = payroll.dat`
		parser := NewIniParser()
		e := parser.LoadFromString(input)
		if e != ErrInvalidContent {
			t.Errorf("got %q want %q", e, ErrInvalidContent)
		}
	})
	t.Run("Empty key or value", func(t *testing.T) {
		input := `
[owner]
= John Doe
organization = Acme Widgets Inc.
database
server = 192.0.2.62     
port = 143
file = payroll.dat`
		parser := NewIniParser()
		e := parser.LoadFromString(input)
		if e != ErrEmptyKeyORValue {
			t.Errorf("got %q want %q", e, ErrEmptyKeyORValue)
		}
	})
	t.Run("duplicate section name", func(t *testing.T) {
		input := `
[owner]
name = John Doe
organization = Acme Widgets Inc.
[owner]
server = 192.0.2.62     
port = 143
file = payroll.dat`
		parser := NewIniParser()
		e := parser.LoadFromString(input)
		if e != ErrSectionAlreadyThere {
			t.Errorf("got %q want %q", e, ErrSectionAlreadyThere)
		}
	})
	t.Run("loading duplicate key for same section", func(t *testing.T) {
		input := `
[owner]
name = John Doe
organization = Acme Widgets Inc.
[database]
server = 192.0.2.62     
server = 143
file = payroll.dat`
		parser := NewIniParser()
		e := parser.LoadFromString(input)
		if e != ErrKeyAlreadyExists {
			t.Errorf("got %q want %q", e, ErrKeyAlreadyExists)
		}
	})
}
func TestLoadFromFile(t *testing.T) {
	//check on ini else nafs el f loadfromstring
	t.Run("not an ini file", func(t *testing.T) {
		parser := NewIniParser()
		e := parser.LoadFromFile("iniparser/file.txt")
		if e != ErrNotINI {
			t.Errorf("got %q want %q", e, ErrNotINI)
		}
	})
	t.Run("error when reading file", func(t *testing.T) {
		parser := NewIniParser()
		e := parser.LoadFromFile("iniparser/file.ini")
		if e != ErrReadingFile {
			t.Errorf("got %q want %q", e, ErrReadingFile)
		}
	})
	t.Run("otherwise it works", func(t *testing.T) {
		parser := NewIniParser()
		e := parser.LoadFromFile("/root/INIParser-NabilaSherif/pkg/testdata/toloadfrom.ini")
		if e != nil {
			t.Errorf("got %q want %q", e, "nil")
		}
		m := map[string]map[string]string{
			"owner": {
				"name":         "John Doe",
				"organization": "Acme Widgets Inc.",
			},
			"database": {
				"server": "192.0.2.62",
				"port":   "143",
				"file":   "payroll.dat",
			},
		}
		if !reflect.DeepEqual(parser.GetSections(), m) {
			t.Errorf("got %q want %q", parser.GetSections(), m)
		}
	})
	//this subtest is enough as for all other cases, the same method is used for loading to string and file
}
func TestGetSectionName(t *testing.T) {
	input := `
[owner]
name = John Doe
organization = Acme Widgets Inc.
[database]
server = 192.0.2.62     
port = 143
file = payroll.dat`
	parser := NewIniParser()
	e := parser.LoadFromString(input)
	if e != nil {
		t.Errorf("got %q want %q", e, "nil")
	}
	got := parser.GetSectionNames()
	want := []string{"owner", "database"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q want %q", got, want)
	}
}
func TestGet(t *testing.T) {
	t.Run("section is not there", func(t *testing.T) {
		input := `
[owner]
name = John Doe
organization = Acme Widgets Inc.
[database]
server = 192.0.2.62     
port = 143
file = payroll.dat`
		parser := NewIniParser()
		e := parser.LoadFromString(input)
		if e != nil {
			t.Errorf("got %q want %q", e, "nil")
		}
		_, got := parser.Get("dumm", "key")
		if got != ErrSectionNotThere {
			t.Errorf("got %q want %q", got, ErrSectionNotThere)
		}
	})
	t.Run("section is there, but key is not there", func(t *testing.T) {
		input := `
[owner]
name = John Doe
organization = Acme Widgets Inc.
[database]
server = 192.0.2.62     
port = 143
file = payroll.dat`
		parser := NewIniParser()
		e := parser.LoadFromString(input)
		if e != nil {
			t.Errorf("got %q want %q", e, "nil")
		}
		_, got := parser.Get("owner", "key")
		if got != ErrKeyNotThere {
			t.Errorf("got %q want %q", got, ErrKeyNotThere)
		}
	})
	t.Run("section is there and key is there", func(t *testing.T) {
		input := `
[owner]
name =John Doe
organization = Acme Widgets Inc.
[database]
server = 192.0.2.62     
port = 143
file = payroll.dat`
		parser := NewIniParser()
		e := parser.LoadFromString(input)
		if e != nil {
			t.Errorf("got %q want %q", e, "nil")
		}
		data, got := parser.Get("owner", "name")
		if got != nil {
			t.Errorf("got %q as an error, want %q", got, "nil")
		}
		if data != "John Doe" {
			t.Errorf("got %q , want %q", data, "John Doe")
		}
	})
}
func TestSet(t *testing.T) {
	t.Run("section is there already and new key", func(t *testing.T) {
		input := `
[owner]
name = John Doe
organization = Acme Widgets Inc.
[database]
server = 192.0.2.62     
port = 143
file = payroll.dat`
		parser := NewIniParser()
		e := parser.LoadFromString(input)
		if e != nil {
			t.Errorf("got %q want %q", e, "nil")
		}
		got := parser.Set("owner", "year", "2019")
		if got != nil {
			t.Errorf("got %q as an error, want %q", got, "nil")
		}
		section := parser.data["owner"]
		value := section["year"]
		if value != "2019" {
			t.Errorf("got %q , want %q", value, "2019")
		}
	})
	t.Run("section not there", func(t *testing.T) {
		input := `
[owner]
name = John Doe
organization = Acme Widgets Inc.
[database]
server = 192.0.2.62     
port = 143
file = payroll.dat`
		parser := NewIniParser()
		e := parser.LoadFromString(input)
		if e != nil {
			t.Errorf("got %q want %q", e, "nil")
		}
		got := parser.Set("section", "key", "value")
		if got != nil {
			t.Errorf("got %q as an error, want %q", got, "nil")
		}
		section := parser.data["section"]
		value := section["key"]
		if value != "value" {
			t.Errorf("got %q , want %q", value, "2019")
		}
	})
	t.Run("section there but key already exists", func(t *testing.T) {
		input := `
[owner]
name = John Doe
organization = Acme Widgets Inc.
[database]
server = 192.0.2.62     
port = 143
file = payroll.dat`
		parser := NewIniParser()
		e := parser.LoadFromString(input)
		if e != nil {
			t.Errorf("got %q want %q", e, "nil")
		}
		got := parser.Set("owner", "name", "nabila")
		if got != ErrKeyAlreadyExists {
			t.Errorf("got %q as an error, want %q", got, ErrKeyAlreadyExists)
		}
	})
}

func TestToString(t *testing.T) {
	t.Run("correct String() behaviour", func(t *testing.T) {
		input := `[owner]
name=John Doe
organization=Acme Widgets Inc.
[database]
server=192.0.2.62
port=143
file=payroll.dat`
		parser := NewIniParser()
		e := parser.LoadFromString(input)
		if e != nil {
			t.Errorf("got %q want %q", e, "nil")
		}
		strOutput := parser.String()
		splittedStrOut := strings.Split(strOutput, "\n")
		fmt.Println((len(splittedStrOut)))
		input = strings.ReplaceAll(input, `" "`, "")
		splittedInput := strings.Split(input, "\n")
		if !reflect.DeepEqual(splittedInput, splittedStrOut) {
			t.Errorf("got %q , want %q", splittedStrOut, splittedInput)
		}
	})
}
func TestSaveToFile(t *testing.T) {
	t.Run("save to an ini", func(t *testing.T) {
		path := "/root/INIParser-NabilaSherif/iniparser/file2.ini"
		input := `
[owner]
name = John Doe
organization = Acme Widgets Inc.
[database]
server = 192.0.2.62     
port = 143
file = payroll.dat`
		parser := NewIniParser()
		e := parser.LoadFromString(input)
		if e != nil {
			t.Errorf("got %q want %q", e, "nil")
		}
		got := parser.SaveToFile(path)
		if got != nil {
			t.Errorf("got %q as an error, want %q", got, "nil")
		}
	})
	t.Run("saving to a diffrent extension than ini", func(t *testing.T) {
		path := "iniparser/file.txt"
		input := `
[owner]
name = John Doe
organization = Acme Widgets Inc.
[database]
server = 192.0.2.62     
port = 143
file = payroll.dat`
		parser := NewIniParser()
		e := parser.LoadFromString(input)
		if e != nil {
			t.Errorf("got %q want %q", e, "nil")
		}
		got := parser.SaveToFile(path)
		if got != ErrNotINI {
			t.Errorf("got %q as an error, want %q", got, ErrNotINI)
		}
	})
	t.Run("saving to a non existant ini file", func(t *testing.T) {
		path := "iniparser/nonexistant.ini"
		input := `
[owner]
name = John Doe
organization = Acme Widgets Inc.
[database]
server = 192.0.2.62     
port = 143
file = payroll.dat`
		parser := NewIniParser()
		e := parser.LoadFromString(input)
		if e != nil {
			t.Errorf("got %q want %q", e, "nil")
		}
		got := parser.SaveToFile(path)
		if got != ErrOpeningFile {
			t.Errorf("got %q as an error, want %q", got, ErrOpeningFile)
		}
	})
}
