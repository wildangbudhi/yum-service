package sql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/resto"
)

type restoRepository struct {
	db *sql.DB
}

func NewRestoRepository(db *sql.DB) resto.RestoRepository {
	return &restoRepository{
		db: db,
	}
}

func (repo *restoRepository) GetRestoByID(id *domain.UUID) (*resto.Resto, error, domain.RepositoryErrorType) {

	var err error

	var queryString string = `
	SELECT name, apn_key, fcm_key
	FROM resto
	WHERE id = ?
	`

	var queryResult *sql.Row
	var resto *resto.Resto = &resto.Resto{
		ID: id,
	}

	queryResult = repo.db.QueryRow(queryString, id.GetValue())

	err = queryResult.Scan(
		&resto.Name,
		&resto.APNKey,
		&resto.FCMKey,
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

func (repo *restoRepository) UpdateRestoByID(resto *resto.Resto) (error, domain.RepositoryErrorType) {

	var err error
	var queryString string

	queryString = `
	UPDATE 
		resto
	SET 
		name=?,
		apn_key=?, 
		fcm_key=?, 
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
		resto.Name,
		resto.APNKey,
		resto.FCMKey,
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
