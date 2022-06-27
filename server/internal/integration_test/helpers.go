package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/gorilla/securecookie"
	"github.com/saifahn/expenseus/internal/app"
	mock_app "github.com/saifahn/expenseus/internal/app/mocks"
	"github.com/saifahn/expenseus/internal/ddb"
	"github.com/saifahn/expenseus/internal/router"
	"github.com/saifahn/expenseus/internal/sessions"
	"github.com/stretchr/testify/assert"
)

var (
	TestSeanUser = app.User{
		Username: "saifahn",
		Name:     "Sean Li",
		ID:       "sean_id",
	}
	TestTomomiUser = app.User{
		Username: "tomochi",
		Name:     "Tomomi Kinoshita",
		ID:       "tomomi_id",
	}

	TestSeanTxnDetails = app.Transaction{
		Location: "Location 1",
		UserID:   TestSeanUser.ID,
		Amount:   123,
		Date:     1644085875,
		Category: "test.test",
	}

	TestTracker = app.Tracker{
		Name:  "Test Tracker",
		Users: []string{TestSeanUser.ID},
		ID:    "test-id",
	}

	testSessionHashKey  = securecookie.GenerateRandomKey(64)
	testSessionRangeKey = securecookie.GenerateRandomKey(32)
	cookies             = securecookie.New(testSessionHashKey, testSessionRangeKey)
)

const (
	testTableName = "expenseus-integ-test"
)

func SetUpDB(d dynamodbiface.DynamoDBAPI) (app.Store, error) {
	err := ddb.CreateTable(d, testTableName)
	if err != nil {
		return nil, err
	}

	return ddb.New(d, testTableName), nil
}

func TearDownDB(d dynamodbiface.DynamoDBAPI) error {
	err := ddb.DeleteTable(d, testTableName)
	if err != nil {
		return err
	}

	return nil
}

// SetUpTestServer sets up a server with with the real routes and a test
// dynamodb instance, with stubs for the rest of the app
func SetUpTestServer(t *testing.T) (http.Handler, func(t *testing.T)) {
	ddbLocal := ddb.NewDynamoDBLocalAPI()
	db, err := SetUpDB(ddbLocal)
	if err != nil {
		t.Fatalf("could not set up the database: %v", err)
	}

	oauth := &mock_app.MockAuth{}
	session := sessions.New(testSessionHashKey, testSessionRangeKey)
	images := &mock_app.MockImageStore{}
	a := app.New(db, oauth, session, "", images)
	r := router.Init(a)

	return r, func(t *testing.T) {
		err := TearDownDB(ddbLocal)
		if err != nil {
			t.Fatalf("could not tear down the database: %v", err)
		}
	}
}

// CreateCookie uses the same keys as the session manager provided for the
// integration tests to encode a value and provide it in a cookie for the tests
func CreateCookie(userID string) *http.Cookie {
	encoded, err := cookies.Encode(app.SessionCookieKey, userID)
	if err != nil {
		panic(err)
	}
	return &http.Cookie{
		Name:  app.SessionCookieKey,
		Value: encoded,
	}
}

// CreateUser creates a user
func CreateUser(t *testing.T, user app.User, r http.Handler) {
	userJSON, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("failed to marshal the user JSON: %v", err)
	}
	response := httptest.NewRecorder()
	request := app.NewCreateUserRequest(userJSON)
	request.AddCookie(CreateCookie(user.ID))
	r.ServeHTTP(response, request)
	assert.Equal(t, http.StatusAccepted, response.Code)
}

// AssertEqualTxnDetails compares two transactions without their ID as these are
// generated by the database
func AssertEqualTxnDetails(t testing.TB, want, got app.Transaction) {
	wantWithoutID := app.Transaction{
		Location: want.Location,
		UserID:   want.UserID,
		Amount:   want.Amount,
		Date:     want.Date,
		Category: want.Category,
		Details:  want.Details,
	}
	gotWithoutID := app.Transaction{
		Location: got.Location,
		UserID:   got.UserID,
		Amount:   got.Amount,
		Date:     got.Date,
		Category: got.Category,
		Details:  got.Details,
	}
	assert.Equal(t, wantWithoutID, gotWithoutID)
}

// RemoveSharedTxnID returns a  shared transaction without the ID
func RemoveSharedTxnID(txn app.SharedTransaction) app.SharedTransaction {
	txn.ID = ""
	return txn
}

func CreateTestTxn(t *testing.T, r http.Handler, td app.Transaction, userid string) {
	payload := app.MakeTxnRequestPayload(td)
	request := app.NewCreateTransactionRequest(payload)
	request.AddCookie(CreateCookie(userid))
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)
}
