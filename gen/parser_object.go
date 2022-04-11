package gen

import (
	"go/ast"
	"go/token"
)

func ObjectParser(pa *param) (Controllers Objects, Services Objects) {

	Services = Objects{}
	Controllers = Objects{}
	var decl = pa.decl

	switch t := decl.(type) {
	case *ast.GenDecl:

		if t.Tok != token.TYPE {
			return
		}

		var _path, _type, ok = hasComment(t.Doc, `(@Controller|@Service)(\(.*\))?`)

		if !ok {
			return
		}

		if !IsExported(decl) || len(t.Specs) == 0 {
			return
		}

		var spec, has = t.Specs[0].(*ast.TypeSpec)

		if !has {
			return
		}
		//*ast.StructType
		//*ast.InterfaceType
		var objType = GenKind(0)

		switch spec.Type.(type) {

		case *ast.InterfaceType:
			objType = ServiceInterface
		case *ast.StructType:
			if _type == "Service" {
				objType = ServiceStruct
			} else {
				objType = Controller
			}
		}

		if _type == "Service" {
			Services.Add(&Object{
				dir:  pa.dir,
				pkg:  pa.file.Name.Name,
				name: spec.Name.Name,
				kind: objType,
			})
		} else if _type == "Controller" {
			Controllers.Add(&Object{
				dir:  pa.dir,
				url:  _path,
				pkg:  pa.file.Name.Name,
				name: spec.Name.Name,
				kind: Controller,
			})
		}

	}

	return
}
