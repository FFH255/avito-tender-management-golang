package use_cases

import (
	"context"
	"time"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

type RollbackBidUseCase struct {
	employeeRepository repositories.EmployeeRepository
	bidRepository      repositories.BidRepository
}

func NewRollbackBidUseCase(
	employeeRepository repositories.EmployeeRepository,
	bidRepository repositories.BidRepository,
) RollbackBidUseCase {
	return RollbackBidUseCase{
		employeeRepository: employeeRepository,
		bidRepository:      bidRepository,
	}
}

type RollbackBidDTO struct {
	BidID    string
	Version  int
	Username string
}

func (uc RollbackBidUseCase) Execute(dto RollbackBidDTO) (*domain.Bid, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	employee, err := uc.employeeRepository.Get(ctx, repositories.GetEmployeeDTO{
		Username: &dto.Username,
	})
	if err != nil {
		return nil, err
	}

	bidID := domain.ID(dto.BidID)
	bid, err := uc.bidRepository.Get(ctx, repositories.GetBidDTO{
		ID:       bidID,
		AuthorID: &employee.ID,
	})
	if err != nil {
		return nil, err
	}

	if err := bid.Rollback(dto.Version); err != nil {
		return nil, err
	}

	if err := uc.bidRepository.Save(ctx, *bid); err != nil {
		return nil, err
	}

	return bid, nil
}
