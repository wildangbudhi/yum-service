package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
	"github.com/wildangbudhi/yum-service/services/v1/auth/delivery/httprest"
	"github.com/wildangbudhi/yum-service/services/v1/auth/delivery/middleware"
	"github.com/wildangbudhi/yum-service/services/v1/auth/repository/redis"
	"github.com/wildangbudhi/yum-service/services/v1/auth/repository/sms"
	"github.com/wildangbudhi/yum-service/services/v1/auth/repository/sql"
	"github.com/wildangbudhi/yum-service/services/v1/auth/usecase"
	"github.com/wildangbudhi/yum-service/utils"
)

func getAuthUsecase(server *utils.Server) auth.AuthUsecase {

	var authUsecase auth.AuthUsecase

	var sessionRepository auth.SessionRepository
	var phoneVerificationRepository auth.PhoneVerificationRepository
	var customerRepository auth.CustomerRepository
	var restoRepository auth.RestoRepository
	var otpRepository auth.OTPRepository

	sessionRepository = redis.NewSessionRepository(server.RedisDB)
	phoneVerificationRepository = sms.NewPhoneVerificationRepository(server.Config, server.SMSService)
	customerRepository = sql.NewCustomerRepository(server.DB)
	restoRepository = sql.NewRestoRepository(server.DB)
	otpRepository = sql.NewOTPRepository(server.DB)

	authUsecase = usecase.NewAuthUsecase(
		&server.Config,
		sessionRepository,
		phoneVerificationRepository,
		customerRepository,
		restoRepository,
		otpRepository,
	)

	return authUsecase

}

func AuthHTTPRestDI(server *utils.Server) {

	var route *gin.RouterGroup = server.Router.Group("/v1/auth")
	var authUsecase auth.AuthUsecase = getAuthUsecase(server)
	var authMiddlewareDelivery auth.AuthMiddlewareDelivery = AuthMiddlewareDI(server)

	httprest.NewAuthHTTPRestHandler(route, authMiddlewareDelivery, authUsecase)

}

func AuthMiddlewareDI(server *utils.Server) auth.AuthMiddlewareDelivery {

	var delivery auth.AuthMiddlewareDelivery
	var authUsecase auth.AuthUsecase = getAuthUsecase(server)

	delivery = middleware.NewAuthMiddlewareDelivery(authUsecase)

	return delivery

}
