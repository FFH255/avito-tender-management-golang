package use_cases

import (
	"context"
	"time"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

type CreateBidUseCase struct {
	employeeRepository repositories.EmployeeRepository
	tenderRepository   repositories.TenderRepository
	bidRepository      repositories.BidRepository
}

func NewCreateBidUseCase(
	employeeRepository repositories.EmployeeRepository,
	tenderRepository repositories.TenderRepository,
	bidRepository repositories.BidRepository,
) CreateBidUseCase {
	return CreateBidUseCase{
		employeeRepository: employeeRepository,
		tenderRepository:   tenderRepository,
		bidRepository:      bidRepository,
	}
}

type CreateBidDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TenderID    string `json:"tenderId"`
	AuthorType  string `json:"authorType"`
	AuthorID    string `json:"authorId"`
}

func (uc CreateBidUseCase) Execute(dto CreateBidDTO) (*domain.Bid, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Поверка существования Employee
	authorID := domain.ID(dto.AuthorID)
	employee, err := uc.employeeRepository.Get(ctx, repositories.GetEmployeeDTO{
		ID: &authorID,
	})
	if err != nil {
		return nil, err
	}

	// Проверка существования Tender
	tenderID := domain.ID(dto.TenderID)
	tender, err := uc.tenderRepository.Get(ctx, repositories.GetTenderDTO{
		ID: tenderID,
	})
	if err != nil {
		return nil, err
	}

	// Создание Bid
	bid, err := domain.NewBid(dto.Name, dto.Description, dto.AuthorType, tender.ID, employee.ID)
	if err != nil {
		return nil, err
	}

	// Сохранение Bid
	if err := uc.bidRepository.Save(ctx, *bid); err != nil {
		return nil, err
	}

	return bid, nil
}
