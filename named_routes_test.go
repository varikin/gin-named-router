package namedrouter_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/varikin/gin-named-router"
)

func noop(c *gin.Context) {}

func TestNamedRoute_Path(t *testing.T) {
	type route struct {
		name    string
		route   string
		handler gin.HandlerFunc
	}

	routes := []route{
		{"root", "/", noop},
		{"index", "/index", noop},
		{"about", "/about/us", noop},
		{"user", "/user/:id", noop},
		{"user-item", "/user/:id/item/:item", noop},
		{"item-splat", "/item/*splat", noop},
	}

	tests := []struct {
		test      string
		want      string
		wantErr   error
		routeName string
		params    map[string]string
		routes    []route
	}{
		{
			test:      "Root path",
			want:      "/",
			routeName: "root",
			routes:    routes,
		}, {
			test:      "Simple path",
			want:      "/about/us",
			routeName: "about",
			routes:    routes,
		}, {
			test:      "An argument at the end",
			want:      "/user/3",
			routeName: "user",
			params:    map[string]string{"id": "3"},
			routes:    routes,
		}, {
			test:      "Multiple arguments",
			want:      "/user/3/item/book",
			routeName: "user-item",
			params:    map[string]string{"id": "3", "item": "book"},
			routes:    routes,
		}, {
			test:      "A star arguments",
			want:      "/item/records/7",
			routeName: "item-splat",
			params:    map[string]string{"splat": "records/7"},
			routes:    routes,
		}, {
			test:      "No route defined",
			wantErr:   namedrouter.NoRouteDefinedError("unknown"),
			routeName: "unknown",
			routes:    routes,
		}, {
			test:      "Unknown parameter",
			wantErr:   namedrouter.RouteParameterNotSet("id"),
			routeName: "user",
			routes:    routes,
		}, {
			test:      "Unused parameter",
			wantErr:   namedrouter.UnknownRouteParameter("item"),
			routeName: "user",
			params:    map[string]string{"id": "3", "item": "3"},
			routes:    routes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			router := namedrouter.New(gin.Default())
			for _, route := range tt.routes {
				router.Get(route.name, route.route, route.handler)
			}
			route := router.Reverse(tt.routeName)
			for key, value := range tt.params {
				route.With(key, value)
			}
			got, err := route.Path()

			if err != tt.wantErr {
				t.Errorf("Path() unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("Path() got = %v, want %v", got, tt.want)
			}
		})
	}
}
