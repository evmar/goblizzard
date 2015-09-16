package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"strconv"
)

type protoField struct {
	Key   int
	Name  string
	Type  string
	Slice bool
}
type protoStruct struct {
	Name   string
	Fields []*protoField
}

func read(path string) (ps []*protoStruct) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		panic(err)
	}

	for _, d := range f.Decls {
		g, ok := d.(*ast.GenDecl)
		if !ok || g.Tok != token.TYPE {
			continue
		}
		for _, s := range g.Specs {
			s := s.(*ast.TypeSpec)
			t, ok := s.Type.(*ast.StructType)
			if !ok {
				continue
			}
			p := protoStruct{
				Name: s.Name.Name,
			}

			for _, f := range t.Fields.List {
				if f.Tag == nil {
					continue
				}
				pf := &protoField{}
				tag := f.Tag.Value
				tag = tag[1 : len(tag)-1] // strip "
				pf.Key, err = strconv.Atoi(tag)
				if err != nil {
					panic(err)
				}
				pf.Name = f.Names[0].Name
				switch t := f.Type.(type) {
				case *ast.Ident:
					pf.Type = fmt.Sprintf("%s", t)
				case *ast.ArrayType:
					pf.Type = fmt.Sprintf("%s", t.Elt.(*ast.StarExpr).X)
					pf.Slice = true
				default:
					panic(t)
				}
				p.Fields = append(p.Fields, pf)
			}
			if len(p.Fields) > 0 {
				ps = append(ps, &p)
			}
		}
	}
	return
}

type w struct {
	io.Writer
}

func (w *w) Print(s string, args ...interface{}) {
	fmt.Fprintf(w, s+"\n", args...)
}

func main() {
	inf := flag.String("in", "", "input filename")
	flag.Parse()

	ps := read(*inf)
	w := w{os.Stdout}
	w.Print("package replay")
	w.Print("import \"blizzard/blizzval\"")
	for _, p := range ps {
		log.Printf("%s", p.Name)
		w.Print("func read%s(v blizzval.Value) *%s {", p.Name, p.Name)
		w.Print("out := %s{}", p.Name)
		w.Print("out.Raw = v")
		w.Print("m := v.(map[int]blizzval.Value)")
		for _, f := range p.Fields {
			if f.Slice {
				w.Print("if f, ok := m[%d]; ok {", f.Key)
				w.Print("s := f.([]blizzval.Value)")
				w.Print("out.%s = make([]*%s, len(s))", f.Name, f.Type)
				w.Print("for i := 0; i < len(s); i++ {")
				w.Print("out.%s[i] = read%s(s[i])", f.Name, f.Type)
				w.Print("}")
				w.Print("}")
			} else {
				w.Print("out.%s = m[%d].(%s)", f.Name, f.Key, f.Type)
			}
		}
		w.Print("return &out")
		w.Print("}")
	}
}
