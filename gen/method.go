package gen

import (
	"fmt"
	"strings"
)

type Method struct {
	name            string
	path            string
	method          string
	Type            int
	objectValueName string
	patternName     string
	arguments       Arguments
	results         Results
}

func (m *Method) GetServiceParams() string {
	if m.Type != 2 {
		return ""
	}
	var arguments []string

	for _, argument := range m.arguments {
		var start string
		if !argument.hasPtr && argument.kind != Interface {
			start = "*"
		}
		arguments = append(arguments, fmt.Sprintf("%s%s", start, argument.objectValueName))
	}

	var callMethod = fmt.Sprintf("%s(%s)\n", m.CallMethod(), strings.Join(arguments, ","))

	return callMethod
}

func (m *Method) GetIf() string {
	if m.Type != 1 {
		return ""
	}
	var patternName = fmt.Sprintf("%s.MatchString(request.URL.Path)", m.GetPatternName())
	var _if = fmt.Sprintf("\nif %s {", patternName)

	//var paramBuffer = strings.Builder{}
	var arguments []string
	for _, argument := range m.arguments {
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

	if len(m.results) > 0 {
		result = "var result = "
		write = `
			if out, err := json.Marshal(result); err == nil {
				_, _ = responseWriter.Write(out)
			}`
	}

	var callMethod = fmt.Sprintf("%s%s(%s)%s", result, m.CallMethod(), strings.Join(arguments, ","), write)

	var ret = fmt.Sprintf("%s\n%s\n}", _if, callMethod)

	return ret

}

func (m *Method) GetPath() string {
	return m.path
}

func (m *Method) GetPatternName() string {
	return m.patternName
}

func (m *Method) CallMethod() string {
	return fmt.Sprintf("%s.%s", m.objectValueName, m.name)
}

type Methods []*Method

func (ms *Methods) Add(methods ...*Method) {
	*ms = append(*ms, methods...)
}

func NewMethod(name string, arguments Arguments, results Results) *Method {
	return &Method{
		name:      name,
		arguments: arguments,
		results:   results,
	}
}
