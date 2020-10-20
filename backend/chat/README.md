# Chat Service:

## Beschreibung:

+ Der Chat Service ist in der NameQuiz Anwendung dafür zuständig Chat Nachrichten an das Frontend zu schicken.

+ Der Chat Service hält selbst keine Daten in einer Datenbank.

+ Der Chat Service ist unter dem Namen: "nq.ChatService" ansprechbar. Wobei "nq" der Namespace der Namequiz App ist.

+ Der Chat Service ist nicht skalierbar. Er ist durch die Ports, die wegen der Websocket Connections aus dem docker-compose Netzwerk exportiert werden müssen eingeschränkt. Zusätzlich ist er durch seine Datenhaltung nicht Zustandslos.

+ Der Chat Service kann ausschließlich über sein definiertes Interface angesprochen werden.

## Daten:
In socketClients wird eine Session ID zu einem Hub Struct gemappt.
Dies wird benötigt, um zur passenden Session Nachrichten an das Frontend zu schicken. 
<pre>
socketClients map[int32]Hub
</pre>

## Methoden:
das Interface ist im [api.Proto](../api/api.proto) File folgendermaßen definiert:

<pre>
service ChatService {
  rpc BroadcastMessage(BroadcastMessageRequest) returns (BroadcastMessageResponse) {}
  rpc StartWebsocket(StartWebsocketRequest) returns (StartWebsocketResponse) {}
}
 </pre>

## Abhängigkeiten:
- NATS (Message Broker)
- etcd (Service Discovery)

## Good to know:

+ Im [docker-compose.yaml](../docker-compose.yaml) File werden dem Microservice die Adresse der Service Discovery und des Message Brokers folgendermaßen übergeben:
<pre>
 environment:
       - ALLOW_NONE_AUTHENTICATION=yes
       - MICRO_BROKER_ADDRESS=nats:4222
       - MICRO_REGISTRY_ADDRESS=etcd:2379
</pre>

+ Da `boolVar := false` in Protocol Buffer als default Wert gilt, wird dieser nicht encoded. Deswegen werden Boolean Parameter als `int32` übertragen. Wobei `1` für `true` und `2` für `false` stehen.

