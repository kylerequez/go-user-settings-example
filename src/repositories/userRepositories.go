package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/kylerequez/go-user-settings-example/src/models"
)

type UserRepository struct {
	db   *pgx.Conn
	Name string
}

const (
	AUTH_ADMIN = "admin"
	AUTH_USER  = "user"
)

var AUTHORITIES = []string{
	AUTH_ADMIN,
	AUTH_USER,
}

func NewUserRepository(db *pgx.Conn, name string) *UserRepository {
	return &UserRepository{
		db:   db,
		Name: name,
	}
}

func (ur *UserRepository) GetAllUsers(ctx context.Context, id uuid.UUID) (*[]models.User, error) {
	sql := fmt.Sprintf(`
      SELECT
        u.id,
        u.name,
    		u.email,
    		u.authority,
        u.createdAt,
        u.updatedAt,
    		s.id,
    		s.theme
      FROM
        %s AS u,
    		%s AS s
			WHERE
				u.id = s.user_id
			AND 
				u.id != $1;
    	;
    `, ur.Name, "user_settings")

	_, err := ur.db.Prepare(ctx, "GetAllUsers", sql)
	if err != nil {
		return nil, err
	}

	res, err := ur.db.Query(ctx, "GetAllUsers",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	users := []models.User{}
	for res.Next() {
		user := new(models.User)

		if err := res.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Authority,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Settings.ID,
			&user.Settings.Theme,
		); err != nil {
			return nil, err
		}

		users = append(users, *user)
	}
	defer res.Close()

	return &users, nil
}

func (ur *UserRepository) CreateUser(ctx context.Context, user models.User) error {
	tx, err := ur.db.Begin(ctx)
	if err != nil {
		return err
	}

	sql := fmt.Sprintf(`
      INSERT INTO
        %s
      (
        name,
    		email,
    		password,
    		authority,
        createdAt,
        updatedAt
      )
      VALUES
      (
        $1,
        $2,
        $3,
        $4,
   			$5,
    		$6
      )
    	RETURNING id;
    `, ur.Name)

	_, err = tx.Prepare(ctx, "CreateUser", sql)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	id := new(uuid.UUID)

	if err := tx.QueryRow(ctx, "CreateUser",
		user.Name,
		user.Email,
		user.Password,
		user.Authority,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&id); err != nil {
		tx.Rollback(ctx)
		return err
	}

	sql = fmt.Sprintf(`
      INSERT INTO
    		%s
      (
    		user_id
      )
      VALUES
      (
        $1
      );
    `, "user_settings")

	_, err = tx.Prepare(ctx, "CreateUserSettings", sql)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx, "CreateUserSettings", id)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (ur *UserRepository) UpdateUser(ctx context.Context, user models.User) error {
	tx, err := ur.db.Begin(ctx)
	if err != nil {
		return err
	}

	sql := fmt.Sprintf(`
			UPDATE
				%s
			SET
				name = $1,
				email = $2,
				authority = $3,
				updatedAt = $4
			WHERE
				id = $5;
		`, ur.Name)

	_, err = tx.Prepare(ctx, "UpdateUser", sql)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	res, err := tx.Exec(ctx, "UpdateUser",
		user.Name,
		user.Email,
		user.Authority,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		tx.Rollback(ctx)
		return nil
	}

	count := res.RowsAffected()
	if count < 1 {
		tx.Rollback(ctx)
		return errors.New("no rows affected during update")
	}

	sql = fmt.Sprintf(`
			UPDATE
				%s
			SET
				theme = $1
			WHERE
				user_id = $2;
		`, "user_settings")

	_, err = tx.Prepare(ctx, "UpdateUserSettings", sql)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	res, err = tx.Exec(ctx, "UpdateUserSettings",
		user.Settings.Theme,
		user.ID,
	)
	if err != nil {
		tx.Rollback(ctx)
		return nil
	}

	return tx.Commit(ctx)
}

func (ur *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	sql := fmt.Sprintf(`
			DELETE FROM
				%s
			WHERE
				id = $1;
		`, ur.Name)

	_, err := ur.db.Prepare(ctx, "DeleteUser", sql)
	if err != nil {
		return err
	}

	res, err := ur.db.Exec(ctx, "DeleteUser",
		id,
	)
	if err != nil {
		return err
	}

	count := res.RowsAffected()
	if count < 1 {
		return errors.New("no rows affected during delete")
	}

	return nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	sql := fmt.Sprintf(`
			SELECT
        u.id,
        u.name,
    		u.email,
    		u.password,
    		u.authority,
        u.createdAt,
        u.updatedAt,
    		s.id,
    		s.theme
      FROM
        %s AS u,
    		%s AS s
			WHERE
				u.email = $1
			AND
				u.id = s.user_id
			LIMIT
				1;
		`, ur.Name, "user_settings")

	_, err := ur.db.Prepare(ctx, "GetUserByEmail", sql)
	if err != nil {
		return nil, err
	}

	res, err := ur.db.Query(ctx, "GetUserByEmail",
		email,
	)
	if err != nil {
		return nil, err
	}

	user := &models.User{}
	for res.Next() {
		if err := res.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Authority,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Settings.ID,
			&user.Settings.Theme,
		); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (ur *UserRepository) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	sql := fmt.Sprintf(`
      SELECT
        u.id,
        u.name,
    		u.email,
    		u.password,
    		u.authority,
        u.createdAt,
        u.updatedAt,
    		s.id,
    		s.theme
      FROM
        %s AS u,
    		%s AS s
			WHERE
				u.id = $1
			AND
				u.id = s.user_id
			LIMIT
				1;
		`, ur.Name, "user_settings")

	_, err := ur.db.Prepare(ctx, "GetUserById", sql)
	if err != nil {
		return nil, err
	}

	res, err := ur.db.Query(ctx, "GetUserById",
		id,
	)
	if err != nil {
		return nil, err
	}

	user := &models.User{}
	for res.Next() {
		if err := res.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Authority,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Settings.ID,
			&user.Settings.Theme,
		); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (ur *UserRepository) GetUsersByQuery(
	ctx context.Context,
	id uuid.UUID,
	keyword *string,
	filter *string,
	sort *string,
) (*[]models.User, error) {
	sql := fmt.Sprintf(`
			SELECT
				id,
				name,
				email,
				password,
				authority,
				createdAt,
				updatedAt
			FROM
				%s
		`, ur.Name)

	params := []interface{}{}
	paramIndex := 1

	if keyword != nil {
		sql += fmt.Sprintf(`
			WHERE
				id != $%d
				AND (
					name LIKE $%d
					OR
					email LIKE $%d
				)
		`, paramIndex, paramIndex+1, paramIndex+2)
		params = append(params, id)
		params = append(params, fmt.Sprintf("%%%s%%", *keyword))
		paramIndex += 3
	} else {
		sql += fmt.Sprintf(`
			WHERE
				id != $%d
		`, paramIndex)
		params = append(params, id)
		paramIndex++
	}

	if filter != nil {
		if keyword == nil {
			sql += " AND "
		} else {
			sql += " AND "
		}
		sql += fmt.Sprintf(`authority = $%d`, paramIndex)
		params = append(params, *filter)
		paramIndex++
	}

	if sort != nil {
		sql += fmt.Sprintf(` ORDER BY %s`, *sort)
	}

	log.Println("SQL: ", sql)

	_, err := ur.db.Prepare(ctx, sql, sql)
	if err != nil {
		return nil, err
	}

	res, err := ur.db.Query(ctx, sql, params...)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	users := []models.User{}
	for res.Next() {
		user := new(models.User)

		if err := res.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Authority,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, *user)
	}

	return &users, nil
}
