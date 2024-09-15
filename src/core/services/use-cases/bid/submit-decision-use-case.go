package use_cases

import (
	"context"
	"time"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

type SubmitDecisionUseCase struct {
	employeeRepository repositories.EmployeeRepository
	orgRespRepository  repositories.OrganizationResponsibleRepository
	bidRepository      repositories.BidRepository
	tenderRepository   repositories.TenderRepository
	decisionRepository repositories.DecisionRepository
}

func NewSubmitDecisionUseCase(
	employeeRepository repositories.EmployeeRepository,
	orgRespRepository repositories.OrganizationResponsibleRepository,
	bidRepository repositories.BidRepository,
	tenderRepository repositories.TenderRepository,
	decisionRepository repositories.DecisionRepository,
) SubmitDecisionUseCase {
	return SubmitDecisionUseCase{
		employeeRepository: employeeRepository,
		orgRespRepository:  orgRespRepository,
		bidRepository:      bidRepository,
		tenderRepository:   tenderRepository,
		decisionRepository: decisionRepository,
	}
}

type SubmitDecisionDTO struct {
	BidID    string
	Decision string
	Username string
}

func (uc SubmitDecisionUseCase) Execute(dto SubmitDecisionDTO) (*domain.Bid, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tenderOwnerEmployee, err := uc.employeeRepository.Get(ctx, repositories.GetEmployeeDTO{
		Username: &dto.Username,
	})
	if err != nil {
		return nil, err
	}

	tenderOwnerOrgResp, err := uc.orgRespRepository.Get(ctx, repositories.GetOrganizationResponsibleDTO{
		EmployeeID: tenderOwnerEmployee.ID,
	})
	if err != nil {
		return nil, err
	}

	bidID := domain.ID(dto.BidID)
	bid, err := uc.bidRepository.Get(ctx, repositories.GetBidDTO{
		ID: bidID,
	})
	if err != nil {
		return nil, err
	}

	tender, err := uc.tenderRepository.Get(ctx, repositories.GetTenderDTO{
		ID: bid.TenderID,
	})
	if err != nil {
		return nil, err
	}

	decisions, err := uc.decisionRepository.GetList(ctx, repositories.GetDecisionListDTO{
		TenderID: &tender.ID,
	})
	if err != nil {
		return nil, err
	}

	tenderQuorum, err := uc.orgRespRepository.GetList(ctx, repositories.GetOrganizationResponsiblesListDTO{
		OrganizationID: &tender.OrganizationID,
	})
	if err != nil {
		return nil, err
	}

	decision, err := domain.NewDecision(decisions, tenderOwnerEmployee.ID, bidID, tender.ID, dto.Decision)
	if err != nil {
		return nil, err
	}

	domain.MakeFinalDecision(len(tenderQuorum), len(decisions), *decision, *tenderOwnerOrgResp, tender, bid)

	if err := uc.tenderRepository.Save(ctx, *tender); err != nil {
		return nil, err
	}

	if err := uc.bidRepository.Save(ctx, *bid); err != nil {
		return nil, err
	}

	if err := uc.decisionRepository.Save(ctx, *decision); err != nil {
		return nil, err
	}

	return bid, nil
}
