package use_cases

import (
	"context"
	"github.com/FFH255/avito-tender-management-golang/internal/core/domain"
	"github.com/FFH255/avito-tender-management-golang/internal/core/services/repositories"
)

type CreateTenderUseCaseDTO struct {
	Name            string
	Description     string
	ServiceType     string
	Status          string
	OrganizationID  string
	CreatorUsername string
}

type CreateTenderUseCase struct {
	organizationResponsibleRepository repositories.OrganizationResponsibleRepository
	tenderRepository                  repositories.TenderRepository
}

func (uc *CreateTenderUseCase) Execute(dto CreateTenderUseCaseDTO) error {

	orgResponsible, err := uc.organizationResponsibleRepository.GetOne(
		context.TODO(),
		repositories.GetOneOrganizationResponsibleDTO{
			Username: dto.CreatorUsername,
		})

	if err != nil {
		return err
	}

	tender, err := domain.NewTender(dto.Name, dto.Description, dto.ServiceType, orgResponsible.OrganizationID)

	if err != nil {
		return err
	}

	err = uc.tenderRepository.Save(context.TODO(), tender)

	return err
}

func NewCreateTenderUseCase(orr repositories.OrganizationResponsibleRepository, tr repositories.TenderRepository) CreateTenderUseCase {
	return CreateTenderUseCase{
		organizationResponsibleRepository: orr,
		tenderRepository:                  tr,
	}
}
