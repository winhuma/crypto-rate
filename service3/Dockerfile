FROM golang:1.18-alpine

WORKDIR /app

COPY . .
COPY libs .
RUN go mod init crypto-rate
RUN go mod tidy

EXPOSE 8080
RUN go build -o myservice
CMD [ "./myservice" ]