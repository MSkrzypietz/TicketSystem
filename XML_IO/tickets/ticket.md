# Aufbau der XML-Datei eines Tickets

**Ticket** Übergeordnetes Tag und enthält alle weiteren Tags zu einem Ticket

**ID** Enthält die ID des dazugehörigen Tickets
**Client** Enthält die E-Mail des Mitarbeiters, der das Ticket erstellt hat
**Reference** Betreff des Tickets
**Status** Aktueller Status des Tickets: 0 = offen, 1 = in Bearbeitung, 2 = geschlossen
**Editor** Aktueller Bearbeiter des Tickets

**MessageList** Liste mit all den Nachrichten, die dem Tickets hinzugefügt wurden

**Message** Definiert eine Nachricht
**CreationDate** Erstellungsdatum der Nachricht
**Actor** E-Mail oder Mitarbeiternummer der Person, die die Nachricht hinzugefügt hat
**Text** Text der Nachricht