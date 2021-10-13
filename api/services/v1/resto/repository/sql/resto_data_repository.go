package sql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/resto"
)

type restoDataRepository struct {
	db *sql.DB
}

func NewRestoDataRepository(db *sql.DB) resto.RestoDataResporitory {
	return &restoDataRepository{
		db: db,
	}
}

func (repo *restoDataRepository) GetRestoDataByRestoID(restoID *domain.UUID) (*resto.RestoData, error, domain.RepositoryErrorType) {

	var err error

	var queryString string = `
	SELECT 
		address, 
		is_free_wifi, 
		is_free_parking, 
		is_physical_distancing_applied, 
		is_using_ac
	FROM 
		resto_data
	WHERE 
		resto_id = ?
	`

	var queryResult *sql.Row
	var restoData *resto.RestoData = &resto.RestoData{
		RestoID: restoID,
	}

	queryResult = repo.db.QueryRow(queryString, restoID.GetValue())

	err = queryResult.Scan(
		&restoData.Address,
		&restoData.IsFreeWifi,
		&restoData.IsFreeParking,
		&restoData.IsPhysicalDistancingApplied,
		&restoData.IsUsingAC,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Data Not Found"), domain.RepositoryDataNotFound
		}

		log.Println(err)
		return nil, fmt.Errorf("Services Unavailable"), domain.RepositoryError
	}

	return restoData, nil, 0

}

func (repo *restoDataRepository) UpdateRestoDataByRestoID(restoData *resto.RestoData) (error, domain.RepositoryErrorType) {

	var err error
	var queryString string

	queryString = `
	UPDATE resto_data
	SET 
		address=?, 
		is_free_wifi=?, 
		is_free_parking=?, 
		is_physical_distancing_applied=?, 
		is_using_ac=?, 
		updated_at=CURRENT_TIMESTAMP()
	WHERE 
		resto_id = ?
	`

	statement, err := repo.db.Prepare(queryString)

	if err != nil {
		log.Println(err)
		return fmt.Errorf("Services Unavailable"), domain.RepositoryError
	}

	var res sql.Result

	res, err = statement.Exec(
		restoData.Address,
		restoData.IsFreeWifi,
		restoData.IsFreeParking,
		restoData.IsPhysicalDistancingApplied,
		restoData.IsUsingAC,
		restoData.RestoID,
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
