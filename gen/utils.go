package gen

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var logger = log.New(os.Stdout, "Info:", log.Lshortfile)

type METHODS interface {
	GetMethod(key string) []IFS
	GetMethods() []IFS
	GetObjects() []IObject
}

type IObject interface {
	GetValueName() string
	GetType() GenKind
	Dir() string
	Pkg() string
	Name() string
	Url() string
}

type IFS interface {
	GetIf() string
	GetPath() string
	CallMethod() string
	GetPatternName() string
	GetServiceParams() string
}

func Getwd() string {
	if path, err := os.Getwd(); err != nil {
		return ""
	} else {
		return path
	}
}

func IsExported(decl ast.Decl) bool {
	switch t := decl.(type) {
	case *ast.GenDecl:
		if spec, ok := t.Specs[0].(*ast.TypeSpec); ok {
			return ast.IsExported(spec.Name.Name)
		}
		return false
	case *ast.FuncDecl:
		return ast.IsExported(t.Name.Name)
	}

	return false
}

func hasComment(group *ast.CommentGroup, pattern string) (string, string, bool) {
	if group == nil || len(group.List) == 0 {
		return "", "", false
	}

	var list = group.List

	var reg = regexp.MustCompile(pattern)

	for _, comment := range list {
		var text = trim(comment.Text)

		if find := reg.FindStringSubmatch(text); len(find) == 3 {
			var (
				arg      = trimPattern(find[2])
				function = trimPattern(find[1])
			)

			return arg, function, true
		}
		return "", "", false
	}

	return "", "", false
}

func trim(s string) string {
	return strings.Replace(s, " ", "", -1)
}

func trimPattern(s string) string {
	return regexp.MustCompile(`(\(|\)|@|\s)`).ReplaceAllString(s, "")
}

func parseGoFiles(fset *token.FileSet) (map[string]*ast.Package, error) {

	var pkgs = make(map[string]*ast.Package)

	return pkgs, filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if matched, err := filepath.Match("*.go", filepath.Base(path)); err != nil {
			return err
		} else if matched {

			if src, err := parser.ParseFile(fset, path, nil, parser.ParseComments); err == nil {

				if src.Name.Name == "main" || len(src.Comments) == 0 {
					return nil
				}

				var (
					filename = filepath.Base(path)
					dir      = filepath.Dir(path)
				)

				pkg, found := pkgs[dir]

				if !found {
					pkg = &ast.Package{
						Name:  src.Name.Name,
						Files: make(map[string]*ast.File),
					}
					pkgs[dir] = pkg
				}

				pkg.Files[filename] = src
			} else {
				return err
			}
		}

		return nil
	})
}
