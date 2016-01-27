package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/euank/gotmpl"

	"gopkg.in/yaml.v2"
)

func main() {
	vars := make(map[string]interface{})
	for _, arg := range os.Args[1:] {
		avar := make(map[string]interface{})
		f, err := os.Open(arg)
		if err != nil {
			log.Fatal("Unable to open file: " + arg)
		}
		data, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal("Unable to read file: " + arg)
		}

		yaml.Unmarshal(data, avar)

		for k, v := range avar {
			vars[k] = v
		}
	}

	strVars := make(map[string]string)
	for k, v := range vars {
		if v == nil {
			strVars[k] = ""
		} else {
			strVars[k] = fmt.Sprintf("%v", v)
		}
	}

	if err := gotmpl.Template(os.Stdin, os.Stdout, gotmpl.MapLookup(strVars)); err != nil {
		log.Fatal(err)
	}
}
