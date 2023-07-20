package main

import "errors"

type ticketRequest struct {
	toggleRollStatus    OPS_ROLL_STATUS
	toggleDeployFlavour OPS_JIRA_DEPLOY_FLAVOUR
	toggleProject       OPS_JIRA_PROJECT
	ticket              opsTicket
	viewID              string
}

func (t *ticketRequest) SetViewID(view string) {
	t.viewID = view
}

func (t *ticketRequest) New() *ticketRequest {
	return &ticketRequest{}
}

// ticket queue manager
type ticketQueue struct {
	tickets []ticketRequest
}

func (q *ticketQueue) Add(newTicket ticketRequest) {
	q.tickets = append(q.tickets, newTicket)
}

func (q *ticketQueue) GetTicketRequestByViewID(viewID string) (ticketRequest, error) {
	for _, ticket := range q.tickets {
		if ticket.viewID == viewID {
			return ticket, nil
		}
	}
	return ticketRequest{}, errors.New("ticket not found (" + viewID + ")")
}

func (q *ticketQueue) SetTicketRequestByViewID(viewID string, ticket ticketRequest) error {
	error := errors.New("no ticket to update, viewid not found (" + viewID + ")")
	for i, existingTickets := range q.tickets {

		if existingTickets.viewID == viewID {
			q.tickets[i] = ticket
			error = nil
		}

	}
	return error
}

func (q *ticketQueue) RemoveTicketRequestByViewID(viewID string) error {
	error := errors.New("no ticket to remove, viewid not found (" + viewID + ")")
	for i, ticket := range q.tickets {

		if ticket.viewID == viewID {
			q.tickets[i] = q.tickets[len(q.tickets)-1]
			q.tickets = q.tickets[:len(q.tickets)-1]
			error = nil
		}

	}
	return error
}
