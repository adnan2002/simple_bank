package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/db/mock"
	"example.com/db/sqlc"
	"example.com/db/util"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type TestCases struct {
	Name           string
	AccountId      interface{} // Changed to interface{} to support invalid IDs
	BuildStub      func(*mock.MockStore)
	CheckResponse  func(*testing.T, *httptest.ResponseRecorder)
}

func TestGetAccount(t *testing.T) {
	account := randomAccount()

	testCases := []TestCases{
		{
			Name:      "OK",
			AccountId: account.ID,
			BuildStub: func(ms *mock.MockStore) {
				ms.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			CheckResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusAccepted, rr.Code)
				requireBodyMatchAccount(t, rr, account)
			},
		},
		{
			Name:      "Not Found",
			AccountId: account.ID,
			BuildStub: func(ms *mock.MockStore) {
				ms.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, pgx.ErrNoRows)
			},
			CheckResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rr.Code)
				requireBodyMatchError(t, rr)
			},
		},
		{
			Name:      "Internal Server Error",
			AccountId: account.ID,
			BuildStub: func(ms *mock.MockStore) {
				ms.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, fmt.Errorf("database connection failed"))
			},
			CheckResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rr.Code)
				requireBodyMatchError(t, rr)
			},
		},
		{
			Name:      "Bad Request - Invalid ID",
			AccountId: "invalid", // This will cause ShouldBindUri to fail
			BuildStub: func(ms *mock.MockStore) {
				// No expectations since the request should fail before reaching the store
			},
			CheckResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rr.Code)
				requireBodyMatchError(t, rr)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mock.NewMockStore(ctrl)
			tc.BuildStub(mockStore)

			server := NewServer(mockStore)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%v", tc.AccountId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.CheckResponse(t, recorder)
		})
	}
}

func requireBodyMatchAccount(t *testing.T, recorder *httptest.ResponseRecorder, account db.Account) {
	var body db.Account
	err := json.Unmarshal(recorder.Body.Bytes(), &body)
	require.NoError(t, err)
	require.Equal(t, account, body)
}

func requireBodyMatchError(t *testing.T, recorder *httptest.ResponseRecorder) {
	// Just verify that the response body contains an error message
	// The exact structure depends on your errorResponse function
	var body map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &body)
	require.NoError(t, err)
	require.NotEmpty(t, body)
}

func randomAccount() db.Account {
	return db.Account{
		ID:    util.RandomInt(1, 10),
		Owner: util.RandomOwner(),
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}