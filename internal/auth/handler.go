package auth

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/unrolled/render"
	"go.uber.org/zap"
	impl "kiramishima/ionix/internal/interfaces"
	"kiramishima/ionix/internal/models"
	httpUtils "kiramishima/ionix/internal/pkg/utils"
	"net/http"
)

var _ impl.AuthHandlers = (*handler)(nil)

// NewAuthHandlers creates an instance of auth handlers
func NewAuthHandlers(r *chi.Mux, logger *zap.Logger, s impl.AuthService, render *render.Render, validate *validator.Validate) {
	handler := &handler{
		logger:   logger,
		service:  s,
		response: render,
		validate: validate,
	}

	r.Route("/v1/auth", func(r chi.Router) {
		r.Post("/sign-in", handler.LoginHandler)
		r.Post("/sign-up", handler.SignUpHandler)
	})
}

type handler struct {
	logger   *zap.Logger
	service  impl.AuthService
	response *render.Render
	validate *validator.Validate
}

func (h handler) SignUpHandler(w http.ResponseWriter, req *http.Request) {
	var form = &models.RegisterForm{}

	err := httpUtils.ReadJSON(w, req, &form)

	if err != nil {
		h.logger.Error(err.Error())
		_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "La petición es invalidad o esta mal formateada"})
		return
	}
	h.logger.Info("Form", zap.Any("data", form))
	// Validate form
	err = form.Validate(h.validate)
	if err != nil {
		h.logger.Error(err.Error())
		switch err.Error() {
		case "Password. This field is required":
			_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: ErrMissingPassword.Error()})
		case "Email. This field is required":
			_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: ErrMissingEmail.Error()})
		case "Email. Bad email format":
			_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: ErrInvalidEmail.Error()})
		default:
			_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: err.Error()})
		}
		return
	}
	ctx := req.Context()

	err = h.service.SignUp(ctx, form)
	if err != nil {
		// h.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			_ = h.response.JSON(w, http.StatusGatewayTimeout, models.ErrorResponse{ErrorMessage: err.Error()})
		default:
			if errors.Is(err, ErrServiceAuth) {
				_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Ocurrio un error por favor intente más tarde"})
			} else if errors.Is(err, ErrUserExist) {
				_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Ya existe una cuenta registrada con este correo"})
			} else {
				_ = h.response.JSON(w, http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: "Ocurrio un error por favor intente más tarde"})
			}
		}
		return
	}

	if err := h.response.JSON(w, http.StatusOK, models.Message{Message: "Registro exitoso."}); err != nil {
		h.logger.Error("[ERROR]", zap.Error(err))
		_ = h.response.JSON(w, http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: "Internal Server Error"})
		return
	}
}

func (h handler) LoginHandler(w http.ResponseWriter, req *http.Request) {
	var form = &models.AuthForm{}

	err := httpUtils.ReadJSON(w, req, &form)

	if err != nil {
		h.logger.Error(err.Error())
		_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "La petición es invalida o esta mal formateada"})
		return
	}

	h.logger.Info("[INFO]", zap.Any("FormData", form))
	// Validate data
	err = form.Validate(h.validate)
	if err != nil {
		h.logger.Error(err.Error())
		switch err.Error() {
		case "Password. This field is required":
			_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: ErrMissingPassword.Error()})
		case "Email. This field is required":
			_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: ErrMissingEmail.Error()})
		case "Email. Bad email format":
			_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: ErrInvalidEmail.Error()})
		default:
			_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: err.Error()})
		}
		return
	}
	ctx := req.Context()

	// Service
	resp, err := h.service.SignIn(ctx, form)
	if err != nil {

		select {
		case <-ctx.Done():
			_ = h.response.JSON(w, http.StatusGatewayTimeout, models.ErrorResponse{ErrorMessage: err.Error()})
		default:

			if errors.Is(err, ErrUserNotFound) {
				_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "El email y/o contraseña son erroneos"})
			} else if errors.Is(err, ErrInvalidPassword) {
				_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Contraseña invalida"})
			} else if errors.Is(err, ErrServiceAuth) {
				_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Ocurrio un error por favor intente más tarde"})
			} else {
				_ = h.response.JSON(w, http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: "Ocurrio un error por favor intente más tarde"})
			}
		}
		return
	}

	// response
	if err := h.response.JSON(w, http.StatusOK, resp); err != nil {
		h.logger.Error("[ERROR]", zap.Error(err))
		_ = h.response.JSON(w, http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: "Internal Server Error"})
		return
	}
}
