package employee_repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

func (r EmployeeRepository) Get(ctx context.Context, dto repositories.GetEmployeeDTO) (*domain.Employee, error) {
	query := `SELECT id, username, first_name, last_name FROM employee WHERE True`

	args := make([]interface{}, 0)
	i := 1

	if dto.ID != nil {
		args = append(args, *dto.ID)
		query += fmt.Sprintf(" AND id = $%d", i)
		i++
	}

	if dto.Username != nil {
		args = append(args, *dto.Username)
		query += fmt.Sprintf(" AND username = $%d", i)
		i++
	}

	row := r.client.QueryRow(ctx, query, args...)

	var employee domain.Employee

	err := row.Scan(&employee.ID, &employee.Username, &employee.FirstName, &employee.LastName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(domain.ErrUserNotFound, "employee not found")
		}
		return nil, err
	}

	return &employee, nil
}
