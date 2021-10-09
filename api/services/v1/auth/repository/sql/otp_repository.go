package sql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
)

type otpRepository struct {
	db *sql.DB
}

func NewOTPRepository(db *sql.DB) auth.OTPRepository {
	return &otpRepository{
		db: db,
	}
}

func (repo *otpRepository) CountOTPWithin30Second(id string, userType int, phoneNumber string) (int, error, domain.RepositoryErrorType) {

	var err error

	var queryString string = `
	SELECT 
		COUNT(id)
	FROM 
		yum.otp
	WHERE
		id = ?
		AND type = ?
		AND phone_number = ?
		AND created_at BETWEEN DATE_SUB(CURRENT_TIMESTAMP(), INTERVAL 30 SECOND) AND CURRENT_TIMESTAMP()  
	`

	var queryResult *sql.Row
	var countOTP int = 0

	queryResult = repo.db.QueryRow(queryString, id, userType, phoneNumber)
	err = queryResult.Scan(&countOTP)

	if err != nil {

		if err == sql.ErrNoRows {
			return -1, fmt.Errorf("Data Not Found"), domain.RepositoryDataNotFound
		}

		log.Println(err)
		return -1, fmt.Errorf("Services Unavailable"), domain.RepositoryError
	}

	return countOTP, nil, 0

}

func (repo *otpRepository) CreateNewOTP(otp *auth.OTP) (error, domain.RepositoryErrorType) {

	var err error

	var queryString string = `
	INSERT INTO otp
	(id, type, phone_number, sid, create_verification_resp_json, verification_check_resp_json)
	VALUES(?, ?, ?, ?, ?, ?)
	`

	statement, err := repo.db.Prepare(queryString)

	if err != nil {
		log.Println(err)
		return fmt.Errorf("Services Unavailable"), domain.RepositoryError
	}

	var res sql.Result

	res, err = statement.Exec(
		otp.ID,
		otp.Type,
		otp.PhoneNumber,
		otp.SID,
		otp.CreateVerificationRespJSON,
		otp.VerificationCheckRespJSON,
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
		return fmt.Errorf("Failed to Save Partner Data"), domain.RepositoryCreateDataFailed
	}

	return nil, 0

}

func (repo *otpRepository) UpdateOTP(otp *auth.OTP) (error, domain.RepositoryErrorType) {

	var err error
	var queryString string

	queryString = `
	UPDATE otp
	SET verification_check_resp_json=?, updated_at=CURRENT_TIMESTAMP
	WHERE 
		id = ?
		AND type = ?
		AND sid = ?
		AND phone_number = ?
	`

	statement, err := repo.db.Prepare(queryString)

	if err != nil {
		log.Println(err)
		return fmt.Errorf("Services Unavailable"), domain.RepositoryError
	}

	var res sql.Result

	res, err = statement.Exec(
		otp.VerificationCheckRespJSON,
		otp.ID,
		otp.Type,
		otp.SID,
		otp.PhoneNumber,
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
