package resto

import (
	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
)

type RestoUsecase interface {
	GetRestoData(authHeader *auth.ValidateAuthTokenResponse, restoID *domain.UUID)
}
