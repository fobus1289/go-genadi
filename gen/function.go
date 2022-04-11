package gen

import (
	"fmt"
	"path"
	"strings"
)

type Functions []*Function

func (fns *Functions) Add(functions ...*Function) {
	*fns = append(*fns, functions...)
}

func (fns Functions) GetMethod(key string) []IFS {
	var ifs []IFS

	for _, fn := range fns {
		if fn.method == key {
			ifs = append(ifs, fn)
		}
	}

	return ifs
}

func (fns Functions) GetMethods() []IFS {
	var ifs []IFS

	for _, fn := range fns {
		ifs = append(ifs, fn)
	}

	return ifs
}

func (f *Function) GetPath() string {
	return f.path
}

func (Functions) GetObjects() []IObject {
	return []IObject{}
}

func (f *Function) GetConfiguration() string {
	if f.object == nil {
		return ""
	}
	return fmt.Sprintf("%s =%s()\n", f.object.GetValueName(), f.CallMethod())
}

func (f *Function) GetConfigurationName() string {
	var dir = strings.Replace(path.Clean(f.dir), "\\", "_", -1)
	return fmt.Sprintf("configuration_%s_%s_%s_%s", dir, f.pkg, f.name, f.method)
}

func (f *Function) GetServiceParams() string {
	return ""
}

func (f *Function) GetIf() string {
	if f.Type != 0 {
		return ""
	}

	var patternName = fmt.Sprintf("%s.MatchString(request.URL.Path)", f.GetPatternName())
	var _if = fmt.Sprintf("\nif %s {", patternName)

	var arguments []string

	for _, argument := range f.arguments {
		var start string
		if !argument.hasPtr && argument.kind != Interface {
			start = "*"
		}
		arguments = append(arguments, fmt.Sprintf("%s%s", start, argument.objectValueName))
	}

	var (
		result string
		write  string
	)

	if len(f.results) > 0 {
		result = "var result = "
		write = `
			if out, err := json.Marshal(result); err == nil {
				_, _ = responseWriter.Write(out)
			}`
	}

	var callMethod = fmt.Sprintf("%s%s(%s)%s", result, f.CallMethod(), strings.Join(arguments, ","), write)

	//var callMethod = fmt.Sprintf("%s(%s)", f.CallMethod(), strings.Join(arguments, ","))
	var ret = fmt.Sprintf("%s\n%s\n}", _if, callMethod)

	return ret
}

func (f *Function) CallMethod() string {
	return fmt.Sprintf("%s.%s", f.pkg, f.name)
}

func (f *Function) GetPatternName() string {
	if f.Type != 0 {
		return ""
	}
	var dir = strings.Replace(path.Clean(f.dir), "\\", "_", -1)
	return fmt.Sprintf("pattern_function_%s_%s_%s_%s", dir, f.pkg, f.name, f.method)
}

type Function struct {
	dir         string
	pkg         string
	name        string
	method      string
	serviceName string
	servicePkg  string
	object      *Object
	Type        int
	path        string
	arguments   Arguments
	results     Results
}

func NewFunction(dir, pkg, name string, arguments Arguments, results Results) *Function {
	return &Function{
		dir:       dir,
		pkg:       pkg,
		name:      name,
		arguments: arguments,
		results:   results,
	}
}
