package gotmpl

import "testing"

func TestMapLookupTemplate(t *testing.T) {
	s := "echo ${foo}"
	l := map[string]string{"foo": "bar"}

	res, err := TemplateString(s, MapLookup(l))
	if err != nil {
		t.Fatal(err)
	}
	if res != "echo bar" {
		t.Errorf("Expected `echo bar`, got %v", res)
	}
}
