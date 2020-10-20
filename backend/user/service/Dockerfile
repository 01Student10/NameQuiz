FROM obraun/golang-micro-ci as builder
COPY . /app
WORKDIR /app
RUN go build -o user/service/user-service user/service/main.go

FROM alpine
COPY --from=builder /app/user/service/user-service /app/user-service
EXPOSE 8091
ENTRYPOINT [ "/app/user-service" ]
RUN chmod +x ./app/user-service