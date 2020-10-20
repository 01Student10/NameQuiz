# Database Services:

## Beschreibung:

+ Die Database Services stellen für die Services die Daten in einer Datenbank halten einen Server zur Verfügung.
+ Die Datenbanken sind redis Caches die eine zusätzliche Persistenz mitbringen. 
+ Im [redis.conf](dbgame/redis.conf) (hier am Beispiel des dbgame Datenbank Servers) wird der Redis Datenbank Server konfiguriert.
+ Um dies zu erreichen muss in jedem Container in dem die Datenbank läuft das entsprechende Config File kopiert werden.
+ Um eine Persistenz zu erhalten setzt redis auf zwei Verfahren:
1. RDB - in bestimmten Zeitintervall wird eine Momentaufnahme des Datensets generiert und in eine Datei [dump.rdb](dbgame/data/dump.rdb) (Bsp.: dbgame) geschrieben.
2. AOF - jeder Schreibzugriff auf den Server wird in Log-Datei [appendonly.aof](dbgame/data/appendonly.aof) (Bsp.: dbgame) gespeichert. 



## Konfiguration (am Beispiel dbgame):
Das [Dockerfile](dbgame/service/Dockerfile)sieht folgendermaßen aus:  
<pre>
 FROM bitnami/redis:latest
 COPY database/dbgame/redis.conf /opt/bitnami/redis/mounted-etc/redis.conf
 CMD ["redis-server","/opt/bitnami/redis/mounted-etc/redis.conf", "--appendonly", "yes"]
 USER 1001
 EXPOSE 6380
</pre>
+ Es ist wichtig das der Server mit der Konfigurationsdatei und dem flag `--appendonly=yes` gestartet wird.
+ Zusätzlich muss der User im Docker Container 1001 sein, da dieser im redis Container Zugriff auf die Daten hat.

Im [docker-compose.yaml](../docker-compose.yaml) File muss der data Ordner gemountet werden.

<pre>
volumes:
    - ./database/dbgame/data/:/opt/bitnami/redis/data/
</pre>

Dies stellt sicher, das bei einem herunterfahren des docker-compose Netzwerk die Daten der Services 
in die beiden Dateien [appendonly.aof](dbgame/data/appendonly.aof) und  [dump.rdb](dbgame/data/dump.rdb) gespeichert werden.

**Wichtig**

* Bevor der Datenbank Service gestartet werden kann, muss die [redis.conf](dbgame/redis.conf) im Repository ausführbar sein.
Dies wird mit folgendem Befehl bewerkstelligt:
`sudo chmod 777 redis.conf` 
* Desweiteren muss dem data Ordner lesende und schreibende rechte zugewiesen werden.
`sudo chmod 777 data`
* Zum Schluss muss der data Ordner dem User 1001 übertragen werden.
`sudo chown 1001:1001 data`
* Das Sudo Passwort ist gleich dem Usernamen: "vss".

* Die [redis.conf](dbgame/redis.conf) unterscheiden sich in Ihrem Port:
Zeile 92:
`port 6380`
Zeile 244:
`pidfile /var/run/redis_6380.pid`


## Good to know:
+ Wenn die Datenbank zurückgesetzt werden soll, müssen im data Ordner die Dateien [appendonly.aof](dbgame/data/appendonly.aof) und  [dump.rdb](dbgame/data/dump.rdb) gelöscht werden:
`sudo rm appendonly.aof dump.rdb`
