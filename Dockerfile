FROM golang:1.23.2-alpine

RUN apk add --no-cache git curl

RUN go install github.com/air-verse/air@latest

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app

COPY . .

RUN go mod vendor

CMD ["air"]
