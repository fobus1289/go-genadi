package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"text/template"

	"github.com/fobus1289/go-genadi/gen"
	"golang.org/x/tools/imports"
)

var logger = log.New(os.Stdout, "Info:", log.Lshortfile)

//go:embed stubl/router.gen
var routerTemplate string

type genIt struct {
	genMethods  []gen.METHODS
	genServices []gen.METHODS
	objects     gen.Objects
	Functions   gen.Functions
}

func GetPattern(ifs gen.IFS) string {
	if ifs.GetServiceParams() != "" || ifs.GetPatternName() == "" {
		return ""
	}
	return fmt.Sprintf("%s = regexp.MustCompile(`%s`)\n",
		ifs.GetPatternName(),
		getRegular(ifs.GetPath()),
	)
}

func (g *genIt) MakeControllers() string {
	var builder = &strings.Builder{}

	for _, c := range g.genMethods {
		for _, object := range c.GetObjects() {
			var _type = fmt.Sprintf("&%s.%s{}", object.Pkg(), object.Name())
			var _var = fmt.Sprintf("%s=%s\n", object.GetValueName(), _type)
			builder.WriteString(_var)
		}
	}

	return builder.String()
}

func (g *genIt) MakeServices() string {
	var builder = &strings.Builder{}

	for _, c := range g.genServices {
		for _, object := range c.GetObjects() {
			if object.Name() == "ResponseWriter" || object.Name() == "Request" || object.GetType() == gen.Configure {
				continue
			}
			var _type = fmt.Sprintf("&%s.%s{}", object.Pkg(), object.Name())
			var _var = fmt.Sprintf("%s=%s\n", object.GetValueName(), _type)
			builder.WriteString(_var)
		}
	}

	return builder.String()
}

func (g *genIt) MakeActionPattern() string {
	var builder = &strings.Builder{}

	for _, genMethod := range g.genMethods {
		for _, function := range genMethod.GetMethods() {
			builder.WriteString(GetPattern(function))
		}
	}

	return builder.String()
}

func (g *genIt) GetConfigure() string {

	var (
		builder = &strings.Builder{}
	)

	for _, fn := range g.Functions {

		if fn.Type != 1 {
			continue
		}

		builder.WriteString(fn.GetConfiguration())
	}

	return builder.String()
}

func (g *genIt) GetServicesParams() string {

	var (
		builder = &strings.Builder{}
	)

	for _, method := range g.objects {
		for _, _if := range method.Methods() {
			if result := _if.GetServiceParams(); result != "" {
				builder.WriteString(result)
			}
		}
	}

	return builder.String()
}

func (g *genIt) GetMethods() string {

	var (
		builder     = &strings.Builder{}
		actions     = map[string][]gen.IFS{}
		httpMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	)

	for _, key := range httpMethods {
		for _, genMethod := range g.genMethods {
			if its := genMethod.GetMethod(key); len(its) != 0 {
				actions[key] = append(actions[key], its...)
			}
		}
	}

	for key, methods := range actions {

		switch key {
		case "GET":
			builder.WriteString("\ncase http.MethodGet:\n")
		case "POST":
			builder.WriteString("\ncase http.MethodPost:\n")
		case "PUT":
			builder.WriteString("\ncase http.MethodPut:\n")
		case "PATCH":
			builder.WriteString("\ncase http.MethodPatch:\n")
		case "DELETE":
			builder.WriteString("\ncase http.MethodDelete:\n")
		}

		var validSlice []string

		for _, method := range methods {
			validSlice = append(validSlice, method.GetIf())
		}

		builder.WriteString(strings.Join(validSlice, " else "))
		builder.WriteString("\n")
		//builder.WriteString("\nreturn")
	}

	return builder.String()
}

func te() {
	fset := token.NewFileSet()

	src := `package main
		type A interface {
		
		}

		type B struct {
			
		}
	`

	astFile, err := parser.ParseFile(fset, "", src, parser.ParseComments)

	if err != nil {
		log.Fatalln(err)
	}

	for _, decl := range astFile.Decls {
		switch t := decl.(type) {
		case *ast.GenDecl:
			var spec = t.Specs[0]
			switch s := spec.(type) {
			case *ast.TypeSpec:
				switch st := s.Type.(type) {
				case *ast.InterfaceType:
					log.Println(st.Methods.List)
				default:
					log.Println(reflect.TypeOf(st))
				}
			default:
				log.Println(reflect.TypeOf(s))
			}
		default:
			log.Println(reflect.TypeOf(t))
		}

	}

}

func main() {

	var pars = &gen.Parser{
		Once:        &sync.Once{},
		Controllers: gen.Objects{},
		Services:    gen.Objects{},
		Functions:   gen.Functions{},
	}

	if err := pars.Parse(); err != nil {
		logger.Fatalln(err)
	}

	var genIts = &genIt{
		genMethods: []gen.METHODS{
			pars.Controllers, pars.Functions,
		},
		genServices: []gen.METHODS{
			pars.Services, pars.Functions,
		},
		objects:   append(pars.Controllers, pars.Services...),
		Functions: pars.Functions,
	}

	var tmp = template.Must(template.New("").Parse(routerTemplate))
	var bu = &bytes.Buffer{}
	tmp.Execute(bu, genIts)

	var re, err2 = imports.Process("", bu.Bytes(), nil)

	if err2 != nil {
		logger.Fatalln(err2)
	}

	var fd1, fdErr1 = os.Create("example/router/init.go")

	if fdErr1 != nil {
		logger.Fatalln(fdErr1)
	}

	defer fd1.Close()

	if _, writeErr1 := fd1.Write(re); writeErr1 != nil {
		logger.Fatalln(writeErr1)
	}

	logger.Println("success generate router")

}

func getRegular(urlPath string) string {
	urlPath = path.Clean(urlPath)
	{
		if urlPath == "." || urlPath == "/" {
			return "^(/)$"
		}
		urlPath = strings.TrimPrefix(urlPath, "/")
	}

	if index := strings.Index(urlPath, "*"); index != -1 {
		urlPath = urlPath[0 : index+1]
	}
	urlPath = strings.Replace(urlPath, "*", "(.*)?", -1)
	var regular = regexp.MustCompile(`(:[a-zA-Z]+)`)
	urlPath = regular.ReplaceAllString(urlPath, `([0-9a-zA-Z]+)`)
	return fmt.Sprintf("^(/?%s/?)$", urlPath)
}
