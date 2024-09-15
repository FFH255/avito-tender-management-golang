package use_cases

import (
	"context"
	"github.com/pkg/errors"
	"time"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

type ChangeBidStatusUseCase struct {
	employeeRepository repositories.EmployeeRepository
	bidRepository      repositories.BidRepository
}

func NewChangeBidStatusUseCase(
	employeeRepository repositories.EmployeeRepository,
	bidRepository repositories.BidRepository,
) ChangeBidStatusUseCase {
	return ChangeBidStatusUseCase{
		employeeRepository: employeeRepository,
		bidRepository:      bidRepository,
	}
}

type ChangeBidStatusDTO struct {
	BidID    string
	Status   string
	Username string
}

func (uc ChangeBidStatusUseCase) Execute(dto ChangeBidStatusDTO) (*domain.Bid, error) {
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

	// Изменение статуса Bid
	if err := bid.ChangeStatus(dto.Status); err != nil {
		return nil, err
	}

	// Сохранение Bid
	if err := uc.bidRepository.Save(ctx, *bid); err != nil {
		return nil, err
	}

	return bid, nil
}
