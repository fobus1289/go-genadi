package gen

import (
	"go/ast"
	"go/token"
	"sync"
)

type param struct {
	dir      string
	pkg      string
	filename string
	services Objects
	decl     ast.Decl
	file     *ast.File
}

type params []*param

type Parser struct {
	*sync.Once
	Packages    map[string]*ast.Package
	FileSet     *token.FileSet
	Params      params
	Controllers Objects
	Services    Objects
	Functions   Functions
}

func (p *Parser) Object() {

	for _, prm := range p.Params {
		var controllers, services = ObjectParser(prm)
		p.Controllers.AddRange(controllers)
		p.Services.AddRange(services)
	}

	p.Services.AddRange(Objects{
		&Object{
			pkg:  "http",
			name: "ResponseWriter",
			kind: ServiceInterface,
		},
		&Object{
			pkg:  "http",
			name: "Request",
			kind: ServiceStruct,
		},
	})

	for _, p2 := range p.Params {
		p2.services = p.Services
	}

}

func (p *Parser) Method() {
	for _, service := range p.Services {
		for _, prm := range p.Params {
			methodParser(prm, service)
		}
	}
	for _, controller := range p.Controllers {
		for _, prm := range p.Params {
			methodParser(prm, controller)
		}
	}
}

func (p *Parser) Function() {
	for _, prm := range p.Params {
		var function = functionParser(prm)
		if function != nil {
			if function.method == "Configuration" {
				object := &Object{
					dir:  function.dir,
					pkg:  function.pkg,
					name: function.serviceName,
					kind: Configure,
				}
				function.object = object
				p.Services.Add(object)
			}
			p.Functions.Add(function)
		}
	}
	for _, p2 := range p.Params {
		p2.services = p.Services
	}

}

func (p *Parser) Parse() error {

	if err := p.findDir(); err != nil {
		return err
	}

	p.Object()
	p.Function()
	p.Method()

	return nil
}

func (p *Parser) findDir() (err error) {
	p.Do(func() {
		p.FileSet = token.NewFileSet()
		p.Packages, err = parseGoFiles(p.FileSet)

		if err != nil {
			return
		}

		for dir, _package := range p.Packages {
			for filename, file := range _package.Files {
				for _, decl := range file.Decls {
					var pkg = file.Name.Name
					var parm = &param{
						pkg:      pkg,
						dir:      dir,
						decl:     decl,
						filename: filename,
						file:     file,
					}
					p.Params = append(p.Params, parm)
				}
			}
		}

	})
	return err
}

func (p *Parser) Each(fn func(*param)) (err error) {

	if err := p.findDir(); err != nil {
		return err
	}

	for _, p2 := range p.Params {
		fn(p2)
	}

	return nil
}
