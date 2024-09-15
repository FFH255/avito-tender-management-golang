package use_cases

import (
	"context"
	"time"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

type GetUserTendersUseCase struct {
	employeeRepository       repositories.EmployeeRepository
	orgResponsibleRepository repositories.OrganizationResponsibleRepository
	tenderRepository         repositories.TenderRepository
}

type GetUserTendersDTO struct {
	Limit    *int   `json:"limit"`
	Offset   *int   `json:"offset"`
	Username string `json:"username"`
}

func (uc GetUserTendersUseCase) Execute(dto GetUserTendersDTO) ([]domain.Tender, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	employee, err := uc.employeeRepository.Get(ctx, repositories.GetEmployeeDTO{
		Username: &dto.Username,
	})
	if err != nil {
		return nil, err
	}

	orgResponsible, err := uc.orgResponsibleRepository.Get(ctx, repositories.GetOrganizationResponsibleDTO{
		EmployeeID: employee.ID,
	})
	if err != nil {
		return nil, err
	}

	limit := repositories.NewLimit(dto.Limit)
	offset := repositories.NewOffset(dto.Offset)
	tenders, err := uc.tenderRepository.GetList(ctx, repositories.GetTendersListDTO{
		OrganizationID: &orgResponsible.OrganizationID,
		Offset:         &offset,
		Limit:          &limit,
	})
	if err != nil {
		return nil, err
	}

	return tenders, nil
}

func NewGetUserTendersUseCase(
	employeeRepository repositories.EmployeeRepository,
	orgResponsibleRepository repositories.OrganizationResponsibleRepository,
	tendersRepository repositories.TenderRepository,
) GetUserTendersUseCase {
	return GetUserTendersUseCase{
		employeeRepository:       employeeRepository,
		orgResponsibleRepository: orgResponsibleRepository,
		tenderRepository:         tendersRepository,
	}
}
