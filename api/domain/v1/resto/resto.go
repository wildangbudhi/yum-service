package resto

import "github.com/wildangbudhi/yum-service/domain"

type Resto struct {
	ID     *domain.UUID `json:"id"`
	Name   *string      `json:"name"`
	APNKey *string      `json:"-"`
	FCMKey *string      `json:"-"`
}

type RestoRepository interface {
	GetRestoByID(id *domain.UUID) (*Resto, error, domain.RepositoryErrorType)
	UpdateRestoByID(resto *Resto) (error, domain.RepositoryErrorType)
}

type RestoData struct {
	RestoID                     *domain.UUID `json:"-"`
	Address                     *string      `json:"address"`
	IsFreeWifi                  *bool        `json:"is_free_wifi"`
	IsFreeParking               *bool        `json:"is_free_parking"`
	IsPhysicalDistancingApplied *bool        `json:"is_physical_distancing_applied"`
	IsUsingAC                   *bool        `json:"is_using_ac"`
}

type RestoDataResporitory interface {
	GetRestoDataByRestoID(restoID *domain.UUID) (*RestoData, error, domain.RepositoryErrorType)
	UpdateRestoDataByRestoID(restoData *RestoData) (error, domain.RepositoryErrorType)
}
