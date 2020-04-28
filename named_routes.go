// Package namedrouter adds the ability to name routes in Gin.
//
// A namedrouter wraps a gin.Engine instance while still exposing the wrapped gin.Engine instance.
//
// See also
//
// https://pkg.go.dev/github.com/gin-gonic/gin
package namedrouter

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// NamedRouter wraps a gin.Engine instance while also maintaining a mapping of names to routes.
// Use namedrouter.New(engine *gin.Engine) to create a new instance.
type NamedRouter struct {
	*gin.Engine
	names map[string]string
}

// NamedRoute is an instance of named route with parameters for a specific route.
type NamedRoute struct {
	name       string
	route      string
	parameters map[string]string
}

type NamedGroup struct {
	*gin.RouterGroup
	namedRouter *NamedRouter
	basePath string
}

// New returns a new instance which wraps the given gin.Engine.
// NamedRouter does not have any requirements on how the gin.Engine is constructed, instead
// preferring to have instance passed in.
func New(engine *gin.Engine) *NamedRouter {
	return &NamedRouter{
		Engine: engine,
		names:  make(map[string]string),
	}
}

// NamedGroup returns a new NamedGroup for a NamedRouter.
// A NamedGroups wrap gin.RouterGroup.
func (nr *NamedRouter) NamedGroup(path string, handlers ...gin.HandlerFunc) *NamedGroup {
	group := nr.Group(path, handlers...)
	return &NamedGroup{
		RouterGroup: group,
		namedRouter: nr,
		basePath: path,
	}
}

// NamedGroup returns a nested NamedGroup of another NamedGroup.
func (ng *NamedGroup) NamedGroup(path string, handlers ...gin.HandlerFunc) *NamedGroup {
	group := ng.RouterGroup.Group(path, handlers...)
	basePath := ng.basePath + "/" + path
	return &NamedGroup{
		RouterGroup: group,
		namedRouter: ng.namedRouter,
		basePath:    basePath,
	}
}

// Reverse returns a NamedRoute for the given name.
func (nr *NamedRouter) Reverse(name string) NamedRoute {
	return NamedRoute{
		name:       name,
		route:      nr.names[name],
		parameters: make(map[string]string),
	}

}

// With adds a parameter by name to the NamedRoute.
// Returns the same named route so it can be used as a fluent-api.
func (nr NamedRoute) With(key, value string) NamedRoute {
	nr.parameters[key] = value
	return nr
}

// Path returns the absolute path to the named URL or an error if the path can't be constructed.
// All the named parameters in the URL will be replaced with parameter values.
// The error will returned if route doesn't exist by name, a parameter doesn't exists by name,
// or all the parameters are not used.
// The constructed path does not include the trailing slash nor the domain and scheme portion.
// For example, the route to http://example.org/users will return /users.
func (nr NamedRoute) Path() (string, error) {
	if nr.route == "" {
		return "", NoRouteDefinedError(nr.name)
	}

	var url strings.Builder
	url.WriteString("/")

	parts := strings.Split(nr.route, "/")
	lastPartIndex := len(parts) - 1
	for i, part := range parts {
		if part == "" {
			continue
		} else if part[0:1] == ":" || part[0:1] == "*" {
			param := part[1:]
			value, exists := nr.parameters[param]
			if !exists {
				return "", RouteParameterNotSet(param)
			}
			url.WriteString(value)

			// Consume the parameter to later ensure everything is used
			delete(nr.parameters, param)
		} else {
			url.WriteString(part)
		}
		// Only add the last slash between parts
		if i < lastPartIndex {
			url.WriteString("/")
		}
	}
	if len(nr.parameters) > 0 {
		for key := range nr.parameters {
			return "", UnknownRouteParameter(key)
		}
	}

	return url.String(), nil
}

// Post wraps gin.Engine.POST(path, handle) and names the wrapped route.
func (nr *NamedRouter) Post(name string, relativePath string, handlers ...gin.HandlerFunc) {
	nr.names[name] = relativePath
	nr.POST(relativePath, handlers...)
}

// Get wraps gin.Engine.GET(path, handle) and names the wrapped route.
func (nr *NamedRouter) Get(name string, relativePath string, handlers ...gin.HandlerFunc) {
	nr.names[name] = relativePath
	nr.GET(relativePath, handlers...)
}

// Delete wraps gin.Engine.DELETE(path, handle) and names the wrapped route.
func (nr *NamedRouter) Delete(name string, relativePath string, handlers ...gin.HandlerFunc) {
	nr.names[name] = relativePath
	nr.DELETE(relativePath, handlers...)
}

// Patch wraps gin.Engine.Patch(path, handle) and names the wrapped route.
func (nr *NamedRouter) Patch(name string, relativePath string, handlers ...gin.HandlerFunc) {
	nr.names[name] = relativePath
	nr.PATCH(relativePath, handlers...)
}

// Put wraps gin.Engine.Put(path, handle) and names the wrapped route.
func (nr *NamedRouter) Put(name string, relativePath string, handlers ...gin.HandlerFunc) {
	nr.names[name] = relativePath
	nr.PUT(relativePath, handlers...)
}

// Options wraps gin.Engine.Options(path, handle) and names the wrapped route.
func (nr *NamedRouter) Options(name string, relativePath string, handlers ...gin.HandlerFunc) {
	nr.names[name] = relativePath
	nr.OPTIONS(relativePath, handlers...)
}

// Head wraps gin.Engine.Head(path, handle) and names the wrapped route.
func (nr *NamedRouter) Head(name string, relativePath string, handlers ...gin.HandlerFunc) {
	nr.names[name] = relativePath
	nr.HEAD(relativePath, handlers...)
}

// Post wraps gin.Engine.POST(path, handle) and names the wrapped route.
func (ng *NamedGroup) Post(name string, relativePath string, handlers ...gin.HandlerFunc) {
	ng.namedRouter.names[name] = ng.basePath + relativePath
	ng.POST(relativePath, handlers...)
}

// Get wraps gin.Engine.GET(path, handle) and names the wrapped route.
func (ng *NamedGroup) Get(name string, relativePath string, handlers ...gin.HandlerFunc) {
	ng.namedRouter.names[name] = ng.basePath + relativePath
	ng.GET(relativePath, handlers...)
}

// Delete wraps gin.Engine.DELETE(path, handle) and names the wrapped route.
func (ng *NamedGroup) Delete(name string, relativePath string, handlers ...gin.HandlerFunc) {
	ng.namedRouter.names[name] = ng.basePath + relativePath
	ng.DELETE(relativePath, handlers...)
}

// Patch wraps gin.Engine.Patch(path, handle) and names the wrapped route.
func (ng *NamedGroup) Patch(name string, relativePath string, handlers ...gin.HandlerFunc) {
	ng.namedRouter.names[name] = ng.basePath + relativePath
	ng.PATCH(relativePath, handlers...)
}

// Put wraps gin.Engine.Put(path, handle) and names the wrapped route.
func (ng *NamedGroup) Put(name string, relativePath string, handlers ...gin.HandlerFunc) {
	ng.namedRouter.names[name] = ng.basePath + relativePath
	ng.PUT(relativePath, handlers...)
}

// Options wraps gin.Engine.Options(path, handle) and names the wrapped route.
func (ng *NamedGroup) Options(name string, relativePath string, handlers ...gin.HandlerFunc) {
	ng.namedRouter.names[name] = ng.basePath + relativePath
	ng.OPTIONS(relativePath, handlers...)
}

// Head wraps gin.Engine.Head(path, handle) and names the wrapped route.
func (ng *NamedGroup) Head(name string, relativePath string, handlers ...gin.HandlerFunc) {
	ng.namedRouter.names[name] = ng.basePath + relativePath
	ng.HEAD(relativePath, handlers...)
}

type NoRouteDefinedError string

func (s NoRouteDefinedError) Error() string {
	return fmt.Sprintf("no route defined for %s", string(s))
}

type UnknownRouteParameter string

func (s UnknownRouteParameter) Error() string {
	return fmt.Sprintf("unknown route parameter: %s", string(s))
}

type RouteParameterNotSet string

func (s RouteParameterNotSet) Error() string {
	return fmt.Sprintf("route parameter not set: %s", string(s))
}