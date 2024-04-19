package drugs

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

var _ impl.DrugsHandlers = (*handler)(nil)

// NewDrugHandlers creates a instance of drug handlers
func NewDrugHandlers(r *chi.Mux, logger *zap.SugaredLogger, s impl.DrugService, render *render.Render, validate *validator.Validate) {
	var tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("JWT_PRIVATE_KEY")), nil)
	// logger.Info("token ->", tokenAuth)
	handler := &handler{
		logger:   logger,
		service:  s,
		response: render,
		validate: validate,
	}

	r.Route("/v1/drugs", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Get("/", handler.ListDrugsHandler)
		r.Post("/", handler.CreateDrugHandler)
		r.Put("/{id}", handler.UpdateDrugHandler)
		r.Delete("/{id}", handler.DeleteDrugHandler)
	})
}

type handler struct {
	logger   *zap.SugaredLogger
	service  impl.DrugService
	response *render.Render
	validate *validator.Validate
}

func (h handler) ListDrugsHandler(w http.ResponseWriter, req *http.Request) {
	// context
	ctx := req.Context()
	// Call Service
	resp, err := h.service.GetListDrugs(ctx)
	h.logger.Info(resp)
	if err != nil {
		select {
		case <-ctx.Done():
			_ = h.response.JSON(w, http.StatusGatewayTimeout, models.ErrorResponse{ErrorMessage: "El tiempo para procesar su petición ha excedido"})
		default:
			h.logger.Info(err.Error())
			if errors.Is(err, ErrNoRecords) {
				_ = h.response.JSON(w, http.StatusOK, models.ResponseWrapper[[]*models.Drug]{Data: make([]*models.Drug, 0)})
			} else if errors.Is(err, ErrExecuteStatement) {
				_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Error al procesar su petición"})
			} else {
				_ = h.response.JSON(w, http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: "Ocurrio un error interno. Por favor intente más tarde"})
			}
		}
		return
	}

	if err := h.response.JSON(w, http.StatusOK, models.ResponseWrapper[[]*models.Drug]{Data: resp}); err != nil {
		h.logger.Error(err)
		_ = h.response.JSON(w, http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: "Ocurrio un error interno. Por favor intente más tarde"})
		return
	}
}

func (h handler) CreateDrugHandler(w http.ResponseWriter, req *http.Request) {
	var form = &models.DrugForm{}

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

	err = h.service.NewDrug(ctx, form)
	if err != nil {
		// h.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			_ = h.response.JSON(w, http.StatusGatewayTimeout, models.ErrorResponse{ErrorMessage: "Tiempo de ejecución"})
		default:
			if errors.Is(err, ErrDuplicateDrug) {
				_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Este medicamento ya se ha dado de alta con anterioridad"})
			} else if errors.Is(err, ErrExecuteStatement) {
				_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Ocurrio un errro por favor intente más tarde"})
			} else {
				_ = h.response.JSON(w, http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: "Ocurrio un errro por favor intente más tarde"})
			}
		}
		return
	}

	if err := h.response.JSON(w, http.StatusOK, models.Message{Message: "Se ha registrado el nuevo medicamento de manera exitosa"}); err != nil {
		h.logger.Error(err)
		_ = h.response.JSON(w, http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: InternalServerError.Error()})
		return
	}
}

func (h handler) UpdateDrugHandler(w http.ResponseWriter, req *http.Request) {
	var DrugID, _ = strconv.ParseInt(chi.URLParam(req, "id"), 10, 10)
	var form = &models.DrugForm{}

	err := httpUtils.ReadJSON(w, req, &form)

	if err != nil {
		h.logger.Error(err.Error())
		_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: ErrInvalidRequestBody.Error()})
		return
	}

	h.logger.Info(form)
	// context
	ctx := req.Context()

	err = h.service.UpdateDrug(ctx, int(DrugID), form)
	if err != nil {
		// h.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			_ = h.response.JSON(w, http.StatusGatewayTimeout, models.ErrorResponse{ErrorMessage: "Tiempo de ejecución"})
		default:
			if errors.Is(err, ErrDuplicateDrug) {
				_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Este medicamento ya se ha dado de alta con anterioridad"})
			} else if errors.Is(err, ErrDrugNotFound) {
				_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Este medicamento no existe"})
			} else if errors.Is(err, ErrExecuteStatement) {
				_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Ocurrio un errro por favor intente más tarde"})
			} else {
				_ = h.response.JSON(w, http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: "Ocurrio un errro por favor intente más tarde"})
			}
		}
		return
	}

	if err := h.response.JSON(w, http.StatusOK, models.Message{Message: "Se ha actualizado la información del medicamento de manera exitosa"}); err != nil {
		h.logger.Error(err)
		_ = h.response.JSON(w, http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: InternalServerError.Error()})
		return
	}
}

func (h handler) DeleteDrugHandler(w http.ResponseWriter, req *http.Request) {
	var DrugID, _ = strconv.ParseInt(chi.URLParam(req, "id"), 10, 10)
	// context
	ctx := req.Context()

	err := h.service.DeleteDrug(ctx, int(DrugID))
	if err != nil {
		select {
		case <-ctx.Done():
			_ = h.response.JSON(w, http.StatusGatewayTimeout, models.ErrorResponse{ErrorMessage: "Tiempo de ejecución"})
		default:
			if errors.Is(err, ErrDrugNotFound) {
				_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Este medicamento no existe"})
			} else if errors.Is(err, ErrExecuteStatement) {
				_ = h.response.JSON(w, http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Ocurrio un error por favor intente más tarde"})
			} else {
				_ = h.response.JSON(w, http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: "Ocurrio un errro por favor intente más tarde"})
			}
		}
		return
	}

	if err := h.response.JSON(w, http.StatusOK, models.Message{Message: "Se ha eliminado el medicamento de manera exitosa"}); err != nil {
		h.logger.Error(err)
		_ = h.response.JSON(w, http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: InternalServerError.Error()})
		return
	}
}
