package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {

		return c.JSON(200, map[string]interface{}{
			"pod":     os.Getenv("HOSTNAME"),
			"message": os.Getenv("SYSTEM_NAME"),
		})
	})

	e.GET("/helloworld", func(c echo.Context) error {
		response, err := http.Get("http://rm-service")
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"error": err.Error(),
			})
		}

		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"error": err.Error(),
			})
		}

		return c.JSON(200, map[string]interface{}{
			"message": string(contents),
		})
	})

	e.GET("/secret", func(c echo.Context) error {

		return c.JSON(200, map[string]interface{}{
			"message": os.Getenv("SECRET_KEY"),
		})
	})

	e.GET("/configmap", func(c echo.Context) error {

		return c.JSON(200, map[string]interface{}{
			"message": os.Getenv("ENV"),
		})
	})

	e.GET("/sleep", func(c echo.Context) error {
		if c.QueryParam("isSleep") == "true" {
			time.Sleep(10 * time.Second)

			return c.JSON(200, map[string]interface{}{
				"pod":     os.Getenv("HOSTNAME"),
				"message": "sleep 10 second",
			})
		}

		return c.JSON(200, map[string]interface{}{
			"pod":     os.Getenv("HOSTNAME"),
			"message": "no sleep",
		})
	})

	e.GET("/panic", func(c echo.Context) error {
		client := redis.NewClient(&redis.Options{
			Addr:     "rm-redis:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		_, err := client.Ping().Result()
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"error": err.Error(),
			})
		}

		val, err := client.Get("RM_COUNT").Result()
		if err != nil {
			client.Set("RM_COUNT", os.Getenv("HOSTNAME"), 10*time.Minute).Err()
			log.Fatal("hello")
		}

		return c.JSON(200, map[string]interface{}{
			"message":    fmt.Sprintf("%s this pod is no response", val),
			"currentPod": os.Getenv("HOSTNAME"),
		})
	})

	e.GET("/reset", func(c echo.Context) error {
		client := redis.NewClient(&redis.Options{
			Addr:     "rm-redis:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		_, err := client.Ping().Result()
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"error": err.Error(),
			})
		}

		client.Del("RM_COUNT").Err()

		return c.JSON(200, map[string]interface{}{
			"message": "Reset the redis count to 0",
		})
	})

	e.Logger.Fatal(e.Start(":5000"))
}
