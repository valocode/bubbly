package main

import (
	"github.com/swaggo/swag/example/basic/docs"
	"github.com/verifa/bubbly/cmd"
)

// @title Bubbly
// @version 1.0  // Change version here
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@bubbly.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @termsOfService http://bubbly.io/terms/

// @host localhost:8080
func main() {
	setSwaggerInfo()
	cmd.Execute()
}

func setSwaggerInfo() {
	docs.SwaggerInfo.Title = "Bubbly Api"
	docs.SwaggerInfo.Description = "API schema and information for the bubbly server"
	docs.SwaggerInfo.Version = "1.0"
	// TODO(server): Have host be defined by environment variables
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}
