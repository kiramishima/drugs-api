package vaccinations

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/unrolled/render"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"kiramishima/ionix/internal/mocks"
	"kiramishima/ionix/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_CreateAssignHandler(t *testing.T) {
	testCases := map[string]struct {
		ID            any
		invest        int32
		buildStubs    func(uc *mocks.MockCreditService)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		"OK": {
			ID:     1,
			invest: 3000,
			buildStubs: func(uc *mocks.MockCreditService) {
				uc.EXPECT().
					Assign(gomock.Any(), int32(3000)).
					Times(1).
					Return(&models.Credit{Invest: 3000, Credit300: 2, Credit500: 2, Credit700: 2, Status: 1}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				assert.Equal(t, recorder.Body.String(), `{"credit_type_300":2,"credit_type_500":2,"credit_type_700":2}`)
				t.Log(recorder.Body.String())
			},
		},
		"No valid": {
			ID:     2,
			invest: 50,
			buildStubs: func(uc *mocks.MockCreditService) {
				uc.EXPECT().
					Assign(gomock.Any(), int32(50)).
					Times(1).
					Return(nil, errors.New("investment needs be multiply of 100"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
				assert.Equal(t, recorder.Body.String(), `{"error":"investment needs be multiply of 100"}`)
			},
		},
		"No valid - Less than 300": {
			ID:     3,
			invest: 200,
			buildStubs: func(uc *mocks.MockCreditService) {
				uc.EXPECT().
					Assign(gomock.Any(), int32(200)).
					Times(1).
					Return(nil, errors.New("Can't assign credits with this amount of investment"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
				assert.Equal(t, recorder.Body.String(), `{"error":"Can't assign credits with this amount of investment"}`)
			},
		},
		"No valid - Less than 650": {
			ID:     4,
			invest: 650,
			buildStubs: func(uc *mocks.MockCreditService) {
				uc.EXPECT().
					Assign(gomock.Any(), int32(650)).
					Times(1).
					Return(nil, errors.New("investment needs be multiply of 100"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
				assert.Equal(t, recorder.Body.String(), `{"error":"investment needs be multiply of 100"}`)
			},
		},
		"Valid - 600": {
			ID:     5,
			invest: 600,
			buildStubs: func(uc *mocks.MockCreditService) {
				uc.EXPECT().
					Assign(gomock.Any(), int32(600)).
					Times(1).
					Return(&models.Credit{Invest: 600, Credit300: 2, Credit500: 0, Credit700: 0, Status: 1}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				t.Log(recorder.Body.String())
				assert.Equal(t, recorder.Body.String(), `{"credit_type_300":2,"credit_type_500":0,"credit_type_700":0}`)
			},
		},
		"No negatives": {
			ID:     6,
			invest: -3000,
			buildStubs: func(uc *mocks.MockCreditService) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
				t.Log(recorder.Body.String())
				assert.Equal(t, recorder.Body.String(), `{"error":"The value needs to be more than 0 and non negative"}`)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc := mocks.NewMockCreditService(ctrl)
			tc.buildStubs(uc)

			recorder := httptest.NewRecorder()

			url := "/v1/credits/credit-assignment"
			data := models.CreditPostFormRequest{
				Investment: tc.invest,
			}
			// marshall data to json (like json_encode)
			marshalled, err := json.Marshal(data)
			if err != nil {
				log.Fatalf("impossible to marshall form: %s", err)
			}

			request := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(marshalled))

			// assert.NoError(t, err)

			router := chi.NewRouter()
			logger, _ := zap.NewProduction()
			slogger := logger.Sugar()
			validate := validator.New()
			r := render.New()

			NewCreditHandlers(router, slogger, uc, r, validate)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandler_GetStatistics(t *testing.T) {
	testCases := map[string]struct {
		ID            any
		buildStubs    func(uc *mocks.MockCreditService)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		"OK": {
			ID: 1,
			buildStubs: func(uc *mocks.MockCreditService) {
				var item = &models.Stats{
					TotalAssigns:        26,
					TotalSuccessAssigns: 16,
					TotalFailAssigns:    10,
					AVGSuccessAssigns:   70.13,
					AVGFailAssigns:      29.87,
				}
				uc.EXPECT().
					Stats(gomock.Any()).
					Times(1).
					Return(item, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				assert.Equal(t, recorder.Body.String(), `{"total_assigns":26,"total_success_assigns":16,"total_fail_assigns":10,"avg_success_assigns":70.13,"avg_fail_assigns":29.87}`)
				t.Log(recorder.Body.String())
			},
		},
		"Empty": {
			ID: 2,
			buildStubs: func(uc *mocks.MockCreditService) {
				uc.EXPECT().
					Stats(gomock.Any()).
					Times(1).
					Return(&models.Stats{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				t.Log(recorder.Body.String())
				assert.Equal(t, recorder.Body.String(), `{"total_assigns":0,"total_success_assigns":0,"total_fail_assigns":0,"avg_success_assigns":0,"avg_fail_assigns":0}`)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc := mocks.NewMockCreditService(ctrl)
			tc.buildStubs(uc)

			recorder := httptest.NewRecorder()

			url := "/v1/credits/statistics"
			request := httptest.NewRequest(http.MethodGet, url, nil)

			// assert.NoError(t, err)

			router := chi.NewRouter()
			logger, _ := zap.NewProduction()
			slogger := logger.Sugar()
			validate := validator.New()
			r := render.New()

			NewCreditHandlers(router, slogger, uc, r, validate)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
