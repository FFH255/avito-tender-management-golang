package use_cases

import (
	"context"
	"time"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

type EditTenderUseCaseDTO struct {
	TenderID    string
	Username    string
	Name        *string
	Description *string
	ServiceType *string
}

type EditTenderUseCase struct {
	employeeRepository                repositories.EmployeeRepository
	organizationResponsibleRepository repositories.OrganizationResponsibleRepository
	tenderRepository                  repositories.TenderRepository
}

func (uc EditTenderUseCase) Execute(dto EditTenderUseCaseDTO) (*domain.Tender, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	employee, err := uc.employeeRepository.Get(ctx, repositories.GetEmployeeDTO{
		Username: &dto.Username,
	})
	if err != nil {
		return nil, err
	}

	orgResponsible, err := uc.organizationResponsibleRepository.Get(ctx, repositories.GetOrganizationResponsibleDTO{
		EmployeeID: employee.ID,
	})
	if err != nil {
		return nil, err
	}

	tenderID := domain.ID(dto.TenderID)
	tender, err := uc.tenderRepository.Get(ctx, repositories.GetTenderDTO{
		ID: tenderID,
	})
	if err != nil {
		return nil, err
	}

	if err = tender.Edit(*orgResponsible, dto.Name, dto.Description, dto.ServiceType); err != nil {
		return nil, err
	}

	if err = uc.tenderRepository.Save(ctx, *tender); err != nil {
		return nil, err
	}

	return tender, nil
}

func NewEditTenderUseCase(
	employeeRepository repositories.EmployeeRepository,
	organizationResponsibleRepository repositories.OrganizationResponsibleRepository,
	tenderRepository repositories.TenderRepository,
) EditTenderUseCase {
	return EditTenderUseCase{
		organizationResponsibleRepository: organizationResponsibleRepository,
		tenderRepository:                  tenderRepository,
		employeeRepository:                employeeRepository,
	}
}
