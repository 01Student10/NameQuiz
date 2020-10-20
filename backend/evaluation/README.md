# Evaluation Service:

## Beschreibung:

+ Der Evaluation Service ist in der NameQuiz Anwendung dafür zuständig, eingehende Nachrichten zu evaluieren. Hierbei werden vom Benutzer angegebene Namen mit dem Richtigen Namen verglichen und anschließend an den GameSession Service gesendet.

+ Da der Evaluation Service selbst keine Daten hält, implementiert er einen quizContent Service Client, über den er die Richtigkeit der agbegebenen Tipps der Spieler überprüfen kann.

+ Der Evaluation Service ist unter dem Namen: "nq.EvaluationService" ansprechbar. Wobei "nq" der Namespace der Namequiz App ist.

+ Der Evaluation Service ist voll skalierbar. Da eine Evaluation Service Instanz der NATS Queue Group "evaluation" als Subscriber beitritt, kann der Nats Server Nachrichten unter den Instanzen lastverteilen.

+ Der Evaluation Service Subscribed dem Topic "nq.api.Chat" um Messages vom apigateway-event Service zu erhalten.

+ Der Evaluation Service Published Nachrichten zum Topic "nq.ChatMessage" die vom GameSession Service empfangen werden.

## Daten:
Messages, die innerhalb des Backends übertragen werden haben folgendes Format und sind im [api.Proto](../api/api.proto) File folgendermaßen definiert:
<pre>
message UserMessage {
  int32 userId = 1;
  string userName = 2;
  string guess = 3;
  int32 sessionId = 4;
  string listId = 5;
  string entryId = 6;
  int32 wasRight = 7;

}
</pre>

Messages die vom apigateway-event Service übertragen werden haben folgendes Format und sind im [api.proto](https://github.com/micro/go-micro/blob/master/api/proto/api.proto) Files des go-micro Frameworks definiert:

<pre>
// A HTTP event as RPC
// Forwarded by the event handler
message Event {
	// e.g login
	string name = 1;
	// uuid
	string id = 2;
	// unix timestamp of event
	int64 timestamp = 3;
	// event headers
        map<string, Pair> header = 4;
	// the event data
	string data = 5;
}
</pre>

## Methoden:
Wird eine Message zum Topic "np.api.Chat" vom apigateway-event veröffentlicht, wird folgende Methode aufgerufen:

<pre>
func (s *Evaluation) Handle(_ context.Context, input *gatewayProto.Event) error {...}
 </pre>

## Abhängigkeiten:
- QuizContent Service
- NATS (Message Broker)
- etcd (Service Discovery)

## Good to know:

+ Im [docker-compose.yaml](../docker-compose.yaml) File werden dem Microservice die Adresse der Service Discovery und des Message Brokers folgendermaßen übergeben:

 `command: --registry_address=etcd:2379 --broker_address=nats:4222`

+ Da `boolVar := false` in Protocol Buffer als default Wert gilt, wird dieser nicht encoded. Deswegen werden Boolean Parameter als `int32` übertragen. Wobei `1` für `true` und `2` für `false` stehen.

+ Wenn in der `micro.RegisterSubscriber()` Methode in der [main.go](service/main.go) ein Evaluation Objekt übergeben wird, werden beim Veröffentlichen zu einem Topic alle in dem Objekt definierten Methoden aufgerufen.
Im Falle des Evaluation Services ist dies nur die Methode `Handle()`.


