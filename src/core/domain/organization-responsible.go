package domain

type OrganizationResponsible struct {
	ID             ID
	OrganizationID ID
	UserID         ID
}

func NewOrganizationResponsible(organizationID ID, userID ID) OrganizationResponsible {

	id := NewID()

	return OrganizationResponsible{
		ID:             id,
		OrganizationID: organizationID,
		UserID:         userID,
	}
}
