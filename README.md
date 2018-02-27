# gotmpl

A stupid-simple template substitution tool.

## Template language

This code supports the following three features:

1. `${var}` &mdash; a template variable which will be replaced with the value of `var`
1. `\$` &mdash; an escape, which will be replaced with a `$`
1. `\\` &mdash; an escape, which will be replaced with a `\`

No other special syntax is supported.

## Usage

### As a CLI tool

By default, the `gotmpl` cli tool will resolve template variables from the current environment. As an argument, it takes a file to template and prints the result to stdout.

For example:

```sh
$ cat input_file
hello ${VAR}

$ export VAR=world

$ gotmpl input_file
hello world

$ gotmpl input_file > output_file
$ cat output_file
hello world

```

### As a library

A template may be evaluated by providing anything which satisfies the `Lookup` interface. The most trivial thing to use is a go map via the `MapLookup`:

```go
result := gotmpl.TemplateString("foo is ${foo}", gotmpl.MapLookup(map[string]string{"foo": "value"}))

fmt.Println(result)
// Prints: foo is value
```

# See Also

* [envsubst](https://www.gnu.org/software/gettext/manual/html_node/envsubst-Invocation.html)
* [Apache Commons StrSubstitutor](https://commons.apache.org/proper/commons-lang/apidocs/org/apache/commons/lang3/text/StrSubstitutor.html)
* sh
