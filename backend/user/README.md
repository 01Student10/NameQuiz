# User Service:

## Beschreibung:

+ Der User Service repräsentiert einen Nutzer in der Namequiz Anwendung.

+ Der User Service hält Daten, die er in einem Datenbank  Service verwaltet.
+ Um einen neuen Eintrag in die Datenbank zu schreiben benötigt er eine eindeutige ID, die er sich vom ID Service holt.

+ Der User Service ist unter dem Namen: "nq.UserService" ansprechbar. Wobei "nq" der Namespace der Namequiz App ist.

+ Der User Service ist voll skalierbar. Durch ein ClientWrapper Plugin wird er mithilfe der Service Discovery lastverteilt.

+ Der User Service kann ausschließlich über sein definiertes Interface angesprochen werden.

## Daten:
Das Struct ist im [api.Proto](../api/api.proto) File folgendermaßen definiert:
<pre>
message User {
  int32 id = 1;
  string username = 2;
  string password = 3;
  string mail = 4;
}
</pre>

## Methoden:
das Interface ist im [api.Proto](../api/api.proto) File folgendermaßen definiert:

<pre>
service UserService{  
   rpc CreateUser(CreateUserRequest) returns (CreateUserResponse) {}  
   rpc GetUser(GetUserRequest) returns (GetUserResponse) {}    
   rpc GetAllUsers(GetAllUsersRequest) returns (GetAllUsersResponse) {}  
   rpc HasUser(HasUserRequest) returns (HasUserResponse) {}  
   rpc Login(LoginRequest) returns (LoginResponse) {}    
 }
 </pre>

## Abhängigkeiten:
- Id Service
- etcd (Service Discovery)
- redis Store Service (Datenbank)

## Good to know:

+ Im [docker-compose.yaml](../docker-compose.yaml) File werden dem Microservice die Adresse der Service Discovery und seines Datenbank Services folgendermaßen übergeben:

 `command: --registry_address=etcd:2379 --store_address=db-user:6383`

+ Da `boolVar := false` in Protocol Buffer als default Wert gilt, wird dieser nicht encoded. Deswegen werden Boolean Parameter als `int32` übertragen. Wobei `1` für `true` und `2` für `false` stehen.
