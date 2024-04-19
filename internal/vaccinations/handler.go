package vaccinations

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/unrolled/render"
	"go.uber.org/zap"
	impl "kiramishima/ionix/internal/interfaces"
	"kiramishima/ionix/internal/models"
	httpUtils "kiramishima/ionix/internal/pkg/utils"
	"net/http"
	"os"
	"strconv"
)

var _ impl.VaccinationsHandlers = (*handler)(nil)

// NewVaccionationHandlers creates a instance of vaccination handlers
func NewVaccionationHandlers(r *chi.Mux, logger *zap.SugaredLogger, s impl.VaccinationService, render *render.Render, validate *validator.Validate) {
	var tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("JWT_PRIVATE_KEY")), nil)
	logger.Info("token ->", tokenAuth)
	handler := &handler{
		logger:   logger,
		service:  s,
		response: render,
		validate: validate,
	}

	r.Route("/v1/vaccination", func(r chi.Router) {
		r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator(tokenAuth)).Get("/", handler.ListVaccinationsHandler)
		r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator(tokenAuth)).Post("/", handler.CreateVaccinationHandler)
		r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator(tokenAuth)).Put("/{id}", handler.UpdateVaccinationHandler)
		r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator(tokenAuth)).Delete("/{id}", handler.DeleteVaccinationHandler)
	})
}

type handler struct {
	logger   *zap.SugaredLogger
	service  impl.VaccinationService
	response *render.Render
	validate *validator.Validate
}

func (h handler) ListVaccinationsHandler(w http.ResponseWriter, req *http.Request) {
	// context
	ctx := req.Context()
	// Call Service
	resp, err := h.service.GetListVaccinations(ctx)
	h.logger.Info(resp)
	if err != nil {
		select {
		case <-ctx.Done():
			_ = h.response.JSON(w, http.StatusGatewayTimeout, models.ErrorResponse{ErrorMessage: "El tiempo para procesar su petición ha excedido"})
		default:
			h.logger.Info(err.Error())
			if errors.Is(err, ErrNoRecords) {
				_ = h.response.JSON(w, http.StatusOK, models.ResponseWrapper[[]*models.Vaccination]{Data: make([]*models.Vaccination, 0)})
			} else if errors.Is(err, ErrExecuteStatement) {
				_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Error al procesar su petición"})
			} else {
				_ = h.response.JSON(w, http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: "Ocurrio un error interno. Por favor intente más tarde"})
			}
		}
		return
	}

	if err := h.response.JSON(w, http.StatusOK, models.ResponseWrapper[[]*models.Vaccination]{Data: resp}); err != nil {
		h.logger.Error(err)
		_ = h.response.JSON(w, http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: "Ocurrio un error interno. Por favor intente más tarde"})
		return
	}
}

func (h handler) CreateVaccinationHandler(w http.ResponseWriter, req *http.Request) {
	var form = &models.VaccinationForm{}

	err := httpUtils.ReadJSON(w, req, &form)

	if err != nil {
		h.logger.Error(err.Error())
		_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: ErrInvalidRequestBody.Error()})
		return
	}
	h.logger.Info(form)
	// Validate form
	err = form.Validate(h.validate)
	if err != nil {
		h.logger.Error(err.Error())
		_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: err.Error()})
		return
	}
	// context
	ctx := req.Context()

	err = h.service.NewVaccination(ctx, form)
	if err != nil {
		// h.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			_ = h.response.JSON(w, http.StatusGatewayTimeout, models.ErrorResponse{ErrorMessage: "Tiempo de ejecución"})
		default:
			if errors.Is(err, ErrDuplicateVaccination) {
				_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Este registro ya se había dado de alta con anterioridad"})
			} else if errors.Is(err, ErrExecuteStatement) {
				_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Ocurrio un error por favor intente más tarde"})
			} else {
				_ = h.response.JSON(w, http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: "Ocurrio un errro por favor intente más tarde"})
			}
		}
		return
	}

	if err := h.response.JSON(w, http.StatusOK, models.Message{Message: "Se ha registrado de manera exitosa"}); err != nil {
		h.logger.Error(err)
		_ = h.response.JSON(w, http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: InternalServerError.Error()})
		return
	}
}

func (h handler) UpdateVaccinationHandler(w http.ResponseWriter, req *http.Request) {
	var VacID, _ = strconv.ParseInt(chi.URLParam(req, "id"), 10, 10)
	var form = &models.VaccinationForm{}

	err := httpUtils.ReadJSON(w, req, &form)

	if err != nil {
		h.logger.Error(err.Error())
		_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: ErrInvalidRequestBody.Error()})
		return
	}

	h.logger.Info(VacID, form)
	// context
	ctx := req.Context()

	err = h.service.UpdateVaccination(ctx, int(VacID), form)
	if err != nil {
		// h.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			_ = h.response.JSON(w, http.StatusGatewayTimeout, models.ErrorResponse{ErrorMessage: "Tiempo de ejecución"})
		default:
			if errors.Is(err, ErrDuplicateVaccination) {
				_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Este registro ya se ha dado de alta con anterioridad"})
			} else if errors.Is(err, ErrVaccinationNotFound) {
				_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Este registro no existe"})
			} else if errors.Is(err, ErrExecuteStatement) {
				_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Ocurrio un error por favor intente más tarde"})
			} else {
				_ = h.response.JSON(w, http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: "Ocurrio un errro por favor intente más tarde"})
			}
		}
		return
	}

	if err := h.response.JSON(w, http.StatusOK, models.Message{Message: "Se ha actualizado la información de manera exitosa"}); err != nil {
		h.logger.Error(err)
		_ = h.response.JSON(w, http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: InternalServerError.Error()})
		return
	}
}

func (h handler) DeleteVaccinationHandler(w http.ResponseWriter, req *http.Request) {
	var VacID, _ = strconv.ParseInt(chi.URLParam(req, "id"), 10, 10)
	// context
	ctx := req.Context()

	err := h.service.DeleteVaccination(ctx, int(VacID))
	if err != nil {
		select {
		case <-ctx.Done():
			_ = h.response.JSON(w, http.StatusGatewayTimeout, models.ErrorResponse{ErrorMessage: "Tiempo de ejecución"})
		default:
			if errors.Is(err, ErrVaccinationNotFound) {
				_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Este registro no existe"})
			} else if errors.Is(err, ErrExecuteStatement) {
				_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Ocurrio un error por favor intente más tarde"})
			} else {
				_ = h.response.JSON(w, http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: "Ocurrio un error por favor intente más tarde"})
			}
		}
		return
	}

	if err := h.response.JSON(w, http.StatusOK, models.Message{Message: "Se ha eliminado el registro de manera exitosa"}); err != nil {
		h.logger.Error(err)
		_ = h.response.JSON(w, http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: InternalServerError.Error()})
		return
	}
}
