FROM obraun/golang-micro-ci as builder
COPY . /app
WORKDIR /app
#ENV GO111MODULE=on
RUN go build -o client/contentpreloader/service/contentpreloader client/contentpreloader/service/main.go

FROM alpine
COPY --from=builder /app/client/contentpreloader/service/contentpreloader /app/contentpreloader
EXPOSE 8091
ENTRYPOINT [ "/app/contentpreloader" ]
RUN chmod +x ./app/contentpreloader