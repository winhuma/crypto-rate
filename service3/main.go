package main

import (
	"crypto-rate/libs/myconnect"
	"crypto-rate/libs/myfunc"
	"crypto-rate/libs/mymodels"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func HelloService(c *fiber.Ctx) error {
	return c.Status(200).SendString("hello")
}

func GetPriceRateAll(c *fiber.Ctx) error {
	token := c.Query("token")
	myredis := myconnect.RedisGetInstance()
	mydata, err := myredis.Get("rate").Result()
	if err != nil && err != redis.Nil {
		return myfunc.MyErrFormat(err)
	}
	var result = map[string]interface{}{}
	var myres interface{}
	if mydata == "" {
		var getRate mymodels.DBCryptoRate
		mydb := myconnect.DBInstance()
		err = mydb.Table("crypto_rate").Last(&getRate).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return myfunc.MyErrFormat(err)
		}
		if !getRate.Data.Valid {
			return c.Status(200).JSON(map[string]interface{}{
				"message": "success",
				"data":    nil,
			})
		}
		mydata = getRate.Data.String
	}

	err = json.Unmarshal([]byte(mydata), &result)
	if err != nil {
		return myfunc.MyErrFormat(err)
	}

	myres = result
	if token != "" {
		myres = result[token]
	}
	return c.Status(200).JSON(map[string]interface{}{
		"message": "success",
		"data":    myres,
	})
}

func SetRoutes(app *fiber.App) {
	app.Get("/", HelloService)
	app.Get("/crypto/get", GetPriceRateAll)
}

func NewApp() *fiber.App {

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	return app
}

func RunServe(app *fiber.App, port string) {
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Info().Msg(".env file not found")
	}

	myconnect.RedisConnect(os.Getenv("REDIS_URL"))
	myconnect.NewDb(os.Getenv("DB_URL"))

	app := NewApp()
	SetRoutes(app)
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprint(8080)
	}
	RunServe(app, port)
}
