package gen

import "fmt"

type Argument struct {
	pkg             string
	key             string
	value           string
	objectValueName string
	kind            Kind
	hasPtr          bool
}

func NewArgument(pkg, key, value string, kind Kind, hasPtr bool) *Argument {
	return &Argument{
		pkg:    pkg,
		key:    key,
		value:  value,
		kind:   kind,
		hasPtr: hasPtr,
	}
}

type Arguments []*Argument

func (args *Arguments) Add(arg ...*Argument) {
	*args = append(*args, arg...)
}

func (arg *Argument) IsPtr() bool {
	return arg.hasPtr
}

func (arg *Argument) HasPkg() bool {
	return arg.pkg != ""
}

func (arg *Argument) GetPackage() string {
	return arg.pkg
}

func (arg *Argument) GetKey() string {
	return arg.key
}

func (arg *Argument) GetValue() string {
	return arg.value
}

func (arg *Argument) GetType() Kind {
	return arg.kind
}

func (arg *Argument) GetParam() string {

	var (
		pkg    = arg.pkg
		value  = arg.value
		result string
	)

	if arg.HasPkg() {
		result = fmt.Sprintf("%s.", pkg)
	}

	result = fmt.Sprintf("%s%s", result, value)

	if arg.IsPtr() {
		result = fmt.Sprintf("&%s", result)
	}

	return result
}
