package iniparser

//importing the needed packages
import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

// defining the errors that may happen
var ErrNotINI = errors.New("not an INI file")
var ErrSectionNotThere = errors.New("section not there")
var ErrKeyAlreadyExists = errors.New("for the provided section, the provided key already exists")
var ErrKeyNotThere = errors.New("key not there")
var ErrReadingFile = errors.New("can't read file")
var ErrOpeningFile = errors.New("can't open file")
var ErrInvalidContent = errors.New("invalid content to be stored in an ini parser")
var ErrSectionAlreadyThere = errors.New("there is already a section with the provided name")
var ErrWritingToFile = errors.New("problem when writing to file")
var ErrEmptyKeyORValue = errors.New("key or value is empty")

// defining a data type for the ini parser
type IniParser struct {
	data map[string](map[string]string)
}

// initialize a new ini parser and returning a pointer to it
func NewIniParser() *IniParser {
	myIniParser := IniParser{make(map[string]map[string]string)}
	return &myIniParser
}

// given a string, call the helper to check for valid ini parser content and add the content to the parser
func (p *IniParser) LoadFromString(str string) error {
	err := p.parse(str)
	return err
}

// given a file path , call the helper to check for valid ini parser content and add the content to the parser
func (p *IniParser) LoadFromFile(file string) error {
	if !strings.HasSuffix(file, ".ini") {
		return ErrNotINI
	}
	data, e := os.ReadFile(file)
	if e != nil {
		return ErrReadingFile
	}
	stringData := string(data)
	e = p.parse(stringData)
	return e
}

// helper for LoadFromString and LoadFromFile to check for valid content and to add to the parser
func (p *IniParser) parse(data string) error {
	lines := strings.Split(data, "\n")
	var section map[string]string
	var sectionName string
	for _, l := range lines {
		line := strings.TrimSpace(l)
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}
		if !(strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]")) && !strings.Contains(line, "=") {
			return ErrInvalidContent
		}
		if strings.Contains(line, "[") && strings.Contains(line, "]") {
			sectionName = line[1:(len(line) - 1)]
			if p.data[sectionName] != nil {
				return ErrSectionAlreadyThere
			}
			section = make(map[string]string)
			p.data[sectionName] = section
		} else {
			currentKandV := strings.SplitN(line, "=", 2)
			if slices.Contains(currentKandV, "") {
				return ErrEmptyKeyORValue
			}
			key := strings.TrimSpace(currentKandV[0])
			value := strings.TrimSpace(currentKandV[1])
			values := p.data[sectionName]
			if _, ok := values[key]; ok {
				return ErrKeyAlreadyExists
			}
			values[key] = value
		}
	}
	return nil
}

// return the sections names in a slice of strings
func (p *IniParser) GetSectionNames() []string {
	var sectionNames []string
	for sectionName := range p.data {
		sectionNames = append(sectionNames, sectionName)
	}
	return sectionNames
}

// retrieve all sections along with its content
func (p *IniParser) GetSections() map[string]map[string]string {
	return p.data
}

// given section name and a key, return the corresponding value if present
func (p *IniParser) Get(sectionName string, key string) (string, error) {
	section, ok := p.data[sectionName]
	if !ok {
		return "", ErrSectionNotThere
	}
	value, ok := section[key]
	if !ok {
		return "", ErrKeyNotThere
	}
	return value, nil
}

// given a section name and a key and a value, set the key of the section with the given name to the given value
func (p *IniParser) Set(sectionName string, key string, value string) error {
	section, ok := p.data[sectionName]
	if !ok {
		section = make(map[string]string)
		p.data[sectionName] = section
	}
	_, ok = section[key]
	if ok {
		return ErrKeyAlreadyExists
	}
	section[key] = value
	return nil
}

// to return ini parser content in the valid format
func (p *IniParser) String() string {
	s := ""
	for sectionName, section := range p.data {
		s += fmt.Sprintf("[%s]\n", sectionName)
		// fmt.Println("***", sectionName)
		for k, v := range section {
			// fmt.Println("***", k, v)
			s += fmt.Sprintf("%s=%s\n", k, v)
		}
	}
	return strings.TrimSpace(s)
}

// save the content of the ini parser to a file given its path
func (p *IniParser) SaveToFile(filepath string) error { //.data bas w create file w return haga?
	if !strings.HasSuffix(filepath, ".ini") {
		return ErrNotINI
	}
	s := p.String()
	f, e := os.Open(filepath)
	if e != nil {
		return ErrOpeningFile
	}
	defer f.Close()
	_, e = f.WriteString(s)
	if e != nil {
		return ErrWritingToFile
	}
	return nil
}
