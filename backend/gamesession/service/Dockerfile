FROM obraun/golang-micro-ci as builder
COPY . /app
WORKDIR /app
#ENV GO111MODULE=on
RUN go build -o gamesession/service/gamesession-service gamesession/service/main.go

FROM alpine
COPY --from=builder /app/gamesession/service/gamesession-service /app/gamesession-service
EXPOSE 8091
ENTRYPOINT [ "/app/gamesession-service" ]
RUN chmod +x ./app/gamesession-service