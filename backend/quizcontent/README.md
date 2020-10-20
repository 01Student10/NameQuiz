# QuizContent Service:

## Beschreibung:

+ Der QuizContent Service repräsentiert den Spiel Inhalt der Namequiz Anwendung.

+ Content kann nur über den Contentpreloader Service in den QuizContent Service geladen werden.

+ Der QuizContent Service hält Daten, die er in einem Datenbank Service verwaltet.
+ Die Id des Eintrags steht bereits im JSON File, deswegen muss keine einzigartige ID vom Id Service angefragt werden.

+ Der QuizContent Service ist unter dem Namen: "nq.QuizContentService" ansprechbar. Wobei "nq" der Namespace der Namequiz App ist.

+ Der QuizContent Service ist voll skalierbar. Durch ein ClientWrapper Plugin wird er mithilfe der Service Discovery lastverteilt.

+ Der QuizContent Service kann ausschließlich über sein definiertes Interface angesprochen werden.

## Daten:
Das Content Struct sieht folgendermaßen aus:
<pre>
// The Content is stored into the store. It is not listed in the api.proto file.
type Content struct {
	// the ID of the playlist.
	ListId string
	// The playlist name.
	Name string
	// Contains id string to a ContentEntry struct(which is defined in api.proto)
	Entries map[string]*api.ContentEntry
}
</pre>

Das ContentEntry Struct ist im [api.Proto](../api/api.proto) File folgendermaßen definiert:
<pre>
message ContentEntry {
  string id = 1;
  string name = 2;
  string path = 3;
  string licence = 4;
}
</pre>

Eine Playlist kann im JSON-Format in den QuizContent Service geladen werden. Die Struktur sieht folgendermaßen aus.

<pre>
{
  "listId": "1",
  "name": "actors",
  "entries": {
    "2": {
      "id": "2",
      "name": "Jack Nicholson",
      "path": "Jack_Nicholson_2001_wikipedia.jpg",
      "licence": "From Kennedy center ( an honoree that year ) December 2, 2001 © copyright John Mathew Smith 2001 -Licence: CC BY-SA 2.0 "
    },
    "3": {
      "id": "3",
      "name": "Uma Thurman",
      "path": "Uma_Thurman_-_Cannes_2000.jpg",
      "licence": "Where: in Cannes Red Carpet in 2000 - Licence: CC BY-SA 2.5 - Creator: Rita Molnár"
    },
    "4": {
      "id": "4",
      "name": "Will Ferrell",
      "path": "Will_Ferrell_Deauville_2014.jpg",
      "licence": "at the Deauville Film Festival - Licence: CC BY-SA 3.0 - Creator: Georges Biard"
    },
    "5": {
      "id": "5",
      "name": "Mark Ruffalo",
      "path": "Mark_Ruffalo_in_2017_by_Gage_Skidmore.jpg",
      "licence": "auf der San Diego Comic-Con International (2017)- Licence: CC BY-SA 2.0 - Creator: Gage Skidmore"
    },
    "6": {
      "id": "6",
      "name": "Emilia Clarke",
      "path": "Emilia_Clarke_by_Gage_Skidmore_2_(cropped).jpg",
      "licence": "auf der San Diego Comic-Con International 2013 - Licence: CC BY-SA 3.0 - Creator: Gage Skidmore "
    },
    "7": {
      "id": "7",
      "name": "Emma Watson",
      "path": "Emma_Watson_2013.jpg",
      "licence": " at the Cannes Film Festival 2013. - Licence: CC BY-SA 3.0 -  Creator: Georges Biard "
    },
    "8": {
      "id": "8",
      "name": "Edward Norton",
      "path": "Edward_Norton_2012.jpg",
      "licence": " in Africa, March 2012. - Licence: CC BY 2.0 - Creator: Steve Jurvetson  "
    },
    "9": {
      "id": "9",
      "name": "Russel Crowe",
      "path": "Russell_Crowe.jpg",
      "licence": "Licence: CC BY-SA 2.0 - Creator: Eva Rinaldi"
    },
    "10": {
      "id": "10",
      "name": "Daisy Ridley",
      "path": "Daisy_Ridley_by_Gage_Skidmore.jpg",
      "licence": "at the 2015 San Diego Comic Con International in San Diego, California. - Licence: CC BY-SA 3.0 - Creator: Gage Skidmore"
    },
    "11": {
      "id": "11",
      "name": "Will Smith",
      "path": "Will_Smith_by_Gage_Skidmore_2.jpg",
      "licence": "bei der Comic-Con (2017) - Licence: CC BY-SA 3.0 - Creator: Gage Skidmore"
    },
    "12": {
      "id": "12",
      "name": "Kevin James",
      "path": "Kevin_James_2011_(Cropped).jpg",
      "licence": "at a ceremony to receive a star on the Hollywood Walk of Fame. - Licence: CC BY-SA 3.0 - Creator: Angela George"
    }
  }
}

</pre>


## Methoden:
das Interface ist im [api.Proto](../api/api.proto) File folgendermaßen definiert:

<pre>
service QuizContentService {
  rpc GetContentEntry(GetContentEntryRequest) returns (GetContentEntryResponse) {}
  rpc GetContentList(GetContentListRequest) returns (GetContentListResponse) {}
  rpc GetAllContentLists(GetAllContentListsRequest) returns (GetAllContentListsResponse) {}
  rpc HasMatch(HasMatchRequest) returns (HasMatchResponse) {}
  rpc CreateNewDataSet(CreateNewDataSetRequest) returns (CreateNewDataSetResponse) {}
}
 </pre>

## Abhängigkeiten:
- etcd (Service Discovery)
- redis Store Service (Datenbank)

## Good to know:

+ Im [docker-compose.yaml](../docker-compose.yaml) File werden dem Microservice die Adresse der Service Discovery und seines Datenbank Services folgendermaßen übergeben:

 `command: --registry_address=etcd:2379 --store_address=db-quizcontent:6381`

+ Da `boolVar := false` in Protocol Buffer als default Wert gilt, wird dieser nicht encoded. Deswegen werden Boolean Parameter als `int32` übertragen. Wobei `1` für `true` und `2` für `false` stehen.

+ Um auf die JSON-Files zugreifen zu können, müssen diese im Docker Container des QuizContent Services verfügbar sein. Folgender Code zeigt den Ausschnitt des Dockerfiles, an dem dies geschieht:
<pre>
COPY client/contentpreloader/data/actors.json .
COPY client/contentpreloader/data/musicians.json .
COPY client/contentpreloader/data/scientists.json .
</pre>