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
	ID          ID
	Name        TenderName
	Description TenderDescription
	ServiceType TenderServiceType
	Version     TenderVersion
	CreatedAt   time.Time
}

func NewTenderSnapshot(name TenderName, description TenderDescription, serviceType TenderServiceType, v TenderVersion) TenderSnapshot {
	return TenderSnapshot{
		ID:          NewID(),
		Name:        name,
		Description: description,
		ServiceType: serviceType,
		Version:     v,
		CreatedAt:   time.Now(),
	}
}

type Tender struct {
	ID             ID                `json:"id"`
	Name           TenderName        `json:"name"`
	Description    TenderDescription `json:"description"`
	OrganizationID ID                `json:"-"`
	Status         TenderStatus      `json:"status"`
	ServiceType    TenderServiceType `json:"serviceType"`
	Version        TenderVersion     `json:"version"`
	Snapshots      []TenderSnapshot  `json:"-"`
	CreatedAt      time.Time         `json:"createdAt"`
}

func (t *Tender) Rollback(executor OrganizationResponsible, version int) error {

	if executor.OrganizationID != t.OrganizationID {
		return errors.Wrap(ErrNoPermission, "Organization responsible has no access to rollback it")
	}

	v, err := NewTenderVersion(version)

	if err != nil {
		return err
	}

	i := slices.IndexFunc(t.Snapshots, func(s TenderSnapshot) bool {
		return s.Version == v
	})

	if i == -1 {
		return errors.Wrapf(ErrNotFound, "Tender snapshot %d not found", v)
	}

	snapshot := t.Snapshots[i]

	t.Snapshots = append(t.Snapshots, NewTenderSnapshot(t.Name, t.Description, t.ServiceType, t.Version))
	t.Name = snapshot.Name
	t.Description = snapshot.Description
	t.ServiceType = snapshot.ServiceType
	t.Version++

	return nil
}

func (t *Tender) Edit(executor OrganizationResponsible, name, description, serviceType *string) error {

	if executor.OrganizationID != t.OrganizationID {
		return errors.Wrap(ErrNoPermission, "Organization responsible has no access to edit Tender")
	}

	t.Snapshots = append(t.Snapshots, NewTenderSnapshot(t.Name, t.Description, t.ServiceType, t.Version))

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

	return nil
}

func (t *Tender) ChangeStatus(executor OrganizationResponsible, status string) error {

	if executor.OrganizationID != t.OrganizationID {
		return errors.Wrap(ErrNoPermission, "Organization responsible has no access to change status of Tender")
	}

	switch status {
	case string(TenderClosedStatus):
		t.Status = TenderClosedStatus
		return nil

	case string(TenderPublishedStatus):
		t.Status = TenderPublishedStatus
		return nil

	case string(TenderCreatedStatus):
		t.Status = TenderCreatedStatus
		return nil
	}

	return errors.Wrap(ErrValidation, "Tender status is invalid")
}

func NewTender(name, description, serviceType, organizationID string, executor OrganizationResponsible) (*Tender, error) {

	orgID := ID(organizationID)

	if executor.OrganizationID != orgID {
		return nil, errors.Wrap(ErrNoPermission, "New Tender responsible has no access to create Tender because of different organization")
	}

	n, err := NewTenderName(name)

	if err != nil {
		return nil, err
	}

	desc, err := NewTenderDescription(description)

	if err != nil {
		return nil, err
	}

	t, err := NewTenderServiceType(serviceType)

	if err != nil {
		return nil, err
	}

	id := NewID()

	return &Tender{
		ID:             id,
		Name:           n,
		Description:    desc,
		Status:         TenderCreatedStatus,
		ServiceType:    t,
		OrganizationID: orgID,
		Version:        TenderVersion(1),
		CreatedAt:      time.Now(),
		Snapshots:      []TenderSnapshot{},
	}, nil
}
