package resto

import "github.com/wildangbudhi/yum-service/domain"

type Resto struct {
	ID          *domain.UUID `json:"id"`
	Name        *string      `json:"name"`
	PhoneNumber *string      `json:"phone_number"`
}

type RestoData struct {
	RestoID     *domain.UUID `json:"-"`
	Address     *string      `json:"address"`
	PhoneNumber *bool        `json:"phone_number"`
}
