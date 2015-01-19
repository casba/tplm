package tplm

import (
	"encoding/json"
	"html/template" // what about text/template, is it viable to support both with one package?
	"os"
)

/* structures to hold our parsed config file. */
type TPLM struct {
	TplDir  string
	Root    string
	Helpers []string
	Tpls    []TPL
}

type TPL struct {
	Name  string
	Files []string
}

func Parse(config string) (*TPLM, error) {

	tplm := &TPLM{}

	file, err := os.Open(config)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(file)

	err = dec.Decode(tplm)
	return tplm, err
}

func (tplm *TPLM) Load() map[string]*template.Template {

	templates = make(map[string]*template.Template)

	/* create the root template, it'll be cloned and used to parse all files. */
	templates["root"] = template.Must(template.New("root").Parse(t.Root))

	/* load helper functions into the root handler. */
	template.Must(templates["root"].ParseFiles(t.Helpers...))

	/* for each template defined in the config file, clone root and parse
	the defined files into the given name */
	for _, tpl := range tplm.Tpls {
		root := templates["root"].Clone()
		templates[tpl.Name] = template.Must(root.ParseFiles(tpl.Files...))
	}

	return templates, err
}

func LoadConfig(config string) (map[string]*template.Template, error) {

	tplm, err := Parse(config)
	if err != nil {
		return nil, err
	}

	templates := tplm.Load()
	return templates, nil
}
