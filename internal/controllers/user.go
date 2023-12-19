package controllers

import (
	"encoding/json"
	"github.com/TechBuilder-360/Auth_Server/internal/common/constant"
	"github.com/TechBuilder-360/Auth_Server/internal/common/types"
	"github.com/TechBuilder-360/Auth_Server/internal/common/utils"
	"github.com/TechBuilder-360/Auth_Server/internal/middlewares"
	"github.com/TechBuilder-360/Auth_Server/internal/services"
	"github.com/TechBuilder-360/Auth_Server/internal/validation"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type IUserController interface {
	UpgradeTierOne(w http.ResponseWriter, r *http.Request)
	RegisterRoutes(router *mux.Router)
}

type UserController struct {
	as services.UserService
}

func (c *UserController) RegisterRoutes(router *mux.Router) {
	apis := router.PathPrefix("/users").Subrouter()
	apis.HandleFunc("/upgrade/tier-one", middlewares.Adapt(http.HandlerFunc(c.UpgradeTierOne),
		middlewares.AuthorizeUserJWT()).ServeHTTP).Methods(http.MethodPost)

}

func DefaultUserController() IUserController {
	return &UserController{
		as: services.NewUserService(),
	}
}

func (c *UserController) UpgradeTierOne(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{constant.RequestIdentifier: utils.GenerateUUID()})
	logger.Info("Upgrading user tiers")
	body := &types.UpgradeUserTierRequest{}
	json.NewDecoder(r.Body).Decode(body)
	logger.Info("Request data: %+v", body)

	if validation.ValidateStruct(w, body, logger) {
		return
	}

	// get user from context
	user, err := middlewares.UserFromContext(r)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.ErrorResponse{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	err = c.as.UpgradeTierOne(body, user, logger)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error(err.Error())
		json.NewEncoder(w).Encode(utils.ErrorResponse{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.SuccessResponse{
		Status:  true,
		Message: "Tier one upgrade successful",
	})
}