package sql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
)

type restoRepository struct {
	db *sql.DB
}

func NewRestoRepository(db *sql.DB) auth.RestoRepository {
	return &restoRepository{
		db: db,
	}
}

func (repo *restoRepository) GetRestoByID(id *domain.UUID) (*auth.Resto, error, domain.RepositoryErrorType) {

	var err error

	var queryString string = `
	SELECT 
		phone_number, password, name, apn_key, fcm_key, phone_verified_at
	FROM 
		resto
	WHERE
		id = ?
	`

	var queryResult *sql.Row
	var resto *auth.Resto = &auth.Resto{
		ID: id,
	}

	queryResult = repo.db.QueryRow(queryString, id.GetValue())

	err = queryResult.Scan(
		&resto.PhoneNumber,
		&resto.Password,
		&resto.Name,
		&resto.APNKey,
		&resto.FCMKey,
		&resto.PhoneVerifiedAt,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Data Not Found"), domain.RepositoryDataNotFound
		}

		log.Println(err)
		return nil, fmt.Errorf("Services Unavailable"), domain.RepositoryError
	}

	return resto, nil, 0

}

func (repo *restoRepository) GetRestoByPhoneNumber(phoneNumber string) (*auth.Resto, error, domain.RepositoryErrorType) {

	var err error

	var queryString string = `
	SELECT 
		id, password, name, apn_key, fcm_key, phone_verified_at
	FROM 
		resto
	WHERE
		phone_number = ?
	`

	var queryResult *sql.Row
	var resto *auth.Resto = &auth.Resto{
		PhoneNumber: &phoneNumber,
	}

	queryResult = repo.db.QueryRow(queryString, phoneNumber)

	err = queryResult.Scan(
		&resto.ID,
		&resto.Password,
		&resto.Name,
		&resto.APNKey,
		&resto.FCMKey,
		&resto.PhoneVerifiedAt,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Data Not Found"), domain.RepositoryDataNotFound
		}

		log.Println(err)
		return nil, fmt.Errorf("Services Unavailable"), domain.RepositoryError
	}

	return resto, nil, 0

}

func (repo *restoRepository) CreateResto(resto *auth.Resto) (*domain.UUID, error, domain.RepositoryErrorType) {

	var err error
	var exisitingResto *auth.Resto

	exisitingResto, err, _ = repo.GetRestoByPhoneNumber(*resto.PhoneNumber)

	if err == nil && exisitingResto != nil {
		return nil, fmt.Errorf("User with same phone number already exsist"), domain.RepositoryCreateDataFailed
	}

	resto.ID = domain.NewUUID()

	var queryString string = `
	INSERT INTO yum.resto
	(id, phone_number, password, name, apn_key, fcm_key, phone_verified_at)
	VALUES(?, ?, ?, ?, ?, ?, ?)
	`

	statement, err := repo.db.Prepare(queryString)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Services Unavailable"), domain.RepositoryError
	}

	var res sql.Result

	res, err = statement.Exec(
		resto.ID,
		resto.PhoneNumber,
		resto.Password,
		resto.Name,
		resto.APNKey,
		resto.FCMKey,
		resto.PhoneVerifiedAt,
	)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Services Unavailable"), domain.RepositoryError
	}

	rowAffected, err := res.RowsAffected()

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Services Unavailable"), domain.RepositoryError
	}

	if rowAffected == 0 {
		return nil, fmt.Errorf("Failed to Save Partner Data"), domain.RepositoryCreateDataFailed
	}

	return resto.ID, nil, 0

}

func (repo *restoRepository) UpdateResto(resto *auth.Resto) (error, domain.RepositoryErrorType) {

	var err error
	var queryString string

	queryString = `
	UPDATE 
		resto
	SET 
		phone_number=?, 
		password=?, 
		name=?,
		apn_key=?, 
		fcm_key=?, 
		phone_verified_at=?,
		updated_at=CURRENT_TIMESTAMP()
	WHERE 
		id = ?
	`

	statement, err := repo.db.Prepare(queryString)

	if err != nil {
		log.Println(err)
		return fmt.Errorf("Services Unavailable"), domain.RepositoryError
	}

	var res sql.Result

	res, err = statement.Exec(
		resto.PhoneNumber,
		resto.Password,
		resto.Name,
		resto.APNKey,
		resto.FCMKey,
		resto.PhoneVerifiedAt,
		resto.ID,
	)

	if err != nil {
		log.Println(err)
		return fmt.Errorf("Services Unavailable"), domain.RepositoryError
	}

	rowAffected, err := res.RowsAffected()

	if err != nil {
		log.Println(err)
		return fmt.Errorf("Services Unavailable"), domain.RepositoryError
	}

	if rowAffected == 0 {
		return fmt.Errorf("Failed to Insert New User"), domain.RepositoryUpdateDataFailed
	}

	return nil, 0

}
