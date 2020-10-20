FROM obraun/golang-micro-ci as builder
COPY . /app
WORKDIR /app
#ENV GO111MODULE=on
RUN go build -o evaluation/service/evaluation-service evaluation/service/main.go

FROM alpine
COPY --from=builder /app/evaluation/service/evaluation-service /app/evaluation-service
EXPOSE 8091
ENTRYPOINT [ "/app/evaluation-service" ]
RUN chmod +x ./app/evaluation-service