package main

import (
	"context"
	"github.com/gin-gonic/gin"
	gofeatureflag "github.com/open-feature/go-sdk-contrib/providers/go-feature-flag/pkg"
	"github.com/open-feature/go-sdk/pkg/openfeature"
	"google.golang.org/appengine/log"
	"net/http"
	"time"
)

const defaultMessage = "Hello!"
const newWelcomeMessage = "Hello, welcome to this OpenFeature-enabled website!"

func main() {
	// Initialize Go Gin
	engine := gin.Default()

	// Setup a simple endpoint
	engine.GET("/hello", func(c *gin.Context) {

		options := gofeatureflag.ProviderOptions{
			Endpoint: "http://localhost:1031",
			HTTPClient: &http.Client{
				Timeout: 1 * time.Second,
			},
		}

		provider, _ := gofeatureflag.NewProvider(options)

		openfeature.SetProvider(provider)
		client := openfeature.NewClient("demo")

		evaluationCtx := openfeature.NewEvaluationContext(
			"test",
			map[string]interface{}{
				"firstname": "mark",
				"lastname":  "",
				"email":     "mark8s@gmail.com",
				"admin":     true,
				"anonymous": false,
			})

		welcomeMessage, err := client.BooleanValue(context.Background(), "test-flag", false, evaluationCtx)
		if err != nil {
			log.Errorf(context.Background(), err.Error())
			return
		}

		if welcomeMessage {
			c.JSON(http.StatusOK, newWelcomeMessage)
			return
		} else {
			c.JSON(http.StatusOK, defaultMessage)
			return
		}
	})

	engine.Run()
}
