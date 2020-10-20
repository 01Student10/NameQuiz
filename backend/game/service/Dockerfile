FROM obraun/golang-micro-ci as builder
COPY . /app
WORKDIR /app
#ENV GO111MODULE=on
RUN go build -o game/service/game-service game/service/main.go

FROM alpine
COPY --from=builder /app/game/service/game-service /app/game-service
EXPOSE 8091
ENTRYPOINT [ "/app/game-service" ]
RUN chmod +x ./app/game-service