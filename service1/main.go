package main

import (
	"context"
	"crypto-rate/libs/myconnect"
	"crypto-rate/libs/myfunc"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var BINANCE_CLIENT *binance.Client

func BinanceConnect() {
	apiKey := viper.GetString("BINANCE_API_KEY")
	secret := viper.GetString("BINANCE_SECRET")
	BINANCE_CLIENT = binance.NewClient(apiKey, secret)
}

func GetAllRate() (map[string]string, error) {
	var ListCurrency = map[string]string{}
	prices, err := BINANCE_CLIENT.NewListPricesService().Do(context.TODO())
	if err != nil {
		return ListCurrency, err
	}
	for _, price := range prices {
		ListCurrency[price.Symbol] = price.Price
	}
	return ListCurrency, nil
}

func MainProcess() error {
	allCurrency, err := GetAllRate()
	if err != nil {
		return myfunc.MyErrFormat(err)
	}
	jsonString, err := json.Marshal(allCurrency)
	if err != nil {
		return myfunc.MyErrFormat(err)
	}
	err = myconnect.KafkaProducer(os.Getenv("KAFKA_URL"), "crypto-rate", "rate", string(jsonString))
	if err != nil {
		return err
	}

	return nil
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Info().Msg(".env file not found")
	}

	log.Info().Msg("KAFKA_URL " + os.Getenv("KAFKA_URL"))

	BinanceConnect()

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Info().Msg("Exiting...")
		done <- true
	}()

	timeCount := time.Now()
	for {
		select {
		case <-done:
			log.Info().Msg("Bye")
			return
		default:
			if time.Since(timeCount).Seconds() >= 10 {
				timeCount = time.Now()
				log.Info().Msg("Working")
				err := MainProcess()
				if err != nil {
					log.Error().Msg(err.Error())
				}
			}
		}
	}
}
