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

func Test_application_updateAccountHandler_success(t *testing.T) {
	reqBody := `
	{
		"id":1,
		"owner":"R",
		"balance": 2000,
		"currency":"EUR"
	}
	`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/accounts/", strings.NewReader(reqBody))
	_updateAccount := mock.UpdateAccount
	defer func() {
		mock.UpdateAccount = _updateAccount
	}()
	{
		// mock calls to db
		mock.UpdateAccount = func(acc *models.Account) error {
			return nil
		}
	}
	handler := http.HandlerFunc(app.updateAccountHandler)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)
	if response.Result().StatusCode != http.StatusOK {
		t.Errorf("expected status code: %d, but got %d", http.StatusOK, response.Result().StatusCode)
	}
}

func Test_application_updateAccountHandler_badly_input(t *testing.T) {
	reqBody := `
	{
		"id":1,
		"owner":"R",
		"balance": "2000",
		"currency":"EUR"
	}
	`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/accounts/", strings.NewReader(reqBody))
	_updateAccount := mock.UpdateAccount
	defer func() {
		mock.UpdateAccount = _updateAccount
	}()
	{
		// mock calls to db
		mock.UpdateAccount = func(acc *models.Account) error {
			return nil
		}
	}
	handler := http.HandlerFunc(app.updateAccountHandler)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)
	if response.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code: %d, but got %d", http.StatusBadRequest, response.Result().StatusCode)
	}
}

func Test_application_updateAccountHandler_record_not_found(t *testing.T) {
	reqBody := `
	{
		"id":1,
		"owner":"R",
		"balance": 2000,
		"currency":"EUR"
	}
	`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/accounts/", strings.NewReader(reqBody))
	_updateAccount := mock.UpdateAccount
	defer func() {
		mock.UpdateAccount = _updateAccount
	}()
	{
		// mock calls to db
		mock.UpdateAccount = func(acc *models.Account) error {
			return store.ErrRecordNotFound
		}
	}
	handler := http.HandlerFunc(app.updateAccountHandler)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)
	if response.Result().StatusCode != http.StatusNotFound {
		t.Errorf("expected status code: %d, but got %d", http.StatusNotFound, response.Result().StatusCode)
	}
}

func Test_application_updateAccountHandler_database_error(t *testing.T) {
	reqBody := `
	{
		"id":1,
		"owner":"R",
		"balance": 2000,
		"currency":"EUR"
	}
	`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/accounts/", strings.NewReader(reqBody))
	_updateAccount := mock.UpdateAccount
	defer func() {
		mock.UpdateAccount = _updateAccount
	}()
	{
		// mock calls to db
		mock.UpdateAccount = func(acc *models.Account) error {
			return errors.New("error")
		}
	}
	handler := http.HandlerFunc(app.updateAccountHandler)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)
	if response.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status code: %d, but got %d", http.StatusInternalServerError, response.Result().StatusCode)
	}
}

func Test_application_updateBalanceHandler_success(t *testing.T) {
	reqBody := `{
		"balance": 5000
	}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/accounts", strings.NewReader(reqBody))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	_getAccountByID := mock.GetAccountByID
	_updateBalance := mock.UpdateBalance
	defer func() {
		mock.GetAccountByID = _getAccountByID
		mock.UpdateBalance = _updateBalance
	}()
	{
		// mock calls to db
		mock.GetAccountByID = func(id int64) (*models.Account, error) {
			toReturn := &models.Account{
				ID:        1,
				Owner:     "Ruthvik",
				Balance:   4000,
				Currency:  "USD",
				CreatedAt: time.Now(),
			}
			return toReturn, nil
		}

		mock.UpdateBalance = func(acc *models.Account) error {
			return nil
		}
	}
	handler := http.HandlerFunc(app.updateBalanceHandler)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)
	if response.Result().StatusCode != http.StatusOK {
		t.Errorf("expected status code: %d, but got %d", http.StatusOK, response.Result().StatusCode)
	}
}

func Test_application_updateBalanceHandler_bad_input(t *testing.T) {
	reqBody := `{
		"balance": "5000"
	}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/accounts", strings.NewReader(reqBody))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	_getAccountByID := mock.GetAccountByID
	_updateBalance := mock.UpdateBalance
	defer func() {
		mock.GetAccountByID = _getAccountByID
		mock.UpdateBalance = _updateBalance
	}()
	handler := http.HandlerFunc(app.updateBalanceHandler)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)
	if response.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code: %d, but got %d", http.StatusBadRequest, response.Result().StatusCode)
	}
}

func Test_application_updateBalanceHandler_record_not_found(t *testing.T) {
	reqBody := `{
		"balance": 5000
	}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/accounts", strings.NewReader(reqBody))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	_getAccountByID := mock.GetAccountByID
	_updateBalance := mock.UpdateBalance
	defer func() {
		mock.GetAccountByID = _getAccountByID
		mock.UpdateBalance = _updateBalance
	}()
	{
		// mock calls to db
		mock.GetAccountByID = func(id int64) (*models.Account, error) {
			return nil, store.ErrRecordNotFound
		}

	}
	handler := http.HandlerFunc(app.updateBalanceHandler)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)
	if response.Result().StatusCode != http.StatusNotFound {
		t.Errorf("expected status code: %d, but got %d", http.StatusNotFound, response.Result().StatusCode)
	}
}
func Test_application_updateBalanceHandler_fetch_account_database_error(t *testing.T) {
	reqBody := `{
		"balance": 5000
	}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/accounts", strings.NewReader(reqBody))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	_getAccountByID := mock.GetAccountByID
	_updateBalance := mock.UpdateBalance
	defer func() {
		mock.GetAccountByID = _getAccountByID
		mock.UpdateBalance = _updateBalance
	}()
	{
		// mock calls to db
		mock.GetAccountByID = func(id int64) (*models.Account, error) {
			return nil, errors.New("error")
		}
	}
	handler := http.HandlerFunc(app.updateBalanceHandler)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)
	if response.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status code: %d, but got %d", http.StatusInternalServerError, response.Result().StatusCode)
	}
}

func Test_application_updateBalanceHandler_update_balance_database_error(t *testing.T) {
	reqBody := `{
		"balance": 5000
	}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/accounts", strings.NewReader(reqBody))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	_getAccountByID := mock.GetAccountByID
	_updateBalance := mock.UpdateBalance
	defer func() {
		mock.GetAccountByID = _getAccountByID
		mock.UpdateBalance = _updateBalance
	}()
	{
		// mock calls to db
		mock.GetAccountByID = func(id int64) (*models.Account, error) {
			toReturn := &models.Account{
				ID:        1,
				Owner:     "Ruthvik",
				Balance:   4000,
				Currency:  "USD",
				CreatedAt: time.Now(),
			}
			return toReturn, nil
		}

		mock.UpdateBalance = func(acc *models.Account) error {
			return errors.New("error")
		}
	}
	handler := http.HandlerFunc(app.updateBalanceHandler)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)
	if response.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status code: %d, but got %d", http.StatusInternalServerError, response.Result().StatusCode)
	}
}

func Test_application_listAccountsHandler_success(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/accounts/", nil)
	_listAccounts := mock.ListAccounts
	defer func() {
		mock.ListAccounts = _listAccounts
	}()
	{
		// mock calls to db
		mock.ListAccounts = func() ([]*models.Account, error) {
			toReturn := []*models.Account{
				{
					ID:        1,
					Owner:     "Ruthvik",
					Balance:   4000,
					Currency:  "USD",
					CreatedAt: time.Now(),
				},
			}
			return toReturn, nil
		}
	}
	handler := http.HandlerFunc(app.listAccountsHandler)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)
	if response.Result().StatusCode != http.StatusOK {
		t.Errorf("expected status code: %d, but got %d", http.StatusOK, response.Result().StatusCode)
	}
}

func Test_application_listAccountsHandler_database_error(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/accounts/", nil)
	_listAccounts := mock.ListAccounts
	defer func() {
		mock.ListAccounts = _listAccounts
	}()
	{
		// mock calls to db
		mock.ListAccounts = func() ([]*models.Account, error) {
			return nil, errors.New("error")
		}
	}
	handler := http.HandlerFunc(app.listAccountsHandler)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)
	if response.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status code: %d, but got %d", http.StatusInternalServerError, response.Result().StatusCode)
	}
}

func Test_application_deleteAccountHandler_success(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/accounts/", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	_deleteAccount := mock.DeleteAccount
	defer func() {
		mock.DeleteAccount = _deleteAccount
	}()
	{
		// mock calls to db
		mock.DeleteAccount = func(id int64) error {
			return nil
		}
	}
	handler := http.HandlerFunc(app.deleteAccountHandler)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Result().StatusCode != http.StatusOK {
		t.Errorf("expected status code: %d, but got %d", http.StatusOK, response.Result().StatusCode)
	}
}

func Test_application_deleteAccountHandler_database_error(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/accounts/", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	_deleteAccount := mock.DeleteAccount
	defer func() {
		mock.DeleteAccount = _deleteAccount
	}()
	{
		// mock calls to db
		mock.DeleteAccount = func(id int64) error {
			return errors.New("error")
		}
	}
	handler := http.HandlerFunc(app.deleteAccountHandler)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status code: %d, but got %d", http.StatusInternalServerError, response.Result().StatusCode)
	}
}

func Test_application_deleteAccountHandler_record_not_found(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/accounts/", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	_deleteAccount := mock.DeleteAccount
	defer func() {
		mock.DeleteAccount = _deleteAccount
	}()
	{
		// mock calls to db
		mock.DeleteAccount = func(id int64) error {
			return store.ErrRecordNotFound
		}
	}
	handler := http.HandlerFunc(app.deleteAccountHandler)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Result().StatusCode != http.StatusNotFound {
		t.Errorf("expected status code: %d, but got %d", http.StatusNotFound, response.Result().StatusCode)
	}
}
