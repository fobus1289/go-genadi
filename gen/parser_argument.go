package gen

import (
	"fmt"
	"go/ast"
	"reflect"
)

func argumentParser(p *param, decl *ast.FuncDecl) Arguments {

	var (
		list      = decl.Type.Params.List
		arguments = Arguments{}
	)

	if len(list) == 0 {
		return arguments
	}

	for _, field := range list {
		switch t := field.Type.(type) {
		case *ast.Ident:
			var service = p.services.FindService(fmt.Sprintf("%s.%s", p.pkg, t.Name))

			if service == nil {
				logger.Fatalln("param not supported 1", t.Name)
			}

			arguments.Add(&Argument{
				pkg:             p.pkg,
				key:             field.Names[0].Name,
				value:           t.Name,
				kind:            getKind(t.Name),
				objectValueName: service.GetValueServiceName(),
				hasPtr:          true,
			})
		case *ast.StarExpr:
			switch x := t.X.(type) {
			case *ast.Ident:

				var service = p.services.FindService(fmt.Sprintf("%s.%s", p.pkg, x.Name))

				if service == nil {
					logger.Println(fmt.Sprintf("%s.%s", p.pkg, x.Name))
					logger.Fatalln("param not supported 2", x.Name)
				}

				arguments.Add(&Argument{
					pkg:             p.pkg,
					key:             field.Names[0].Name,
					value:           x.Name,
					objectValueName: service.GetValueServiceName(),
					kind:            getKind(x.Name),
					hasPtr:          true,
				})

			case *ast.SelectorExpr:
				var ident = x.X.(*ast.Ident)
				var service = p.services.FindService(fmt.Sprintf("%s.%s", ident.Name, x.Sel.Name))
				logger.Println(fmt.Sprintf("%s.%s", ident.Name, x.Sel.Name))
				if service == nil {
					logger.Fatalln("param not supported 3", ident.Name, x.Sel.Name)
				}

				arguments.Add(&Argument{
					pkg:             ident.Name,
					key:             field.Names[0].Name,
					value:           x.Sel.Name,
					objectValueName: service.GetValueServiceName(),
					kind:            getKind(x.Sel.Name),
					hasPtr:          true,
				})

			default:
				logger.Println(reflect.TypeOf(x))
			}

		case *ast.SelectorExpr:

			var ident = t.X.(*ast.Ident)

			var service = p.services.FindService(fmt.Sprintf("%s.%s", ident.Name, t.Sel.Name))
			if service == nil {
				for _, object := range p.services {
					logger.Println(object)
				}
				logger.Fatalln("param not supported 4", fmt.Sprintf("%s.%s", ident.Name, t.Sel.Name))
			}

			arguments.Add(&Argument{
				pkg:             ident.Name,
				key:             field.Names[0].Name,
				value:           t.Sel.Name,
				objectValueName: service.GetValueServiceName(),
				kind:            getKind(t.Sel.Name),
				hasPtr:          false,
			})

		}

	}

	return arguments
}

func getKind(name string) Kind {
	switch name {
	case "ResponseWriter":
		return Interface
	}

	return Struct
}
