package mail

import "TicketSystem/XML_IO"

func CreateTicketFromMail(mail string, reference string, message string) (XML_IO.Ticket, error) {
	tickets := XML_IO.GetTicketsByEditor("../data/tickets/ticket", "definitions.xml", mail)
	for _, actTicket := range tickets {
		if actTicket.Reference == reference {
			//TODO: Ticketstatus ggf aendern
			if actTicket.Status != 1 {
				XML_IO.ChangeStatus("../data/tickets/ticket", actTicket.Id, 1)
			}
			return XML_IO.AddMessage("../data/tickets/ticket", actTicket, mail, message)
		}
		//TODO: schauen ob Ticket bereits vorhanden ist
	}
	return CreateTicketFromMail(mail, reference, message)
}
