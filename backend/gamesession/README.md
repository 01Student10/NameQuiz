# Game Session Service:

## Beschreibung:

+ Der Game Session Service regelt den Ablauf eines Spiels.

+ Will man ein Spiel spielen, kann eine Session erstellt werden. Das Game bleibt erhalten, die Session wird, nach dem Ende des Spiels wieder gelöscht.

+ Der Game Session Service hält keine Daten, die in einer Datenbank gespeichert werden.

+ Der Game Session Service ist unter dem Namen: "nq.GameSessionService" ansprechbar. Wobei "nq" der Namespace der Namequiz App ist.

+ Der Game Session Service ist nicht skalierbar. Er ist durch die Ports, die wegen der Websocket Connections aus dem docker-compose Netzwerk exportiert werden müssen eingeschränkt. Zusätzlich ist er durch seine Datenhaltung nicht Zustandslos.

+ Der Game Session Service kann ausschließlich über sein definiertes Interface angesprochen werden.

## Daten:
Der Service hält folgende Daten lokal:

Das FrontendMessage Struct sind all die Daten, die an das Frontend übermittelt werden.
<pre>
type FrontendMessage struct {
	Status          string
	UserNameToScore map[string]int32
	SessionID       int32
	EntryID         string
	Path            string
	PicInfo         string
	Winner          []Pair
	ErrorMsg        string
	GameTitle       string
}
</pre>
Dieses Struct definiert die Daten, die in einer Session vorhanden sind.
<pre>
type Session struct {
	sessionID      int32
	currentStatus  Status
	userToScore    map[string]int32
	hasScored      map[string]struct{}
	currentGameID  int32
	entryListIndex int32
	hub            Hub
	isRunning      int32
	chatAddress    string
	sessionAddress string
}
</pre>
Das Service Struct enthält die Daten, die der Service braucht um seine Sessions zu verwalten.
<pre>
type Service struct {
	// Clients
	userService        api.UserService
	chatService        api.ChatService
	gameService        api.GameService
	quizContentService api.QuizContentService
	// Variables
	sessionIDIndex     int32
	roundTime          int32
	sessions           map[int32]*Session
	chatSocketPorts    map[string]bool
	sessionSocketPorts map[string]bool
}
</pre>

Das Client Struct enthält eine Websocket Connection und einen byte array Chanel auf den daten geschrieben werden.
<pre>
type Client struct {
	// The websocket connection.
	Conn *websocket.Conn
	// Buffered channel of outbound messages.
	Send chan []byte
}
</pre>

## Methoden:
das Interface ist im [api.Proto](../api/api.proto) File folgendermaßen definiert:

<pre>
service GameSessionService {
  rpc CreateSession(CreateSessionRequest) returns (CreateSessionResponse) {}
  rpc StartSession(StartSessionRequest) returns (StartSessionResponse) {}
  rpc JoinSession(JoinSessionRequest) returns (JoinSessionResponse) {}
  rpc LeaveSession(LeaveSessionRequest) returns (LeaveSessionResponse) {}
  rpc IsSessionCreated(IsSessionCreatedRequest) returns (IsSessionCreatedResponse) {}
}
 </pre>

## Abhängigkeiten:
- Game Service
- QuizContent Service
- User Service
- Chat Service
- etcd (Service Discovery)
- NATS (Message Broker)

## Good to know:

+ Im [docker-compose.yaml](../docker-compose.yaml) File werden dem Microservice die Adresse der Service Discovery und des Message Broker Services folgendermaßen übergeben:

 `command: --registry_address=etcd:2379 --broker_address=nats:4222`

+ Da `boolVar := false` in Protocol Buffer als default Wert gilt, wird dieser nicht encoded. Deswegen werden Boolean Parameter als `int32` übertragen. Wobei `1` für `true` und `2` für `false` stehen.
