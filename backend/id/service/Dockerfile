FROM obraun/golang-micro-ci as builder
COPY . /app
WORKDIR /app
RUN go build -o id/service/id-service id/service/main.go

FROM alpine
COPY --from=builder /app/id/service/id-service /app/id-service
EXPOSE 8091
ENTRYPOINT [ "/app/id-service" ]
RUN chmod +x ./app/id-service