package tplm

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var CONFIG_FORMAT = `{ 	
	"TplDir": "%s", 
	"Root":"layout.html",
	"Tpls": [
		{
			"Name":"poll_view",
			"Files":[
				"view.html"
			]
		},
		{
			"Name":"poll_new",
			"Files":[
				"new.html",
				"message.html"
			]
		}
	]
}`

var config string

type templatefiles struct {
	name     string
	contents string
}

var templates = []templatefiles{

	{"view.html", `{{define "content"}}<p>it's view.</p>{{end}}`},
	{"new.html", `{{define "content"}}<div class="message">{{template "message"}}</div><p>it's new.</p>{{end}}`},
	{"message.html", `{{define "message"}}<p>it's a message</p>{{end}}`},
	{"layout.html", `{{define "root"}}<html><head><title>hi</title></head><body>{{template "content"}}</body></html>{{end}}`},
}

func TestLoadConfig(t *testing.T) {
	dir := createTestDir(templates, CONFIG_FORMAT)
	defer os.RemoveAll(dir)

	t.Log("Starting TestLoadConfig")
	tplm, err := LoadConfig(filepath.Join(dir, "config.json"))
	if err != nil {
		t.Fatal(err)
	}

	for key, _ := range tplm.templates {
		// just pump the rendered template into dev/null, we're not really concerned with the output for this test,
		// just that we can generate our template map and execute them without error.
		err = tplm.templates[key].ExecuteTemplate(ioutil.Discard, "root", nil)
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Log("TestLoadConfig completed successfully.")
}

func write(fpath string, contents string) error {
	f, err := os.Create(fpath)
	if err != nil {
		log.Print("Failed to create file at: " + fpath)
		return err
	}
	defer f.Close()
	_, err = io.WriteString(f, contents)
	if err != nil {
		log.Print("Failed to write file")
		return err
	}
	return nil
}

func createTestDir(files []templatefiles, format string) string {
	dir, err := ioutil.TempDir("", "tpl")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		err = write(filepath.Join(dir, file.name), file.contents)
		if err != nil {
			log.Fatal(err)
		}
	}

	config := fmt.Sprintf(format, dir)
	err = write(filepath.Join(dir, "config.json"), config)
	if err != nil {
		log.Fatal(err)
	}

	return dir
}
