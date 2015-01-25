package tplm

import (
	"encoding/json"
	"fmt"
	"html/template" // what about text/template, is it viable to support both with one package?
	"os"
	"path"
)

/* structures to hold our parsed config file. */
type TPLC struct {
	TplDir  string
	Root    string
	Helpers []string
	Tpls    []TPL
}

type TPL struct {
	Name  string
	Files []string
}

type TplmError struct {
	Where string
	Err   error
}

type TPLM struct {
	Templates map[string]*template.Template
	Root      *template.Template
}

func (e TplmError) Error() string {
	return fmt.Sprintf("%v: %v", e.Where, e.Err)
}

const (
	PARSE_ERROR  string = "TPLC error in Parse"
	APPEND_ERROR string = "TPLC error in appendRoot"
	LOAD_ERROR   string = "TPLC error in Load"
)

func Parse(config string) (*TPLC, error) {

	tplc := &TPLC{}

	file, err := os.Open(config)
	defer file.Close()
	if err != nil {
		return nil, TplmError{PARSE_ERROR, err}
	}
	dec := json.NewDecoder(file)

	err = dec.Decode(tplc)
	if err != nil {
		err = TplmError{PARSE_ERROR, err}
	}
	return tplc, err
}

/* add the TplDir string to the paths if it's supplied */
func (tplc *TPLC) appendRoot() {
	if tplc.TplDir == "" {
		return
	}

	tplc.Root = path.Join(tplc.TplDir, tplc.Root)
	for idx, _ := range tplc.Helpers {
		tplc.Helpers[idx] = path.Join(tplc.TplDir, tplc.Helpers[idx])
	}
	for _, tpl := range tplc.Tpls {
		for idx, _ := range tpl.Files {
			tpl.Files[idx] = path.Join(tplc.TplDir, tpl.Files[idx])
		}
	}
}

func (tpl TPL) String() string {
	return fmt.Sprintf("\n\t{\n\t\tname:%s, \n\t\tfiles:{\n\t\t\t%s\n\t\t}\n\t}\n\t", tpl.Name, tpl.Files)
}

func (t TPLC) String() string {
	var output = ""

	output = fmt.Sprintf("{\nroot:%s, \nhelpers:%v, \ntpls:%v\n}", t.Root, t.Helpers, t.Tpls)
	return output
}

/* get rid of these panics later for actual error returns. */
func (tplc *TPLC) Load() (*TPLM, error) {

	var err error
	tplm := &TPLM{}

	tplm.Templates = make(map[string]*template.Template)

	// adjust our file definitions to take into account the TplDir option from
	//	the config
	tplc.appendRoot()
	/* create the root template, it'll be cloned and used to parse all files. */
	tplm.Root, err = template.New("root").Parse(tplc.Root)
	if err != nil {
		return nil, TplmError{LOAD_ERROR, err}
	}

	/* load helper functions into the root handler. */
	if len(tplc.Helpers) > 0 {
		_, err = tplm.Root.ParseFiles(tplc.Helpers...)
		if err != nil {
			return nil, TplmError{LOAD_ERROR, err}
		}
	}

	// for each template defined in the config file, clone root and parse
	// the defined files into the given name
	for _, tpl := range tplc.Tpls {

		if len(tpl.Files) <= 0 {
			return nil, TplmError{LOAD_ERROR, fmt.Errorf("template definition:%s has no files defined", tpl.Name)}
		}

		root, err := tplm.Root.Clone()
		if err != nil {
			return nil, TplmError{LOAD_ERROR, err}
		}
		tplm.Templates[tpl.Name], err = root.ParseFiles(tpl.Files...)
		if err != nil {
			return nil, TplmError{LOAD_ERROR, err}
		}
	}

	return tplm, nil
}

func LoadConfig(config string) (*TPLM, error) {

	tplc, err := Parse(config)
	if err != nil {
		return nil, err
	}

	tplm, err := tplc.Load()
	if err != nil {
		return nil, err
	}
	return tplm, nil
}
