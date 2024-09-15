package use_cases

import (
	"context"
	"time"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

type GetUserBidsUseCase struct {
	employeeRepository repositories.EmployeeRepository
	bidRepository      repositories.BidRepository
}

func NewGetUserBidsUseCase(
	employeeRepository repositories.EmployeeRepository,
	bidRepository repositories.BidRepository,
) GetUserBidsUseCase {
	return GetUserBidsUseCase{
		employeeRepository: employeeRepository,
		bidRepository:      bidRepository,
	}
}

type GetUserBidsDTO struct {
	Limit    *int
	Offset   *int
	Username string
}

func (uc GetUserBidsUseCase) Execute(dto GetUserBidsDTO) ([]domain.Bid, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Проверка существования Employee
	employee, err := uc.employeeRepository.Get(ctx, repositories.GetEmployeeDTO{
		Username: &dto.Username,
	})
	if err != nil {
		return nil, err
	}

	// Получение списка Bid
	limit := repositories.NewLimit(dto.Limit)
	offset := repositories.NewOffset(dto.Offset)
	bidList, err := uc.bidRepository.GetList(ctx, repositories.GetBidListDTO{
		AuthorID: &employee.ID,
		Limit:    &limit,
		Offset:   &offset,
	})
	if err != nil {
		return nil, err
	}

	return bidList, nil
}
