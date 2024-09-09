package use_cases

import (
	"context"
	"github.com/FFH255/avito-tender-management-golang/internal/core/services/repositories"
)

type EditTenderUseCaseDTO struct {
	TenderID    string
	Username    string
	name        *string
	description *string
	serviceType *string
}

type EditTenderUseCase struct {
	organizationResponsibleRepository repositories.OrganizationResponsibleRepository
	tenderRepository                  repositories.TenderRepository
}

func (uc EditTenderUseCase) Execute(dto EditTenderUseCaseDTO) error {

	tender, err := uc.tenderRepository.GetOne(context.TODO(), repositories.GetOneTenderDTO{
		ID: dto.TenderID,
	})

	if err != nil {
		return err
	}

	orgResponsible, err := uc.organizationResponsibleRepository.GetOne(context.TODO(), repositories.GetOneOrganizationResponsibleDTO{
		Username: dto.Username,
	})

	err = tender.Edit(orgResponsible, dto.name, dto.description, dto.serviceType)

	return err
}

func NewEditTenderUseCase(orr repositories.OrganizationResponsibleRepository, tr repositories.TenderRepository) EditTenderUseCase {
	return EditTenderUseCase{
		organizationResponsibleRepository: orr,
		tenderRepository:                  tr,
	}
}
