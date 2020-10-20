# Content Preloader Service:

## Beschreibung:

+ Der Content Preloader Service hat die Aufgabe neuen Content in die NameQuiz app zu laden.

+ Der Content Preloader Service hält Daten, die er in einem Datenbank  Service verwaltet.
+ Um einen neuen Eintrag in die Datenbank zu schreiben, benötigt er eine eindeutige ID, die er sich vom ID Service holt.

+ Der Content Preloader Service ist unter dem Namen: "nq.Content PreloaderService" ansprechbar. Wobei "nq" der Namespace der Namequiz App ist.

+ Der Content Preloader Service ist voll skalierbar. Durch ein ClientWrapper Plugin wird er mithilfe der Service Discovery lastverteilt.

+ Der Content Preloader Service kann ausschließlich über sein definiertes Interface angesprochen werden.

## Daten:
Im [data](data) Ordner befinden sich die .json Files die dem QuizContent Service übergeben werden. 
Diese müssen in den Docker Container des QuizContent Services kopiert werden wie es im [Dockerfile](../../quizcontent/service/Dockerfile) des QuizContent Service zu sehen ist.

## Methoden:
<pre>
func (ci *ContentPreLoader) InitializeData(path string) {..}
 </pre>
 
Lädt ein neuen Content in den QuizContent Service.

## Abhängigkeiten:
- QuizContent Service
- etcd (Service Discovery)