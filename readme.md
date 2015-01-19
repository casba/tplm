A simple template manager, load a config file and parse it into a map of named templates for use in rendering pages.


Read in a config.json file of the format:

{
	"root":  #the root template, all templates are compiled by cloning this files parse
	"tpldir": #the root dir of the templates
	"helpers": #array of filenames, compiled into root (css.tpl, js.tpl, etc...), optional
	"tpls": # array of template definitions
		{ #template example
			"name": #the name you want to refer to this template by
			"files": #files to compile into root
		}
}

See tplm_test.go for an example of a config file.

Don't currently support adding functions to the templates.