package use_cases

import (
	"context"
	"github.com/pkg/errors"
	"time"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

type GetBidStatusUseCase struct {
	employeeRepository repositories.EmployeeRepository
	bidRepository      repositories.BidRepository
}

func NewGetBidStatusUseCase(
	employeeRepository repositories.EmployeeRepository,
	bidRepository repositories.BidRepository,
) GetBidStatusUseCase {
	return GetBidStatusUseCase{
		employeeRepository: employeeRepository,
		bidRepository:      bidRepository,
	}
}

type GetBidStatusDTO struct {
	BidID    string
	Username string
}

func (uc GetBidStatusUseCase) Execute(dto GetBidStatusDTO) (*domain.BidStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Проверка существования Employee
	employee, err := uc.employeeRepository.Get(ctx, repositories.GetEmployeeDTO{
		Username: &dto.Username,
	})
	if err != nil {
		return nil, err
	}

	// Получение Bid
	bidID := domain.ID(dto.BidID)
	bid, err := uc.bidRepository.Get(ctx, repositories.GetBidDTO{
		ID: bidID,
	})
	if err != nil {
		return nil, err
	}

	// Проверка прав Employee
	if bid.AuthorID != employee.ID {
		return nil, errors.Wrap(domain.ErrNoPermission, "employee is not author of bid")
	}

	return &bid.Status, nil
}
