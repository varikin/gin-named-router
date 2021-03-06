package namedrouter_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/varikin/gin-named-router"
)

func Example() {

	helloFunc := func(c *gin.Context) {
		c.String(http.StatusOK, "Hello")
	}
	// Setup the Gin router with named routes
	engine := gin.Default()
	router := namedrouter.New(engine)
	router.Get("root", "/", helloFunc)
	router.Get("user", "/user/:id", helloFunc)

	// Routes groups be named as well.
	// The group isn't named, but a NamedGroup is needed to register the named routs
	// within that group.
	api := router.NamedGroup("/api")
	api.Get("api-info", "/info", helloFunc)

	// Nested named groups are also supported
	v1 := api.NamedGroup("v1")
	v1.Post("v1-submit", "/submit", helloFunc)


	// Start the router (but not in a simple example because it blocks)
	// router.Run(":8080")

	// Elsewhere in a handler
	rootPath, _ := router.Reverse("root").Path()
	fmt.Println(rootPath)

	path, _ := router.Reverse("user").With("id", "3").Path()
	fmt.Println(path)

	// Only the named router is needed to reverse any routes within groups
	path, _ = router.Reverse("api-info").Path()
	fmt.Println(path)

	path, _ = router.Reverse("v1-submit").Path()
	fmt.Println(path)



	// Output:
	// /
	// /user/3
	// /api/info
	// /api/v1/submit
}

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
