# INIParser
# Project Description
An implemenation of a library, ini parser is a Go package for initializing, parsing, manipulating, generating and dealing with INI files.
# How to Use
### 1. import ini parser in your project 
import "https://github.com/codescalersinternships/INIParser-NabilaSherif"
### 2. initialize a new ini parser 
`parser:=NewIniParser()`
### 3. use the methods on the parser
`NewIniParser()` returns a pointer to the initialized parser to be used later with the package's methods
# API implemented 

### 1.`NewIniParser() *IniParser`
Function to initialize a new ini parser and return a pointer to it
### 2.`(p *IniParser) LoadFromString(str string) error`
Method that takes a string supposidly in ini format and loads it to the ini parser
Returns an error if invalid format 
### 3. `(p *IniParser) LoadFromFile(file string) error`
Method that takes file path, reads the file content and loads it to the ini parser
Returns an error if invalid format, not an ini file or error when reading file content
### 4. `(p *IniParser) GetSectionNames() []string`
Method that returns a slice of all sections in the ini parser
### 5. `(p *IniParser) Get(sectionName string, key string) (string, error)`
Method given the section name and key it should return the value of that key for the indicated section
Returns an error in case sectionName is not there or key is not there
### 6. `(p *IniParser) Set(sectionName string, key string, value string) error`
Method given section name , key and value it sets for the indicted section a key with the passed value
If section is not present it creates a new one
Returns an error in case for the indicated section name, there is already a key with the passed key value
### 7. `(p *IniParser) String() string `
Overridding the method String to return a string in the valid ini format
### 8. `(p *IniParser) SaveToFile(filepath string) error`
Method passed the file path, it save to it the content of the ini parser
Returns an error if the file is not an ini one or if an error is encountered when opening the file 
### 9. `(p *IniParser) GetSections() map[string]map[string]string`
Method to return content of a parser 
