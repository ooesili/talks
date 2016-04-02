package main

import (
	"os"
	"text/template"
)

// START OMIT
const tmplString = `package main

import "fmt"

func main() {
{{- range .}}
	fmt.Println({{printf "%q" .}})
{{- end}}
}
`

func main() {
	tmpl, _ := template.New("gen").Parse(tmplString)
	tmpl.Execute(os.Stdout, []string{"Hello", "world!"}) // HL
}

// END OMIT
