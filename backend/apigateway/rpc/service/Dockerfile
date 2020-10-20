FROM golang as builder
COPY . /app
WORKDIR /app
#ENV GO111MODULE=on

FROM alpine
COPY --from=builder /app/apigateway/rpc/service/apigateway-rpc /app/apigateway-rpc
EXPOSE 8080
ENTRYPOINT [ "/app/apigateway-rpc" ]
RUN chmod +x ./app/apigateway-rpc