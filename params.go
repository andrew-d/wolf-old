package wolf

import (
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

var paramsKey private

func setParamsInContext(ctx context.Context, p httprouter.Params) context.Context {
	// Allocate a map large enough to handle the params
	mm := make(map[string][]string, len(p))

	// Set each one
	var l []string
	var ok bool
	for _, param := range p {
		l, ok = mm[param.Key]
		if !ok {
			l = make([]string, 0, 5) // Default size should generally be "large enough"
		}

		l = append(l, param.Value)
		mm[param.Key] = l
	}

	// Set in context
	return context.WithValue(ctx, &paramsKey, mm)
}

// ParamFrom retrieves the first parameter with the given name from this
// context, along with a boolean indicating whether or not the parameter was
// given.
func ParamFrom(ctx context.Context, name string) (string, bool) {
	if l, ok := AllParamsFrom(ctx, name); ok {
		return l[0], true
	}
	return "", false
}

// AllParamsFrom returns a slice of all parameters with the given name from
// this context, along with a boolean indicating whether or not the parameter
// was given.
func AllParamsFrom(ctx context.Context, name string) ([]string, bool) {
	mm := ctx.Value(&paramsKey).(map[string][]string)
	if l, ok := mm[name]; ok {
		return l, true
	}

	return nil, false
}
