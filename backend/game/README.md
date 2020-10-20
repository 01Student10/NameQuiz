# Game Service:

## Beschreibung:

+ Der Game Service repräsentiert ein Spiel in der Namequiz Anwendung.

+ Der Game Service hält Daten, die er in einem Datenbank  Service verwaltet.
+ Um einen neuen Eintrag in die Datenbank zu schreiben, benötigt er eine eindeutige ID, die er sich vom ID Service holt.

+ Der Game Service ist unter dem Namen: "nq.GameService" ansprechbar. Wobei "nq" der Namespace der Namequiz App ist.

+ Der Game Service ist voll skalierbar. Durch ein ClientWrapper Plugin wird er mithilfe der Service Discovery lastverteilt.

+ Der Game Service kann ausschließlich über sein definiertes Interface angesprochen werden.

## Daten:
Das Struct ist im [api.Proto](../api/api.proto) File folgendermaßen definiert:
<pre>
message Game {
  int32 id = 1;
  string title = 2;
  int32 amountOfRounds = 3;
  string playlistId = 4;
  int32 owner = 6;
}
</pre>

## Methoden:
das Interface ist im [api.Proto](../api/api.proto) File folgendermaßen definiert:

<pre>
service GameService {
  rpc CreateGame(CreateGameRequest) returns (CreateGameResponse) {}
  rpc DeleteGame(DeleteGameRequest) returns (DeleteGameResponse) {}
  rpc GetGame(GetGameRequest) returns (GetGameResponse) {}
  rpc GetAllGames(GetAllGamesRequest) returns (GetAllGamesResponse) {}
}
 </pre>

## Abhängigkeiten:
- Id Service
- etcd (Service Discovery)
- redis Store Service (Datenbank)

## Good to know:

+ Im [docker-compose.yaml](../docker-compose.yaml) File werden dem Microservice die Adresse der Service Discovery und seines Datenbank Services folgendermaßen übergeben:

 `command: --registry_address=etcd:2379 --store_address=db-game:6380`

+ Da `boolVar := false` in Protocol Buffer als default Wert gilt, wird dieser nicht encoded. Deswegen werden Boolean Parameter als `int32` übertragen. Wobei `1` für `true` und `2` für `false` stehen.
