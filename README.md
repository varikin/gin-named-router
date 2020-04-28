# Gin Named Router

Adds the ability to name routes in [Gin](https://github.com/gin-gonic/).

See [GoDocs](https://pkg.go.dev/github.com/varikin/gin-named-router?tab=doc) for the documentation.
Hopefully, the example is enough.
If not, please let me know.

## Should I use this?

NO!

This is incomplete.
It does not yet support groups in Gin which is the biggest missing feature gap.

## Example

```
func Example() {
	// Setup the Gin router with named routes
	engine := gin.Default()
	router := namedrouter.New(engine)
	router.Get("root", "/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello")
	})
	router.Get("user", "/user/:id", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello")
	})

	// Start the router (but not in a simple example because it blocks)
	// router.Run(":8080")

	// Elsewhere in a handler
	rootPath, _ := router.Reverse("root").Path()
	println(rootPath)
	// Output: /

	path, _ := router.Reverse("user").With("id", "3").Path()
	println(path)
	// Output: /user/3
}
```