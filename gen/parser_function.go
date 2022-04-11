package gen

import (
	"go/ast"
	"reflect"
)

func functionParser(p *param) (function *Function) {

	var decl = p.decl

	switch t := decl.(type) {

	case *ast.FuncDecl:

		var _path, method, ok = hasComment(t.Doc, `(@GET|@POST|@PUT|@PATCH|@DELETE|@Configuration)(\(.*\))?`)

		if !ok {
			return nil
		}

		if !IsExported(decl) {
			return nil
		}

		if t.Recv != nil {
			return nil
		}

		if method == "Configuration" {
			if t.Type.Results == nil || len(t.Type.Results.List) > 1 {
				logger.Fatalln("Configuration can be empty or more 1")
			}

			var (
				serviceName string
				servicePkg  = p.file.Name.Name
			)

			switch r := t.Type.Results.List[0].Type.(type) {
			case *ast.Ident:
				serviceName = r.Name
			case *ast.StarExpr:
				switch x := r.X.(type) {
				case *ast.Ident:
					serviceName = x.Name
				}
				logger.Println(reflect.TypeOf(r.X))
			case *ast.SelectorExpr:
				logger.Println(reflect.TypeOf(r.Sel))
			default:
				logger.Println(reflect.TypeOf(r))
			}

			return &Function{
				dir:         p.dir,
				pkg:         p.file.Name.Name,
				name:        t.Name.Name,
				method:      method,
				serviceName: serviceName,
				servicePkg:  servicePkg,
				Type:        1,
			}
		}

		var arguments = argumentParser(p, t)
		var results = resultParser(p, t)

		function = &Function{
			dir:       p.dir,
			pkg:       p.file.Name.Name,
			name:      t.Name.Name,
			method:    method,
			path:      _path,
			arguments: arguments,
			results:   results,
		}

	}

	return function
}
