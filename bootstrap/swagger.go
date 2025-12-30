package bootstrap

import (
	"GoHub-Service/docs"
	"GoHub-Service/pkg/config"
)

// SetupSwaggerDocs sets swagger meta before routes are registered.
// This file is not auto-generated and is safe from swag init overwrites.
func SetupSwaggerDocs() {
	docs.SwaggerInfo.Title = config.GetString("app.name", "GoHub Service API")
	docs.SwaggerInfo.Description = "GoHub Service API Docs"
	docs.SwaggerInfo.Version = config.GetString("app.version", "1.0")

	// Host/BasePath help Swagger UI build request URLs; rely on APP_URL/api prefix when api_domain unset.
	docs.SwaggerInfo.Host = config.GetString("app.api_domain", "localhost:3000")
	if len(docs.SwaggerInfo.Host) == 0 {
		docs.SwaggerInfo.Host = "localhost:3000"
	}
	docs.SwaggerInfo.BasePath = "/api/v1"

	// Support both http/https depending on deployment
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}
