package auth

import (
	"bytes"
	"encoding/json"
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

func TestHandler_LoginHandler(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		ID            any
		form          *models.AuthForm
		buildStubs    func(uc *mocks.MockAuthService)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		"Success new account": {
			ID:   1,
			form: &models.AuthForm{Email: "giny@mail.com", Password: "123456"},
			buildStubs: func(uc *mocks.MockAuthService) {
				uc.EXPECT().
					SignIn(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&models.AuthResponse{AccessToken: "123456"}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				assert.Equal(t, recorder.Body.String(), `{"access_token":"123456"}`)
				// t.Log(recorder.Body.String())
			},
		},
		"Password required": {
			ID:   2,
			form: &models.AuthForm{Email: "giny@mail.com", Password: ""},
			buildStubs: func(uc *mocks.MockAuthService) {
				/*uc.EXPECT().
				SignIn(gomock.Any(), gomock.Any()).
				Times(1).
				Return(nil, ErrMissingPassword)*/
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
				t.Log(recorder.Body.String())
				assert.Equal(t, recorder.Body.String(), `{"error":"El campo password es requerido"}`)
			},
		},
		"Bad Email format": {
			ID:   3,
			form: &models.AuthForm{Email: "giny_mail.com", Password: "123456"},
			buildStubs: func(uc *mocks.MockAuthService) {
				/*uc.EXPECT().
				SignIn(gomock.Any(), gomock.Any()).
				Times(1).
				Return(nil, ErrMissingPassword)*/
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
				t.Log(recorder.Body.String())
				assert.Equal(t, recorder.Body.String(), `{"error":"El campo email es invalido"}`)
			},
		},
		"Email required": {
			ID:   4,
			form: &models.AuthForm{Email: "", Password: "123456"},
			buildStubs: func(uc *mocks.MockAuthService) {
				/*uc.EXPECT().
				SignIn(gomock.Any(), gomock.Any()).
				Times(1).
				Return(nil, ErrMissingPassword)*/
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
				t.Log(recorder.Body.String())
				assert.Equal(t, recorder.Body.String(), `{"error":"El campo email es requerido"}`)
			},
		},
		"User not found": {
			ID:   5,
			form: &models.AuthForm{Email: "giny@mail.com", Password: "123456"},
			buildStubs: func(uc *mocks.MockAuthService) {
				uc.EXPECT().
					SignIn(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, ErrUserNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
				assert.Equal(t, recorder.Body.String(), `{"error":"El email y/o contraseña son erroneos"}`)
				t.Log(recorder.Body.String())
			},
		},
		"Bad Password": {
			ID:   6,
			form: &models.AuthForm{Email: "giny@mail.com", Password: "123456"},
			buildStubs: func(uc *mocks.MockAuthService) {
				uc.EXPECT().
					SignIn(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, ErrInvalidPassword)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
				assert.Equal(t, recorder.Body.String(), `{"error":"Contraseña invalida"}`)
				t.Log(recorder.Body.String())
			},
		},
		"General service": {
			ID:   6,
			form: &models.AuthForm{Email: "giny@mail.com", Password: "123456"},
			buildStubs: func(uc *mocks.MockAuthService) {
				uc.EXPECT().
					SignIn(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, ErrServiceAuth)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
				assert.Equal(t, recorder.Body.String(), `{"error":"Ocurrio un error por favor intente más tarde"}`)
				t.Log(recorder.Body.String())
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc := mocks.NewMockAuthService(ctrl)
			tc.buildStubs(uc)

			recorder := httptest.NewRecorder()

			url := "/v1/auth/sign-in"
			data := tc.form
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

			NewAuthHandlers(router, slogger, uc, r, validate)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandler_SignUpHandler(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		ID            any
		form          *models.RegisterForm
		buildStubs    func(uc *mocks.MockAuthService)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		"OK": {
			ID:   1,
			form: &models.RegisterForm{Email: "giny@mail.com", Password: "123456", Name: "Jhon"},
			buildStubs: func(uc *mocks.MockAuthService) {
				uc.EXPECT().
					SignUp(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				assert.Equal(t, recorder.Body.String(), `{"message":"Registro exitoso."}`)
				t.Log(recorder.Body.String())
			},
		},
		"Password required": {
			ID:   2,
			form: &models.RegisterForm{Email: "giny@mail.com", Password: "", Name: "Jhon"},
			buildStubs: func(uc *mocks.MockAuthService) {
				/*uc.EXPECT().
				SignIn(gomock.Any(), gomock.Any()).
				Times(1).
				Return(nil, ErrMissingPassword)*/
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
				t.Log(recorder.Body.String())
				assert.Equal(t, recorder.Body.String(), `{"error":"El campo password es requerido"}`)
			},
		},
		"Bad Email format": {
			ID:   3,
			form: &models.RegisterForm{Email: "giny[at]mail.com", Password: "123456", Name: "Jhon"},
			buildStubs: func(uc *mocks.MockAuthService) {
				/*uc.EXPECT().
				SignIn(gomock.Any(), gomock.Any()).
				Times(1).
				Return(nil, ErrMissingPassword)*/
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
				t.Log(recorder.Body.String())
				assert.Equal(t, recorder.Body.String(), `{"error":"El campo email es invalido"}`)
			},
		},
		"Email required": {
			ID:   4,
			form: &models.RegisterForm{Email: "", Password: "123456", Name: "Jhon"},
			buildStubs: func(uc *mocks.MockAuthService) {
				/*uc.EXPECT().
				SignIn(gomock.Any(), gomock.Any()).
				Times(1).
				Return(nil, ErrMissingPassword)*/
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
				t.Log(recorder.Body.String())
				assert.Equal(t, recorder.Body.String(), `{"error":"El campo email es requerido"}`)
			},
		},
		"User exists": {
			ID:   5,
			form: &models.RegisterForm{Email: "giny@mail.com", Password: "123456", Name: "Jhon"},
			buildStubs: func(uc *mocks.MockAuthService) {
				uc.EXPECT().
					SignUp(gomock.Any(), gomock.Any()).
					Times(1).
					Return(ErrUserExist)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
				assert.Equal(t, recorder.Body.String(), `{"error":"Ya existe una cuenta registrada con este correo"}`)
				t.Log(recorder.Body.String())
			},
		},
		"General service": {
			ID:   6,
			form: &models.RegisterForm{Email: "giny@mail.com", Password: "123456", Name: "Jhon"},
			buildStubs: func(uc *mocks.MockAuthService) {
				uc.EXPECT().
					SignUp(gomock.Any(), gomock.Any()).
					Times(1).
					Return(ErrServiceAuth)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
				assert.Equal(t, recorder.Body.String(), `{"error":"Ocurrio un error por favor intente más tarde"}`)
				t.Log(recorder.Body.String())
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc := mocks.NewMockAuthService(ctrl)
			tc.buildStubs(uc)

			recorder := httptest.NewRecorder()

			url := "/v1/auth/sign-up"
			data := tc.form
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

			NewAuthHandlers(router, slogger, uc, r, validate)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
