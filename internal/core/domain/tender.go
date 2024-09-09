package domain

import (
	"github.com/pkg/errors"
	"slices"
	"time"
)

type TenderStatus string

const (
	TenderCreatedStatus   TenderStatus = "CREATED"
	TenderPublishedStatus TenderStatus = "PUBLISHED"
	TenderClosedStatus    TenderStatus = "CLOSED"
)

type TenderServiceType string

const (
	TenderConstructionServiceType TenderServiceType = "Construction"
	TenderDeliveryServiceType     TenderServiceType = "Delivery"
	TenderManufactureServiceType  TenderServiceType = "Manufacture"
)

func NewTenderServiceType(t string) (TenderServiceType, error) {
	switch t {
	case string(TenderConstructionServiceType), string(TenderDeliveryServiceType), string(TenderManufactureServiceType):
		return TenderServiceType(t), nil
	}
	return "", errors.Wrap(ErrValidation, "Invalid tender service type")
}

type TenderVersion int

func NewTenderVersion(v int) (TenderVersion, error) {

	if v <= 0 {
		return TenderVersion(1), errors.Wrap(ErrValidation, "Tender version must be positive")
	}

	return TenderVersion(v), nil
}

type TenderName string

func NewTenderName(str string) (TenderName, error) {

	if len(str) > 100 {
		return "", errors.Wrap(ErrValidation, "Tender name cannot exceed 100 characters")
	}

	return TenderName(str), nil
}

type TenderDescription string

func NewTenderDescription(str string) (TenderDescription, error) {

	if len(str) > 500 {
		return "", errors.Wrap(ErrValidation, "Tender description cannot exceed 500 characters")
	}

	return TenderDescription(str), nil
}

type TenderSnapshot struct {
	Name        TenderName
	Description TenderDescription
	ServiceType TenderServiceType
	Version     TenderVersion
	CreatedAt   time.Time
}

func NewTenderSnapshot(name TenderName, description TenderDescription, serviceType TenderServiceType, v TenderVersion) TenderSnapshot {
	return TenderSnapshot{
		Name:        name,
		Description: description,
		ServiceType: serviceType,
		Version:     v,
		CreatedAt:   time.Now(),
	}
}

type Tender struct {
	ID             ID
	Name           TenderName
	Description    TenderDescription
	OrganizationID ID
	Status         TenderStatus
	ServiceType    TenderServiceType
	Version        TenderVersion
	Snapshots      []TenderSnapshot
}

func (t *Tender) Publish(executor OrganizationResponsible) error {

	if executor.OrganizationID != t.OrganizationID {
		return errors.New("Organization responsible has no access to tender to publish it")
	}

	t.Status = TenderPublishedStatus

	return nil
}

func (t *Tender) Close(executor OrganizationResponsible) error {

	if executor.OrganizationID != t.OrganizationID {
		return errors.New("Organization responsible has no access on tender to close it")
	}

	t.Status = TenderClosedStatus

	return nil
}

func (t *Tender) Rollback(executor OrganizationResponsible, version int) error {

	if executor.OrganizationID != t.OrganizationID {
		return errors.New("Organization responsible has no access to rollback it")
	}

	v, err := NewTenderVersion(version)

	if err != nil {
		return err
	}

	i := slices.IndexFunc(t.Snapshots, func(s TenderSnapshot) bool {
		return s.Version == v
	})

	if i == -1 {
		return errors.New("There is snapshot with such version")
	}

	snapshot := t.Snapshots[i]

	t.Name = snapshot.Name
	t.Description = snapshot.Description
	t.ServiceType = snapshot.ServiceType
	t.Version++

	t.Snapshots = append(t.Snapshots, NewTenderSnapshot(t.Name, t.Description, t.ServiceType, t.Version))

	return nil
}

func (t *Tender) Edit(executor OrganizationResponsible, name, description, serviceType *string) error {

	if executor.OrganizationID != t.OrganizationID {
		return errors.New("Organization responsible has no access to edit Tender")
	}

	if name != nil {
		n, err := NewTenderName(*name)

		if err != nil {
			return err
		}
		t.Name = n
	}

	if description != nil {
		d, err := NewTenderDescription(*description)

		if err != nil {
			return err
		}
		t.Description = d
	}

	if serviceType != nil {
		sType, err := NewTenderServiceType(*serviceType)

		if err != nil {
			return err
		}
		t.ServiceType = sType
	}

	t.Version++
	t.Snapshots = append(t.Snapshots, NewTenderSnapshot(t.Name, t.Description, t.ServiceType, t.Version))

	return nil
}

func NewTender(name, description, serviceType string, organizationID ID) (Tender, error) {

	n, err := NewTenderName(name)

	if err != nil {
		return Tender{}, err
	}

	desc, err := NewTenderDescription(description)

	if err != nil {
		return Tender{}, err
	}

	t, err := NewTenderServiceType(serviceType)

	if err != nil {
		return Tender{}, err
	}

	id := NewID()

	return Tender{
		ID:             id,
		Name:           n,
		Description:    desc,
		Status:         TenderCreatedStatus,
		ServiceType:    t,
		OrganizationID: organizationID,
		Version:        TenderVersion(1),
		Snapshots:      []TenderSnapshot{NewTenderSnapshot(n, desc, t, TenderVersion(1))},
	}, nil
}
