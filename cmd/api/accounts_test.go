package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Ruthvik10/simple_bank/internal/logger"
	mock "github.com/Ruthvik10/simple_bank/internal/mock/db"
	"github.com/Ruthvik10/simple_bank/internal/models"
	"github.com/Ruthvik10/simple_bank/internal/store"
	"github.com/go-chi/chi/v5"
)

var app = application{
	store:  store.NewMockStore(),
	logger: logger.New(os.Stdout, logger.LevelError),
}

func Test_application_createAccountHandler_success(t *testing.T) {
	reqBody := `{
					"owner": "Ruthvik",
					"balance": 20000,
					"Currency": "USD"
				}`
	handler := http.HandlerFunc(app.CreateAccountHandler)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/accounts/", strings.NewReader(reqBody))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Result().StatusCode != http.StatusCreated {
		t.Errorf("expected status code: %d, but got %d", http.StatusCreated, response.Result().StatusCode)
	}
}

func Test_application_createAccountHandler_database_error(t *testing.T) {
	handler := http.HandlerFunc(app.CreateAccountHandler)
	reqBody := `{
					"owner": "Ruthvik",
					"balance": 20000,
					"Currency": "USD"
				}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/accounts/", strings.NewReader(reqBody))
	_createAccount := mock.CreateAccount
	defer func() {
		mock.CreateAccount = _createAccount
	}()
	{
		// mock db calls
		mock.CreateAccount = func(acc *models.Account) error {
			return errors.New("error")
		}
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status code: %d, but got %d", http.StatusInternalServerError, response.Result().StatusCode)
	}
}

func Test_application_createAccountHandler_badly_formed_json(t *testing.T) {
	reqBody := `{
					"owner": "Ruthvik"
					"balance": 20000,
					"Currency": "USD"
				}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/accounts/", strings.NewReader(reqBody))

	handler := http.HandlerFunc(app.CreateAccountHandler)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code: %d, but got %d", http.StatusBadRequest, response.Result().StatusCode)
	}
}

func Test_application_getAccountByIDHandler_success(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/accounts/", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	_getAccountByID := mock.GetAccountByID
	defer func() {
		mock.GetAccountByID = _getAccountByID
	}()
	{
		// mock calls to the db
		mock.GetAccountByID = func(id int64) (*models.Account, error) {
			toReturn := &models.Account{
				ID:        1,
				Owner:     "Ruthvik",
				Balance:   4000,
				Currency:  "USD",
				CreatedAt: time.Now(),
				Version:   0,
			}
			return toReturn, nil
		}
	}
	handler := http.HandlerFunc(app.getAccountByIDHandler)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Result().StatusCode != http.StatusOK {
		t.Errorf("expected status code: %d, but got %d", http.StatusOK, response.Result().StatusCode)
	}
}

func Test_application_getAccountByIDHandler_no_records_found(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/accounts/", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "10")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	_getAccountByID := mock.GetAccountByID
	defer func() {
		mock.GetAccountByID = _getAccountByID
	}()
	{
		// mock calls to db
		mock.GetAccountByID = func(id int64) (*models.Account, error) {
			return nil, store.ErrRecordNotFound
		}
	}
	handler := http.HandlerFunc(app.getAccountByIDHandler)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Result().StatusCode != http.StatusNotFound {
		t.Errorf("expected status code: %d, but got %d", http.StatusNotFound, response.Result().StatusCode)
	}
}

func Test_application_getAccountByIDHandler_database_error(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/accounts/", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	_getAccountByID := mock.GetAccountByID
	defer func() {
		mock.GetAccountByID = _getAccountByID
	}()
	{
		// mock calls to db
		mock.GetAccountByID = func(id int64) (*models.Account, error) {
			return nil, errors.New("error")
		}
	}
	handler := http.HandlerFunc(app.getAccountByIDHandler)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status code: %d, but got %d", http.StatusInternalServerError, response.Result().StatusCode)
	}
}