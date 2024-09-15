package use_cases

import (
	"context"
	"time"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

type EditBidUseCase struct {
	employeeRepository repositories.EmployeeRepository
	bidRepository      repositories.BidRepository
}

func NewEditBidUseCase(
	employeeRepository repositories.EmployeeRepository,
	bidRepository repositories.BidRepository,
) EditBidUseCase {
	return EditBidUseCase{
		employeeRepository: employeeRepository,
		bidRepository:      bidRepository,
	}
}

type EditBidDTO struct {
	BidID       string
	Username    string
	Name        *string
	Description *string
}

func (uc EditBidUseCase) Execute(dto EditBidDTO) (*domain.Bid, error) {
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

	if err := bid.Edit(dto.Name, dto.Description); err != nil {
		return nil, err
	}

	if err := uc.bidRepository.Save(ctx, *bid); err != nil {
		return nil, err
	}

	return bid, nil
}
