FROM golang as builder
COPY . /app
WORKDIR /app

FROM alpine
COPY --from=builder /app/apigateway/event/service/apigateway-event /app/apigateway-event
EXPOSE 8080
ENTRYPOINT [ "/app/apigateway-event" ]
RUN chmod +x ./app/apigateway-event