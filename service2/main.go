package main

import (
	"context"
	"crypto-rate/libs/myconnect"
	"crypto-rate/libs/myfunc"
	"crypto-rate/libs/mymodels"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"gopkg.in/guregu/null.v4"
)

func SetUpTableCrypto() {
	mydb := myconnect.DBInstance()
	err := mydb.AutoMigrate(&mymodels.DBCryptoRate{})
	if err != nil {
		log.Fatal().Msg(myfunc.MyErrFormat(err).Error())
	}
	log.Info().Msg("[PROCESS] SetUpTableCrypto success")
}

func InsertHistory(data string) error {
	var insertData = mymodels.DBCryptoRate{
		DateCreated: null.TimeFrom(time.Now()),
		Data:        null.StringFrom(data),
	}

	mydb := myconnect.DBInstance()
	err := mydb.Table("crypto_rate").Create(&insertData).Error
	if err != nil {
		return myfunc.MyErrFormat(err)
	}
	return nil
}

func MainProcess(reader *kafka.Reader, myredis *redis.Client) {

	m, err := reader.ReadMessage(context.TODO())
	if err != nil {
		log.Error().Msg(myfunc.MyErrFormat(err).Error())
	}

	if string(m.Key) == "rate" {
		mydata := string(m.Value)
		err = myredis.Set("rate", mydata, time.Second*60).Err()
		if err != nil {
			log.Error().Msg(myfunc.MyErrFormat(err).Error())
		}
		err = InsertHistory(mydata)
		if err != nil {
			log.Error().Msg(myfunc.MyErrFormat(err).Error())
		}
	}
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Info().Msg(".env file not found")
	}

	configDB := os.Getenv("DB_URL")
	configRedis := os.Getenv("REDIS_URL")
	configKafka := os.Getenv("KAFKA_URL")
	log.Info().Msg("DB_URL " + configDB)
	log.Info().Msg("REDIS_URL " + configRedis)
	log.Info().Msg("KAFKA_URL " + configKafka)

	myconnect.NewDb(configDB)
	SetUpTableCrypto()
	myredis := myconnect.RedisConnect(configRedis)
	reader := myconnect.KafkaReader(configKafka, "crypto-rate")
	defer reader.Close()
	defer myredis.Close()

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Info().Msg("Exiting...")
		done <- true
	}()

	timeCount := time.Now()
	log.Info().Msg("...Start...")
	MainProcess(reader, myredis)
	for {
		select {
		case <-done:
			log.Info().Msg("Bye")
			return
		default:
			if time.Since(timeCount).Seconds() >= 10 {
				timeCount = time.Now()
				log.Info().Msg("Working...")
				MainProcess(reader, myredis)
			}
		}
	}
}
