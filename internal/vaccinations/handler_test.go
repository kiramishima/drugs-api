package vaccinations

import (
	"fmt"
	"kiramishima/ionix/internal/mocks"
	"kiramishima/ionix/internal/models"
	"kiramishima/ionix/internal/pkg/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/unrolled/render"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestHandler_ListDrugsHandler(t *testing.T) {
	t.Setenv("JWT_PRIVATE_KEY", "Megaman")
	t.Setenv("TOKEN_TTL", "3600")

	testCases := map[string]struct {
		ID            any
		buildStubs    func(uc *mocks.MockVaccinationService)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		"Getting Data": {
			ID: 1,
			buildStubs: func(uc *mocks.MockVaccinationService) {
				var data = []*models.Vaccination{
					{
						ID:        1,
						Name:      "Jhon Wick",
						Drug:      "medicament 1",
						DrugID:    1,
						Dose:      1,
						AppliedAt: time.Now(),
					},
					{
						ID:        2,
						Name:      "Jhon Connor",
						Drug:      "medicament 1",
						DrugID:    1,
						Dose:      1,
						AppliedAt: time.Now(),
					},
				}

				uc.EXPECT().
					GetListVaccinations(gomock.Any()).
					Return(data, nil).
					AnyTimes()
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				assert.Equal(t, http.StatusOK, recorder.Code)
				// assert.Equal(t, recorder.Body.String(), `{"access_token":"123456"}`)
			},
		},
		"Getting empty": {
			ID: 2,
			buildStubs: func(uc *mocks.MockVaccinationService) {
				uc.EXPECT().
					GetListVaccinations(gomock.Any()).
					Return(nil, ErrNoRecords).
					AnyTimes()
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				t.Log(recorder.Body.String())
				// assert.Equal(t, recorder.Body.String(), `{"error":"El campo password es requerido"}`)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc := mocks.NewMockVaccinationService(ctrl)
			tc.buildStubs(uc)

			recorder := httptest.NewRecorder()
			url := "/v1/vaccination"

			// Generate Token
			token, _ := utils.GenerateJWT(&models.User{ID: 2})
			jwtToken := fmt.Sprintf("Bearer %s", token)
			t.Log("JWT -> ", jwtToken)
			h := http.Header{}
			h.Set("Authorization", jwtToken)

			request := httptest.NewRequest(http.MethodGet, url, nil)
			// request.Header.Set(k, v[0])
			for k, v := range h {
				request.Header.Set(k, v[0])
				recorder.Header().Set(k, v[0])
			}

			router := chi.NewRouter()
			logger, _ := zap.NewProduction()
			slogger := logger.Sugar()
			validate := validator.New()
			r := render.New()

			NewVaccionationHandlers(router, slogger, uc, r, validate)

			// router.ServeHTTP(recorder, request)
			ts := httptest.NewServer(router)
			defer ts.Close()

			tc.checkResponse(t, recorder)
		})
	}
}
