# Aufgabenstellung
Es soll ein einfaches Ticketsystem implementiert werden. Diese soll es Kunden ermöglichen,
Tickets zu erstellt. Diese sollen dann durch eine Gruppe von Bearbeitern bearbeitet werden
können. Dabei soll allen Bearbeitern der komplette Verlauf zugänglich sein.

# Anforderungen
## Nicht funktional
1. Der Code soll im Paket de/vorlesung/projekt/[Gruppen-ID] liegen, damit
alle Lösungen parallel im Source-Tree gehalten werden können.
2. Es durfen keine Pakete Dritter verwendet werden! Einzige Ausnahme sind
Pakete zur Vereinfachung der Tests und der Fehlerbehandlung.
Empfohlen sei hier github.com/stretchr/testify/assert und
github.com/pkg/errors.
3. Alle Source-Dateien und die PDF-Datei mussen die Matrikelnummern aller 
Gruppenmitglieder enthalten.
## Allgemein
4. Die Anwendung muss nur einen Mandanten bzw. ein Projekt verwalten k¨onnen.
Mehrere Ticketsysteme auf dem selben Server sind nicht erforderlich.
5. Die Anwendung soll unter Windows und Linux lauff¨ahig sein.
6. Es soll sowohl IE 11 als auch Firefox, Chrome und Edge in der jeweils aktuellen
Version unterstutzt werden. Diese Anforderung ist am einfachsten zu 
erfullen, indem Sie auf komplexe JavaScript/CSS "spielereien“ verzichten. ;-)
## Sicherheit
7. Die Web-Seite soll nur per HTTPS erreichbar sein.
8. Der Zugang fur die Bearbeiter soll durch Benutzernamen und Passwort gesch ¨ utzt ¨
werden.
10. Die Passwörter durfen nicht im Klartext gespeichert werden. 
11. Es soll ”salting“ eingesetzt werden.
12. Alle Zugangsdaten sind in einer gemeinsamen Datei zu speichern.
## Ticket
13. Ein Ticket enthält
- eine eindeutige Kennung
- einen Betreff
- den Status des Tickets (offen; in Bearbeitung; geschlossen)
- Kennung des Bearbeiters, wenn im Status ”in Bearbeitung“
- Eine Reihe von Einträgen.
14. Ein Eintrag enthält
- Datum der Erstellung
- Entweder E-Mail des Kunden oder Name des Bearbeiters.
- Text des Eintrags
## Ticketerstellung über eine Web-Seite 
15. Über eine Web-Seite soll ein Ticket erstellt werden können.
16. Die Erzeugung eines Tickets soll ohne eine Authentifizierung möglich sein.
17. Erfasst werden sollen
- E-Mail Adresse
- Betreffzeile
- der eigentlicher Text des Tickets
## E-Mail-Empfang über eine REST-API
Später soll ein Dienst implementiert werden, welcher ein oder mehrere E-MailKonten
fur die Ticketerstellung verwendet. Dieser soll ein E-Mail-Konto abfragen
und alle eingehenden Nachrichten in das Ticketsystem einspeisen. Dieser E-Mail Abfrage-Dienst
selbst ist nicht Bestandteil der Aufgabenstellung. 

18. Es soll eine Funktion geben, uber welche ein externer Dienst Nachrichten
abgeben kann.
19. Die Funktion soll folgende Parameter haben:
- E-Mail Absendeadresse
- Betreffzeile
- der eigentliche Text der Nachricht
20. Beim Aufruf soll uberprüft werden, ob sich die Nachricht auf ein bereits existierendes Ticket bezieht.
21. Wenn ein Ticket gefunden wurde,
- soll es diesem als neuer Eintrag hinzugefugt werden. 
- und das Ticket hat den Status ”geschlossen“, soll es auf ”offen“ zuruckgesetzt werden.
22. Wurde kein bestehendes Ticket gefunden, soll ein neues Ticket aus der eingegangenen
Nachricht erzeugt werden.
23. Es soll ein einfaches Kommandozeilen-Tool geben, mit welchem Nachrichten
an den Server ubertragen werden können.
## E-Mail-Versand uber eine REST-API 
Später soll ein Dienst implementiert werden, welcher ein E-Mail-Konto fur den
Nachrichtenversand verwendet. Dieser soll beim Server zu versendende E-Mails
abfragen, diese verschicken und danach das Verschicken bestätigen. Dieser E-Mail3
Versende-Dienst selbst ist nicht Bestandteil der Aufgabenstellung.

24. Es soll eine Funktion geben, uber welche alle E-Mails abgerufen werden 
können, die noch zu verschicken sind.
25. Es soll eine Funktion geben, uber welche ein externer Dienst dem Server 
mitteilen kann, welche der E-Mails er verschickt hat.
26. E-Mail-Ping-Pong (Abwesenheitsassistent o.ä.) soll verhindert werden.
27. Es soll ein einfaches Kommandozeilen-Tool geben, mit welchem zu versendende
Nachrichten auf der Konsole ausgegeben werden.
## Bearbeitung der Tickets
28. Die Bearbeitung der Tickets soll ausschließlich uber eine WEB-Seite erfolgen. 
29. Bearbeiter sollen ein Ticket ubernehmen können.
30. Bearbeiter sollen alle Tickets einsehen können, welche noch kein Bearbeiter
ubernommen hat. 
31. Bearbeiter sollen Tickets nach der Ubernahme auch freigeben können, so das
diese eine anderer Bearbeiter ubernehmen kann. 
32. Ein Bearbeiter soll ein Ticket einem anderen Bearbeiter zuteilen können.
33. Kommentiert ein Bearbeiter ein Ticket, soll er wählen können, ob dieser Kommentar
an den Kunden versendet wird, oder ob er nur fur andere Bearbeiter sichtbar ist.
## Storage
34. Die Tickets sollen zusammen mit dem kompletten Verlaufs einer Datei im
Dateisystem gespeichert werden.
35. Es sollen nicht alle Tickets in einer gemeinsamen Datei gespeichert werden.
36. Es soll ein geeignetes Caching implementiert werden, d.h. es sollen nicht bei
jedem Request alle Dateien neu eingelesen werden.
## Konfiguration
37. Die Konfiguration soll komplett uber Startparameter erfolgen. (Siehe Package flag)
38. Der Port muss sich uber ein Flag festlegen lassen. 
39. Hart kodierte“ absolute Pfade sind nicht erlaubt.
## Betrieb
40. Wird die Anwendung ohne Argumente gestartet, soll ein sinnvoller default gewählt werden.
41. Nicht vorhandene aber ben¨otigte Order sollen ggfls. angelegt werden.
42. Die Anwendung soll zwar HTTPS und die entsprechenden erforderlichen Zertifikate
unterstutzen, es kann jedoch davon ausgegangen werden, dass geeignete Zertifikate gestellt werden. Fur Ihre Tests können Sie ein 
"self signed“ Zertifikat verwenden. Es ist nicht erforderlich zur Laufzeit Zertifikate zu erstellen
o.ä.. Ebenso ist keine Let’s Encrypt Anbindung erforderlich.
Optionale Anforderungen
## Urlaubsmodus
43. Ein Bearbeiter kann sich in einem ”Urlaubsmodus“ befinden.
44. Ist ein Bearbeiter in diesem Modus, kann ihm kein Ticket zugeteilt werden.
## Zusammenfuhren von Tickets 
45. Es soll möglich sein zwei Tickets zu einem Ticket zu verschmelzen.
46. Dabei sollen einem Ticket alle Einträge eines zweiten Tickets hinzugefugt werden.
47. Das zweite Ticket soll danach gelöscht werden.
48. Möglich soll dieser Vorgang nur sein, wenn beide Tickets den selben Bearbeiter haben.
