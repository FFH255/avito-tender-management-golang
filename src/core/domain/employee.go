package domain

import "github.com/pkg/errors"

type Employee struct {
	ID        ID
	Username  string
	FirstName *string
	LastName  *string
}

func NewEmployee(username string, firstName, lastName *string) (Employee, error) {

	if len(username) > 50 {
		return Employee{}, errors.Wrap(ErrValidation, "Username must not exceed 50 characters")
	}

	if len(*firstName) > 50 {
		return Employee{}, errors.Wrap(ErrValidation, "First name must not exceed 50 characters")
	}

	if len(*lastName) > 50 {
		return Employee{}, errors.Wrap(ErrValidation, "Last name must not exceed 50 characters")
	}

	id := NewID()

	return Employee{
		ID:        id,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
	}, nil
}
