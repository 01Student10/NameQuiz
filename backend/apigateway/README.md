# API Gateways:

## Beschreibung:

+ Die Api Gateways sind dafür zuständig vom Frontend kommende Requests zu übersetzen und an den entsprechenden Service im Backend weiterzuleiten.

## API Gateway RPC:
+ Das API Gateway RPC ist ein bereits bestehender Service aus dem micro framework.
+ Es ist dafür zuständig eingehende HTTP-REST-Requests in RPC-Requests umzuwandeln. Gleiches gilt für Responses.

### Konfiguration:
Die Konfiguration wird im [docker-compose.yaml](../docker-compose.yaml) vorgenommen.
<pre>
 apigateway-rpc:
    image: "micro/micro:v2.9.2"
    environment:
      - ALLOW_NONE_AUTHENTICATION=yes
      - MICRO_API_HANDLER=rpc
      - MICRO_API_REGISTRY=etcd
      - MICRO_API_NAMESPACE="sq"
      - MICRO_API_ENABLE_RPC=true
      - MICRO_API_ADDRESS="8085"
    command: --auth_namespace=sq --registry=etcd --registry_address=etcd:2379 api --handler=rpc --address=:8085 --namespace=sq
    links:
      - etcd
    ports:
      - "8085:8085"
    depends_on:
      - etcd
</pre>

Hierbei ist wichtig, dass das Gateway Anfragen an den port 8085 in RPC-Requests übersetzt. Dies geschieht mit dem Befehl:
`api --handler=rpc --address=:8085`.
Zusätzlich muss der Port exportiert werden.

### Abhängigkeiten:
- etcd (Service Discovery)


## API Gateway event:
+ Das API Gateway event ist ein bereits bestehender Service aus dem micro framework.
+ Es ist dafür zuständig eingehende HTTP-REST-Requests in EventMessages umzuwandeln und an den Message Broker zu senden.

### Konfiguration:
Die Konfiguration wird im [docker-compose.yaml](../docker-compose.yaml) vorgenommen.
<pre>
  apigateway-event:
    image: "micro/micro:v2.9.2"
    environment:
      - ALLOW_NONE_AUTHENTICATION=yes
      - MICRO_API_HANDLER=event
      - MICRO_API_REGISTRY=etcd
      - MICRO_API_NAMESPACE=nq
      - MICRO_API_ENABLE_RPC=true
      - MICRO_API_ADDRESS="8086"
    command: --auth_namespace=nq --broker_address=nats:4222 --broker=nats --registry=etcd --registry_address=etcd:2379 api --handler=event --address=:8086 --namespace=nq
    links:
      - etcd
    ports:
      - "8086:8086"
    depends_on:
      - etcd
      - nats
</pre>
Hierbei ist wichtig, dass das Gateway Anfragen an den port 8085 in RPC-Requests übersetzt. Dies geschieht mit dem Befehl:
`api --handler=event --address=:8086`.
Zusätzlich muss der Port exportiert werden.

### Abhängigkeiten:
- etcd (Service Discovery)
- NATS (Message Broker)

