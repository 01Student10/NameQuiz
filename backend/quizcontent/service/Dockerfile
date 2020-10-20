FROM obraun/golang-micro-ci as builder
COPY . /app
WORKDIR /app
RUN go build -o quizcontent/service/quizcontent-service quizcontent/service/main.go

FROM alpine
COPY client/contentpreloader/data/actors.json .
COPY client/contentpreloader/data/musicians.json .
COPY client/contentpreloader/data/scientists.json .
COPY --from=builder /app/quizcontent/service/quizcontent-service /app/quizcontent-service
EXPOSE 8091
ENTRYPOINT [ "/app/quizcontent-service" ]
RUN chmod +x ./app/quizcontent-service