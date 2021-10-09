package sql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
)

type customerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) auth.CustomerRepository {
	return &customerRepository{
		db: db,
	}
}

func (repo *customerRepository) GetCustomerByID(id *domain.UUID) (*auth.Customer, error, domain.RepositoryErrorType) {

	var err error

	var queryString string = `
	SELECT 
		phone_number, name, password, apn_key, fcm_key, phone_verified_at
	FROM 
		customer
	WHERE
		id = ?
	`

	var queryResult *sql.Row
	var customer *auth.Customer = &auth.Customer{
		ID: id,
	}

	queryResult = repo.db.QueryRow(queryString, id.GetValue())

	err = queryResult.Scan(
		&customer.PhoneNumber,
		&customer.Name,
		&customer.Password,
		&customer.APNKey,
		&customer.FCMKey,
		&customer.PhoneVerifiedAt,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Data Not Found"), domain.RepositoryDataNotFound
		}

		log.Println(err)
		return nil, fmt.Errorf("Services Unavailable"), domain.RepositoryError
	}

	return customer, nil, 0

}

func (repo *customerRepository) GetCustomerByPhoneNumber(phoneNumber string) (*auth.Customer, error, domain.RepositoryErrorType) {

	var err error

	var queryString string = `
	SELECT 
		id, name, password, apn_key, fcm_key, phone_verified_at
	FROM 
		customer
	WHERE
		phone_number = ?
	`

	var queryResult *sql.Row
	var customer *auth.Customer = &auth.Customer{
		PhoneNumber: &phoneNumber,
	}

	queryResult = repo.db.QueryRow(queryString, phoneNumber)

	err = queryResult.Scan(
		&customer.ID,
		&customer.Name,
		&customer.Password,
		&customer.APNKey,
		&customer.FCMKey,
		&customer.PhoneVerifiedAt,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Data Not Found"), domain.RepositoryDataNotFound
		}

		log.Println(err)
		return nil, fmt.Errorf("Services Unavailable"), domain.RepositoryError
	}

	return customer, nil, 0

}

func (repo *customerRepository) CreateCustomer(customer *auth.Customer) (*domain.UUID, error, domain.RepositoryErrorType) {

	var err error
	var exisitingCustomer *auth.Customer

	exisitingCustomer, err, _ = repo.GetCustomerByPhoneNumber(*customer.PhoneNumber)

	if err == nil && exisitingCustomer != nil {
		return nil, fmt.Errorf("User with same phone number already exsist"), domain.RepositoryCreateDataFailed
	}

	customer.ID = domain.NewUUID()

	var queryString string = `
	INSERT INTO customer
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
		customer.ID,
		customer.PhoneNumber,
		customer.Password,
		customer.Name,
		customer.APNKey,
		customer.FCMKey,
		customer.PhoneVerifiedAt,
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

	return customer.ID, nil, 0

}

func (repo *customerRepository) UpdateCustomer(customer *auth.Customer) (error, domain.RepositoryErrorType) {

	var err error
	var queryString string

	queryString = `
	UPDATE 
		customer
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
		customer.PhoneNumber,
		customer.Password,
		customer.Name,
		customer.APNKey,
		customer.FCMKey,
		customer.PhoneVerifiedAt,
		customer.ID,
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
