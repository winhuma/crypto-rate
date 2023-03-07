# crypto-rate

## Description
Project get rate crypto and serv API

## How to Run

Clone this repo and move in to directory.
On file `docker-compose.yml` must set environment variable about binance api.
You can go to `https://www.binance.com` for get this data.

```
BINANCE_API_KEY: <binance_api_key>
BINANCE_SECRET: <binance_secert>
```

Run follow command for build docker images.


```
chmod +x build-service.sh
./build-service.sh
docker-compose up
```

Now system start. Can follow response from `http://127.0.0.1:8080/crypto/get`
Or go to url `http://127.0.0.1:8880` for display with adminer