package domain

import (
	"github.com/pkg/errors"
	"slices"
)

// DecisionStatus Статус Решения
type DecisionStatus string

const (
	DecisionApprovedStatus DecisionStatus = "Approved"
	DecisionRejectedStatus DecisionStatus = "Rejected"
)

func NewDecisionStatus(status string) (DecisionStatus, error) {
	switch status {
	case string(DecisionApprovedStatus):
		return DecisionApprovedStatus, nil

	case string(DecisionRejectedStatus):
		return DecisionRejectedStatus, nil

	default:
		return "", errors.Wrapf(ErrValidation, "invalid decision status - '%s'", status)
	}
}

// Decision Решение по предложению
type Decision struct {
	ID       ID
	AuthorID ID
	BidID    ID
	TenderID ID
	Status   DecisionStatus
}

func NewDecision(bidDecisions []Decision, authorID, bidID, tenderID ID, status string) (*Decision, error) {
	i := slices.IndexFunc(bidDecisions, func(d Decision) bool {
		return d.AuthorID == authorID
	})
	if i != -1 {
		return nil, errors.Wrapf(ErrAlreadyExist, "author already send decision on '%s'", bidID)
	}

	id := NewID()

	s, err := NewDecisionStatus(status)
	if err != nil {
		return nil, err
	}

	return &Decision{
		ID:       id,
		AuthorID: authorID,
		BidID:    bidID,
		TenderID: tenderID,
		Status:   s,
	}, nil
}

func MakeFinalDecision(
	quorumSize int,
	decisionsCount int,
	incomingDecision Decision,
	tenderOwner OrganizationResponsible,
	tender *Tender,
	bid *Bid,
) {
	if incomingDecision.Status == DecisionRejectedStatus {
		bid.Status = BidCanceledStatus
		return
	}

	if quorumSize <= decisionsCount+1 {
		_ = tender.ChangeStatus(tenderOwner, string(TenderClosedStatus))
	}
}
