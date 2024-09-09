package domain

import "github.com/pkg/errors"

type OrganizationType string

const (
	IEOrganizationType  OrganizationType = "IE"
	LLCOrganizationType OrganizationType = "LLC"
	JSCOrganizationType OrganizationType = "JSC"
)

func NewOrganizationType(str string) (OrganizationType, error) {
	switch str {
	case string(IEOrganizationType), string(LLCOrganizationType), string(JSCOrganizationType):
		return OrganizationType(str), nil
	default:
		return "", errors.Wrap(ErrValidation, "Invalid organization type")
	}
}

type Organization struct {
	ID          ID
	Name        string
	Description string
	Type        OrganizationType
}

func NewOrganization(name, description, organizationType string) (Organization, error) {

	if len(name) == 0 {
		return Organization{}, errors.Wrap(ErrValidation, "Organization name must not be empty")
	}

	if len(name) > 100 {
		return Organization{}, errors.Wrap(ErrValidation, "Organization name must not exceed 100 characters")
	}

	orgType, err := NewOrganizationType(organizationType)

	if err != nil {
		return Organization{}, err
	}

	id := NewID()

	return Organization{
		ID:          id,
		Name:        name,
		Description: description,
		Type:        orgType,
	}, nil
}
