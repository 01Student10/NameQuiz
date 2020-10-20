# Id Service:

## Beschreibung:

+ Der Id Service stellt den Microservices eine eindeutige ID zur Verfügung.
+ Die ID ist eine fortlaufende Nummer vom typ `int32`, die in einem Datenbank Service verwaltet wird.
+ Der Id Service ist unter dem Namen: "nq.IdService" ansprechbar. Wobei "nq" der Namespace der Namequiz App ist.
+ Der Id Service ist voll skalierbar. Durch ein ClientWrapper Plugin wird er mithilfe der Service Discovery lastverteilt.
+ Der Id Service kann ausschließlich über sein definiertes Interface angesprochen werden.

## Daten:
Das private Struct ist folgendermaßen definiert:
<pre>
type Id struct {
	idCounter int32
}
</pre>

## Methoden:
das Interface ist im [api.Proto](../api/api.proto) File folgendermaßen definiert:

<pre>
service IdService {
  rpc GetId(GetIdRequest) returns (GetIdResponse) {}
}
 </pre>

## Abhängigkeiten:
- etcd (Service Discovery)
- redis Store Service (Datenbank)

## Good to know:

+ Im [docker-compose.yaml](../docker-compose.yaml) File werden dem Microservice die Adresse der Service Discovery und seines Datenbank Services folgendermaßen übergeben:

 `command: --registry_address=etcd:2379 --store_address=db-id:6382`

+ Da `boolVar := false` in Protocol Buffer als default Wert gilt, wird dieser nicht encoded. Deswegen werden Boolean Parameter als `int32` übertragen. Wobei `1` für `true` und `2` für `false` stehen.
