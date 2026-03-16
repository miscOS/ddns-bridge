package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/miscOS/ddns-bridge/auth"
	db "github.com/miscOS/ddns-bridge/database"
	"github.com/miscOS/ddns-bridge/routes"
	"github.com/miscOS/ddns-bridge/services"
)

func main() {

	// Load the environment variables
	godotenv.Load()

	// Initialize the database
	db.Init()

	// Register services
	services.Register(services.NewExampleDNS())
	services.Register(services.NewCloudflareDNS())

	// Create a Gin router with default middleware (logger and recovery)
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// Enable CORS middleware
	router.Use(cors.Default())

	// Register authenticators
	// Local JWT authentication
	if os.Getenv("DISABLE_LOCAL_AUTH") == "true" {
		log.Println("Local JWT authentication is disabled")
	} else {
		log.Println("Registering local JWT authenticator")

		authenticator, err := auth.NewLocalAuthenticator(os.Getenv("LOCAL_AUTH_JWT_SECRET"), os.Getenv("LOCAL_AUTH_REGISTER_KEY"))
		if err != nil {
			log.Fatalf("Failed to create local authenticator: %s", err)
		}
		auth.Register(authenticator)
	}

	// OIDC authentication
	authenticator, err := auth.NewOIDCAuthenticator(os.Getenv("OIDC_ISSUER"), os.Getenv("OIDC_CLIENT_ID"))
	if err != nil {
		log.Fatalf("Failed to create OIDC authenticator: %s", err)
	}
	auth.Register(authenticator)

	// Add routes
	routes.RegisterRoutes(router)
	routes.RegisterSecureRoutes(router)

	// Start the server on port 8080 or environment variable PORT
	router.Run()
}
