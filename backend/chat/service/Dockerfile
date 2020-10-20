FROM obraun/golang-micro-ci as builder
COPY . /app
WORKDIR /app
#ENV GO111MODULE=on
RUN go build -o chat/service/chat-service chat/service/main.go

FROM alpine
COPY --from=builder /app/chat/service/chat-service /app/chat-service
EXPOSE 8091
ENTRYPOINT [ "/app/chat-service" ]
RUN chmod +x ./app/chat-service