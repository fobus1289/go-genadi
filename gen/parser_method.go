package gen

import (
	"go/ast"
	"path"
)

func methodParser(p *param, controller *Object) {

	var decl = p.decl

	switch t := decl.(type) {
	case *ast.FuncDecl:

		if t.Recv == nil {
			return
		}

		if t.Name.Name == "AfterCreate" && t.Type.Params != nil && len(t.Type.Params.List) > 0 {
			var _object *ast.Ident

			switch r := t.Recv.List[0].Type.(type) {
			case *ast.StarExpr:
				_object = r.X.(*ast.Ident)
			case *ast.Ident:
				_object = r
			default:
				return
			}

			if _object.Name != controller.name || controller.dir != p.dir || controller.pkg != p.file.Name.Name {
				return
			}

			var arguments = argumentParser(p, t)

			controller.methods.Add(&Method{
				Type:            2,
				name:            t.Name.Name,
				arguments:       arguments,
				objectValueName: controller.GetValueName(),
			})
			return
		}

		var _path, method, ok = hasComment(t.Doc, `(@GET|@POST|@PUT|@PATCH|@DELETE)(\(.*\))?`)

		if !ok {
			return
		}

		if !IsExported(decl) {
			return
		}

		var _object *ast.Ident

		switch r := t.Recv.List[0].Type.(type) {
		case *ast.StarExpr:
			_object = r.X.(*ast.Ident)
		case *ast.Ident:
			_object = r
		default:
			return
		}

		if _object.Name != controller.name || controller.dir != p.dir || controller.pkg != p.file.Name.Name {
			return
		}
		var arguments = argumentParser(p, t)
		var results = resultParser(p, t)

		controller.methods.Add(&Method{
			Type:      1,
			name:      t.Name.Name,
			method:    method,
			path:      path.Join(controller.url, _path),
			arguments: arguments,
			results:   results,
		})

	}

}
