FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates bash curl jq yq

WORKDIR /root/

COPY --from=builder /app/cmd/main .

COPY configs ./configs
COPY tools ./tools
RUN chmod +x ./tools/load_envs.sh

ENV ENV=prod

CMD ["./tools/load_envs.sh", "./main"]
