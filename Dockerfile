FROM golang:1.16.2 as builder
WORKDIR /go/src/github.com/pcherednichenko/spam_fighter_bot
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -o spam_fighter_bot ./cmd/spam_fighter_bot/spam_fighter_bot.go

FROM alpine:3.13.5
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/pcherednichenko/spam_fighter_bot/spam_fighter_bot spam_fighter_bot

CMD ["./spam_fighter_bot"]