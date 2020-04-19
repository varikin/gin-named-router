// Package namedrouter adds the ability to name routes in Gin without hiding Gin from the user.
//
// See https://pkg.go.dev/github.com/gin-gonic/gin for more information information on Gin.
package namedrouter

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// NamedRouter wraps a `*gin.Engine` instance while also maintaining a mapping of names to routes.
// Use `namedrouter.New(engine *gin.Engine)` to create a new instance.
type NamedRouter struct {
	*gin.Engine
	Names map[string]string
}

// NamedRoute is an instance of named route with parameters for a specific route.
// See the examples for further information.
type NamedRoute struct {
	Name string
	Route string
	Parameters map[string]string
}

// New returns a new instance of a `*NamedRouter` with the given `*gin.Engine`.
// NamedRouter does not have any requirements on how the `*gin.Engine` is constructed, instead
// preferring to have instance passed in.
func New(engine *gin.Engine) *NamedRouter {
	return &NamedRouter{
		Engine: engine,
		Names: make(map[string]string),
	}
}

// Reverse returns a `NamedRoute` for the given `name`.
// See examples for more information.
func (nr *NamedRouter) Reverse(name string) NamedRoute {
	return NamedRoute{
		Name: name,
		Route: nr.Names[name],
		Parameters: make(map[string]string),
	}

}

// With adds a parameter by name to the `NamedRoute`.
// Returns the same named route so it can be used as a fluent-api.
func (nr NamedRoute) With(key, value string) NamedRoute {
	nr.Parameters[key] = value
	return nr
}

// Path returns the absolute path to the named URL or an error if the path can't be constructed.
// All the named parameters in the URL will be replaced with parameter values.
// The error will returned if route doesn't exist by name, a parameter doesn't exists by name,
// or all the parameters are not used.
// The constructed path does not include the trailing slash nor the domain and scheme portion.
// For example, the route to http://example.org/users will return /users.
func (nr NamedRoute) Path() (string, error) {
	if nr.Route == "" {
		return "", fmt.Errorf("undefined route for %s", nr.Name)
	}
	if len(nr.Parameters) == 0 {
		return nr.Route, nil
	}

	var url strings.Builder
	url.WriteString("/")

	parts := strings.Split(nr.Route, "/")
	lastPartIndex := len(parts) - 1
	for i, part := range parts {
		if part == "" {
			continue
		} else if part[0:1] == ":" || part[0:1] == "*" {
			param := part[1:]
			value, exists := nr.Parameters[param]
			if !exists {
				return "", fmt.Errorf("parameter in named route not set: %s", param)
			}
			url.WriteString(value)

			// Consume the parameter to later ensure everything is used
			delete(nr.Parameters, param)
		} else {
			url.WriteString(part)
		}
		// Only add the last slash between parts
		if i < lastPartIndex {
			url.WriteString("/")
		}
	}
	if len(nr.Parameters) > 0 {
		return "", fmt.Errorf("not all parameters used in the route: %v", nr.Parameters)
	}

	return url.String(), nil
}

// Post wraps gin.Engine.POST(path, handle) and names the wrapped route.
func (nr *NamedRouter) Post(name string, relativePath string, handlers ...gin.HandlerFunc) {
	nr.Names[name] = relativePath
	nr.POST(relativePath, handlers...)
}

// Get wraps gin.Engine.GET(path, handle) and names the wrapped route.
func (nr *NamedRouter) Get(name string, relativePath string, handlers ...gin.HandlerFunc) {
	nr.Names[name] = relativePath
	nr.GET(relativePath, handlers...)
}

// Delete wraps gin.Engine.DELETE(path, handle) and names the wrapped route.
func (nr *NamedRouter) Delete(name string, relativePath string, handlers ...gin.HandlerFunc) {
	nr.Names[name] = relativePath
	nr.DELETE(relativePath, handlers...)
}

// Patch wraps gin.Engine.Patch(path, handle) and names the wrapped route.
func (nr *NamedRouter) Patch(name string, relativePath string, handlers ...gin.HandlerFunc) {
	nr.Names[name] = relativePath
	nr.PATCH(relativePath, handlers...)
}

// Put wraps gin.Engine.Put(path, handle) and names the wrapped route.
func (nr *NamedRouter) Put(name string, relativePath string, handlers ...gin.HandlerFunc) {
	nr.Names[name] = relativePath
	nr.PUT(relativePath, handlers...)
}

// Options wraps gin.Engine.Options(path, handle) and names the wrapped route.
func (nr *NamedRouter) Options(name string, relativePath string, handlers ...gin.HandlerFunc) {
	nr.Names[name] = relativePath
	nr.OPTIONS(relativePath, handlers...)
}

// Head wraps gin.Engine.Head(path, handle) and names the wrapped route.
func (nr *NamedRouter) Head(name string, relativePath string, handlers ...gin.HandlerFunc) {
	nr.Names[name] = relativePath
	nr.HEAD(relativePath, handlers...)
}