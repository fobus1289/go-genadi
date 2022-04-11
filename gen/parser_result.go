package gen

import (
	"go/ast"
	"log"
)

func resultParser(p *param, decl *ast.FuncDecl) Results {

	if decl.Type.Results == nil {
		return Results{}
	}

	var (
		list    = decl.Type.Results.List
		results = Results{}
	)

	if len(list) == 0 {
		return results
	}

	if len(list) > 1 {
		log.Fatalln("function result cannot be more than one")
	}

	results = append(results, &Result{})

	return results
}
