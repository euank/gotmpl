package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/euank/gotmpl"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	env := flag.Bool("env", true, "Pull variables from the environment")
	inplace := flag.Bool("inplace", false, "Replace variables in the given file inplace")
	flag.Parse()
	remainingArgs := flag.Args()

	var tmplReader io.Reader
	var templateFile *os.File
	var inplaceFile string
	defer func() {
		if templateFile != nil {
			_ = templateFile.Close()
		}
	}()

	if shouldReadStdin() {
		if *inplace {
			return errors.New("cannot do inplace replacement of stdin")
		}

		tmplReader = os.Stdin
	} else {
		if len(remainingArgs) == 0 {
			return errors.New("must provide an argument of a file to template")
		}
		fileName := remainingArgs[len(remainingArgs)-1]
		lastFile, err := os.Open(fileName)
		if err != nil {
			return fmt.Errorf("could not open given file (%v) for templating: %w", fileName, err)
		}
		templateFile = lastFile
		tmplReader = lastFile
		if *inplace {
			inplaceFile = fileName
		}
		remainingArgs = remainingArgs[0 : len(remainingArgs)-1]
	}

	resolvers := chainResolver{}

	if *env {
		resolvers = append(resolvers, envLookup{})
	}

	vars := make(map[string]interface{})
	for _, arg := range remainingArgs {
		avar := make(map[string]interface{})
		data, err := os.ReadFile(arg)
		if err != nil {
			return fmt.Errorf("unable to read file %v: %w", arg, err)
		}

		err = json.Unmarshal(data, &avar)
		if err != nil {
			return fmt.Errorf("invalid json file %v: %w", arg, err)
		}

		for k, v := range avar {
			vars[k] = v
		}
	}

	strVars := make(gotmpl.MapLookup)
	for k, v := range vars {
		if v == nil {
			strVars[k] = ""
		} else {
			strVars[k] = fmt.Sprintf("%v", v)
		}
	}

	resolvers = append(resolvers, strVars)

	if inplaceFile == "" {
		return gotmpl.Template(tmplReader, os.Stdout, resolvers)
	}

	var output bytes.Buffer
	if err := gotmpl.Template(tmplReader, &output, resolvers); err != nil {
		return err
	}
	if err := templateFile.Close(); err != nil {
		return fmt.Errorf("close template file: %w", err)
	}
	templateFile = nil
	if err := atomicWriteFile(inplaceFile, output.Bytes()); err != nil {
		return fmt.Errorf("replace template file: %w", err)
	}
	return nil
}

func atomicWriteFile(name string, data []byte) (err error) {
	name, err = filepath.EvalSymlinks(name)
	if err != nil {
		return err
	}
	info, err := os.Stat(name)
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(name), "."+filepath.Base(name)+".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() {
		if tmp != nil {
			_ = tmp.Close()
		}
		_ = os.Remove(tmpName)
	}()

	if err := tmp.Chmod(info.Mode().Perm()); err != nil {
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	tmp = nil
	return os.Rename(tmpName, name)
}

// chainResolver checks for a variable in each element of the chain in order
type chainResolver []gotmpl.Lookup

func (c chainResolver) Resolve(variable string) (string, bool) {
	for _, l := range c {
		if s, ok := l.Resolve(variable); ok {
			return s, ok
		}
	}
	return "", false
}

// envLookup implements a gotmpl.Lookup sources from the environment
type envLookup struct{}

func (envLookup) Resolve(variable string) (string, bool) {
	return os.LookupEnv(variable)
}

// shouldReadStdin determines if stdin should be considered a valid source of data for templating.
func shouldReadStdin() bool {
	stdinStat, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	return stdinStat.Mode()&os.ModeCharDevice == 0
}
