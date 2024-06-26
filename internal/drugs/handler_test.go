package drugs

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/unrolled/render"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"kiramishima/ionix/internal/mocks"
	"kiramishima/ionix/internal/models"
	"kiramishima/ionix/internal/pkg/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_ListDrugsHandler(t *testing.T) {
	t.Setenv("JWT_PRIVATE_KEY", "Megaman")
	t.Setenv("TOKEN_TTL", "3600")

	testCases := map[string]struct {
		ID            any
		buildStubs    func(uc *mocks.MockDrugService)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		"Getting Data": {
			ID: 1,
			buildStubs: func(uc *mocks.MockDrugService) {
				var drugs = []*models.Drug{
					{
						ID:          1,
						Name:        "medicament 1",
						Approved:    true,
						MinDose:     1,
						MaxDose:     5,
						AvailableAt: time.Now(),
					},
					{
						ID:          2,
						Name:        "medicament 2",
						Approved:    true,
						MinDose:     1,
						MaxDose:     5,
						AvailableAt: time.Now(),
					},
				}

				uc.EXPECT().
					GetListDrugs(gomock.Any()).
					Return(drugs, nil).
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
			buildStubs: func(uc *mocks.MockDrugService) {
				uc.EXPECT().
					GetListDrugs(gomock.Any()).
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

			uc := mocks.NewMockDrugService(ctrl)
			tc.buildStubs(uc)

			recorder := httptest.NewRecorder()
			url := "/v1/drugs"

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
			logger := zap.NewNop()
			validate := validator.New()
			r := render.New()

			NewDrugHandlers(router, logger, uc, r, validate)

			// router.ServeHTTP(recorder, request)
			ts := httptest.NewServer(router)
			defer ts.Close()

			tc.checkResponse(t, recorder)
		})
	}
}
