package use_cases

import (
	"context"
	"time"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

type CreateTenderDTO struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	ServiceType     string `json:"serviceType"`
	OrganizationID  string `json:"organizationId"`
	CreatorUsername string `json:"creatorUsername"`
}

type CreateTenderUseCase struct {
	employeeRepository                repositories.EmployeeRepository
	organizationResponsibleRepository repositories.OrganizationResponsibleRepository
	tenderRepository                  repositories.TenderRepository
}

func (uc *CreateTenderUseCase) Execute(dto CreateTenderDTO) (*domain.Tender, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	employee, err := uc.employeeRepository.Get(ctx, repositories.GetEmployeeDTO{
		Username: &dto.CreatorUsername,
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

	tender, err := domain.NewTender(dto.Name, dto.Description, dto.ServiceType, dto.OrganizationID, *orgResponsible)
	if err != nil {
		return nil, err
	}

	err = uc.tenderRepository.Save(ctx, *tender)
	if err != nil {
		return nil, err
	}

	return tender, nil
}

func NewCreateTenderUseCase(
	organizationResponsibleRepository repositories.OrganizationResponsibleRepository,
	tenderRepository repositories.TenderRepository,
	employeeRepository repositories.EmployeeRepository,
) CreateTenderUseCase {
	return CreateTenderUseCase{
		organizationResponsibleRepository: organizationResponsibleRepository,
		tenderRepository:                  tenderRepository,
		employeeRepository:                employeeRepository,
	}
}
