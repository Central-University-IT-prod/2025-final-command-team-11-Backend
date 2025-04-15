package user

import (
	"fmt"
	"strings"
	"time"

	"github.com/nikitaSstepanov/coffee-id/internal/entity"
	types "github.com/nikitaSstepanov/coffee-id/internal/entity/type"
)

func idQuery(id string) string {
	return fmt.Sprintf(
		`
			SELECT * FROM %s 
			WHERE id = '%s';
		`, usersTable, id,
	)
}

func emailQuery(email string) string {
	return fmt.Sprintf(
		`
			SELECT * FROM %s 
			WHERE email = '%s';
		`, usersTable, email,
	)
}

func createQuery(user *entity.User) string {
	if user.OAuth != "" {
		return fmt.Sprintf(
			`
				INSERT INTO %s 
					(email, name, password, birthday, oauth) 
				VALUES 
					('%s', '%s', '%s', '%s', '%s') 
				RETURNING id, role;
			`,
			usersTable, user.Email, user.Name, user.Password, user.Birthday.Format(time.DateOnly), user.OAuth,
		)
	}

	return fmt.Sprintf(
		`
			INSERT INTO %s 
				(email, name, password, birthday) 
			VALUES 
				('%s', '%s', '%s', '%s') 
			RETURNING id, role;
		`,
		usersTable, user.Email, user.Name, user.Password, user.Birthday.Format(time.DateOnly),
	)
}

func updateQuery(user *entity.User) string {
	toUpd := setupValues(user)

	return fmt.Sprintf(
		`
			UPDATE %s 
			SET %s 
			WHERE id = '%s';
		`, usersTable, toUpd, user.Id,
	)
}

func roleQuery(user *entity.User, roles types.Role) string {
	return fmt.Sprintf(
		`
			UPDATE %s 
			SET role = '%s'
			WHERE id = '%s';
		`, usersTable, roles, user.Id,
	)
}

func verifyQuery(verified bool, id string) string {
	return fmt.Sprintf(
		`
			UPDATE %s 
			SET verified = %t
			WHERE id = '%s';
		`, usersTable, verified, id,
	)
}

func deleteQuery() string {
	return fmt.Sprintf(
		`
			DELETE FROM %s 
			WHERE id = $1;
		`, usersTable,
	)
}

func setupValues(user *entity.User) string {
	toUpd := make([]string, 0)

	if user.Email != "" {
		toUpd = append(toUpd, fmt.Sprintf("email = '%s'", user.Email))
	}

	if user.Name != "" {
		toUpd = append(toUpd, fmt.Sprintf("name = '%s'", user.Name))
	}

	if user.Password != "" {
		toUpd = append(toUpd, fmt.Sprintf("password = '%s'", user.Password))
	}

	if !user.Birthday.IsZero() {
		toUpd = append(toUpd, fmt.Sprintf("birthday = '%s'", user.Birthday.Format(time.DateOnly)))
	}

	query := strings.Join(toUpd, ", ")

	return query
}

func redisKey(id string) string {
	return fmt.Sprintf("users:%s", id)
}
