package gen

import (
	"fmt"
	"path"
	"strings"
)

type Objects []*Object

func (ots *Objects) AddRange(oo Objects) {
	*ots = append(*ots, oo...)
}

func (ots Objects) FindService(name string) *Object {

	for _, ot := range ots {
		logger.Println(ot.GetValueName())
		if name == fmt.Sprintf("%s.%s", ot.pkg, ot.name) {
			return ot
		}
	}

	return nil
}

func (ots *Objects) Add(o ...*Object) {
	*ots = append(*ots, o...)
}

func (ots Objects) GetMethods() []IFS {
	var ifs []IFS

	for _, o := range ots {
		for _, method := range o.methods {
			var dir = strings.Replace(path.Clean(o.dir), "\\", "_", -1)
			method.patternName = fmt.Sprintf("pattern_method_%s_%s_%s_%s_%s", dir, o.pkg, strings.Title(o.name), method.name, method.method)
			method.objectValueName = o.GetValueName()
			ifs = append(ifs, method)
		}
	}

	return ifs
}

func (ots Objects) GetMethod(key string) []IFS {
	var ifs []IFS

	for _, o := range ots {
		for _, fn := range o.methods {
			if fn.method == key {
				var dir = strings.Replace(path.Clean(o.dir), "\\", "_", -1)
				fn.patternName = fmt.Sprintf("pattern_method_%s_%s_%s_%s_%s", dir, o.pkg, strings.Title(o.name), fn.name, fn.method)
				fn.objectValueName = o.GetValueName()
				ifs = append(ifs, fn)
			}
		}
	}

	return ifs
}

func (ots Objects) GetObjects() []IObject {
	var objects []IObject
	for _, ot := range ots {
		objects = append(objects, ot)
	}
	return objects
}

type Object struct {
	dir     string
	pkg     string
	name    string
	url     string
	kind    GenKind
	methods Methods
}

func (o *Object) Dir() string {
	return o.dir
}

func (o *Object) Pkg() string {
	return o.pkg
}

func (o *Object) Name() string {
	return o.name
}

func (o *Object) Url() string {
	return o.url
}

func (o *Object) Methods() Methods {
	return o.methods
}

func (o *Object) GetType() GenKind {
	return o.kind
}

func (o *Object) GetValueServiceName() string {
	//responseWriter
	//request
	switch o.name {
	case "Request":
		return "request"
	case "ResponseWriter":
		return "responseWriter"
	}
	var dir = path.Clean(strings.Replace(o.dir, "\\", "_", -1))
	return fmt.Sprintf("%s_%s_%s", dir, o.pkg, o.name)
}

func (o *Object) GetValueName() string {
	var dir = path.Clean(strings.Replace(o.dir, "\\", "_", -1))
	return fmt.Sprintf("%s_%s_%s", dir, o.pkg, o.name)
}

func NewObject(dir, pkg, name string, kind GenKind, methods Methods) *Object {
	return &Object{
		dir:     dir,
		pkg:     pkg,
		name:    name,
		kind:    kind,
		methods: methods,
	}
}
