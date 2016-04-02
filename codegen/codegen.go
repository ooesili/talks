package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"strings"
	"text/template"
)

func main() {
	if err := realMain(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func realMain() error {
	// parse START OMIT
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "fixture/structs.go", nil, parser.ParseComments)
	if err != nil {
		return err
	}
	files := []*ast.File{f}

	info := types.Info{
		Defs: make(map[*ast.Ident]types.Object),
	}
	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check("github.com/ooesili/codegen", fset, files, &info)
	if err != nil {
		return err
	}
	// parse END OMIT

	// collect getters START OMIT
	getters := make([]Getter, 0)
	for _, decl := range f.Decls {
		getters = append(getters, findGetters(decl, info)...)
	}
	// collect getters END OMIT

	// output fmt START OMIT
	buf := &bytes.Buffer{}
	err = GettersTemplate.Execute(buf, GettersInfo{
		Package: pkg.Name(),
		Getters: getters,
	})
	if err != nil {
		return err
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	fmt.Printf("%s", formatted)
	// output fmt END OMIT

	return nil
}

func hasComment(decl *ast.GenDecl, text string) bool {
	if decl.Doc == nil {
		return false
	}

	for _, comment := range decl.Doc.List {
		if comment.Text == text {
			return true
		}
	}

	return false
}

// find genDecl START OMIT
func findGetters(decl ast.Decl, info types.Info) []Getter {
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok {
		return nil
	}
	// find genDecl END OMIT

	// check comment START OMIT
	if !hasComment(genDecl, "// +getters") {
		return nil
	}
	// check comment END OMIT

	// struct decl START OMIT
	typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec)
	if !ok {
		return nil
	}

	obj := info.Defs[typeSpec.Name]
	def, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return nil
	}
	// struct decl END OMIT

	// collect fields START OMIT
	getters := make([]Getter, def.NumFields())

	for i := 0; i < def.NumFields(); i++ {
		v := def.Field(i)
		getters[i] = Getter{
			Name:     v.Name(),
			Reciever: typeSpec.Name.Name,
			Type:     v.Type().String(),
		}
	}

	return getters
	// collect fields END OMIT
}

// template data START OMIT
type GettersInfo struct {
	Package string
	Getters []Getter
}

type Getter struct {
	Reciever string
	Name     string
	Type     string
}

func (g Getter) RecieverArg() string {
	return strings.ToLower(g.Reciever[0:1])
}

func (g Getter) ExportedName() string {
	return strings.ToUpper(g.Name[0:1]) + g.Name[1:]
}

// template data END OMIT

// template START OMIT
var GettersTemplate = template.Must(template.New("getters").Parse(
	`package {{.Package}}

{{range .Getters}}
// {{.ExportedName}} is a getter for {{.Reciever}}
func ({{.RecieverArg}} {{.Reciever}}) {{.ExportedName}}() {{.Type}} {
	return {{.RecieverArg}}.{{.Name}}
}
{{end}}
`))

// template END OMIT
