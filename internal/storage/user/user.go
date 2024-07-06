package user

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/ardihikaru/graphql-example-part-1/internal/service/user/dto"

	mySqlx "github.com/ardihikaru/graphql-example-part-1/pkg/mysqldb"
)

type Store struct {
	*mySqlx.Storage
}

type User struct {
	dto.User
	Passwd string `db:"pass_hash"`
}

const (
	table = "user"
)

// ToService converts the userDoc struct into User struct from the service
func (u *User) ToService() *dto.User {
	acc := &dto.User{
		UserID:   u.UserID,
		UserNm:   u.UserNm,
		IsAdmin:  u.IsAdmin,
		StatusCd: u.StatusCd,
	}

	return acc
}

// InsertUser inserts a new user
func (store *Store) InsertUser(username, hashedPassword, status string, isAdmin int, setCreateUserID int64) (*dto.User, error) {
	query := fmt.Sprintf(`
    	INSERT INTO %s (user_nm, pass_hash, is_admin, status_cd, created_user_id, created_dttm) 
    	VALUES (:user_nm, :pass_hash, :is_admin, :status_cd, :created_user_id, NOW())
    `, table)

	args := map[string]interface{}{
		"user_nm":         username,
		"pass_hash":       hashedPassword,
		"is_admin":        isAdmin,
		"status_cd":       status,
		"created_user_id": setCreateUserID,
	}

	qArgs := mySqlx.QueryArgs{
		Query: query,
		Args:  args,
	}

	rslt, err := store.NamedExec(qArgs)
	if err != nil {
		return nil, err
	}

	lastInsertedId, _ := (*rslt).(sql.Result).LastInsertId()

	// builds results
	userData := &dto.User{
		UserID:   strconv.FormatInt(lastInsertedId, 10),
		UserNm:   username,
		IsAdmin:  isAdmin,
		StatusCd: status,
	}

	return userData, nil
}

// GetUserById fetches user data by uid
func (store *Store) GetUserById(userId int64) (*dto.User, error) {
	var err error
	var row *sql.Row
	record := User{}

	query := fmt.Sprintf(`
		SELECT
		    user_id, user_nm, is_admin, status_cd, pass_hash
		FROM %s
		WHERE user_id = '%v'
	`, table, userId)
	row = store.Db.QueryRow(query)

	err = row.Scan(&record.UserID, &record.UserNm, &record.IsAdmin, &record.StatusCd, &record.Passwd)

	if err != nil && err == sql.ErrNoRows {
		store.Log.Debug(fmt.Sprintf("record is not found in the database"))
		return nil, sql.ErrNoRows
	}
	if err != nil {
		store.Log.Error(fmt.Sprintf("got error in query: %s", err.Error()))
		return nil, fmt.Errorf("QUERY_FAILED")
	}

	return record.ToService(), nil
}

// GetUsers fetches list of user
func (store *Store) GetUsers(userIdStr, statusCd string) ([]*dto.User, error) {
	var err error

	query := fmt.Sprintf(`
		SELECT user_id, user_nm, is_admin, status_cd
		FROM %s
		WHERE true
	`, table)

	if userIdStr != "" {
		query = fmt.Sprintf("%s AND user_id='%s'", query, userIdStr)
	}
	if statusCd != "" {
		query = fmt.Sprintf("%s AND status_cd='%s'", query, statusCd)
	}

	var userList []*dto.User
	rows, err := store.Queryx(query, nil)
	if err != nil && err == sql.ErrNoRows {
		store.Log.Debug(fmt.Sprintf("record is not found in the database"))
		return nil, fmt.Errorf("DATA_NOT_FOUND")
	}
	if err != nil {
		return nil, fmt.Errorf("QUERY_FAILED")
	}

	// starts to assign value to the designated struct
	for rows.Next() {
		usr := dto.User{}
		err = rows.StructScan(&usr)
		if err != nil {
			store.Log.Error(fmt.Sprintf("scan failed: %s", err.Error()))
			return nil, fmt.Errorf("QUERY_SCANNING_FAILED")
		}

		// appends
		userList = append(userList, &usr)
	}

	return userList, nil
}

// GetUserCredByUsername fetches user data by username
func (store *Store) GetUserCredByUsername(username string) (*dto.User, *string, error) {
	var err error
	var row *sql.Row
	record := User{}

	query := fmt.Sprintf(`
		SELECT
		    user_id, user_nm, is_admin, status_cd, pass_hash
		FROM %s
		WHERE user_nm = '%s'
	`, table, username)
	row = store.Db.QueryRow(query)

	err = row.Scan(&record.UserID, &record.UserNm, &record.IsAdmin, &record.StatusCd, &record.Passwd)

	if err != nil && err == sql.ErrNoRows {
		store.Log.Debug(fmt.Sprintf("record is not found in the database"))
		return nil, nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, nil, fmt.Errorf("QUERY_FAILED")
	}

	return record.ToService(), &record.Passwd, nil
}
