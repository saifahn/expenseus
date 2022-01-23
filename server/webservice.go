package expenseus

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"

	"golang.org/x/oauth2"
)

type contextKey int

const (
	CtxKeyExpenseID  contextKey = iota
	CtxKeyUsername   contextKey = iota
	CtxKeyUserID     contextKey = iota
	jsonContentType             = "application/json"
	SessionCookieKey            = "expenseus-session"
)

type ExpenseStore interface {
	GetExpense(id string) (Expense, error)
	GetExpensesByUsername(username string) ([]Expense, error)
	GetAllExpenses() ([]Expense, error)
	CreateExpense(expenseDetails ExpenseDetails) error
	CreateUser(user User) error
	GetUser(id string) (User, error)
	GetAllUsers() ([]User, error)
}

type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	ID       string `json:"id"`
}

type ExpenseDetails struct {
	Name     string `json:"name"`
	UserID   string `json:"userId"`
	ImageKey string `json:"imageKey,omitempty"`
}

type Expense struct {
	ExpenseDetails
	ID       string `json:"id"`
	ImageURL string `json:"imageUrl,omitempty"`
}

type ExpenseusOauth interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	GetInfoAndGenerateUser(state string, code string) (User, error)
}

type SessionManager interface {
	Validate(r *http.Request) bool
	Save(rw http.ResponseWriter, r *http.Request)
	GetUserID(r *http.Request) (string, error)
	Remove(rw http.ResponseWriter, r *http.Request)
}

type ImageStore interface {
	Upload(file multipart.File, header multipart.FileHeader) (string, error)
	Validate(file multipart.File) (bool, error)
	AddImageToExpense(expense Expense) (Expense, error)
}

type WebService struct {
	store       ExpenseStore
	oauthConfig ExpenseusOauth
	sessions    SessionManager
	images      ImageStore
	frontend    string
}

func NewWebService(store ExpenseStore, oauth ExpenseusOauth, sessions SessionManager, frontend string, images ImageStore) *WebService {
	return &WebService{store: store, oauthConfig: oauth, sessions: sessions, frontend: frontend, images: images}
}

// VerifyUser is middleware that checks that the user is logged in and authorized
// before passing the request to the handler with the userID in the context.
func (wb *WebService) VerifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		sessionIsAuthorized := wb.sessions.Validate(r)
		if !sessionIsAuthorized {
			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		id, err := wb.sessions.GetUserID(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		ctx := context.WithValue(r.Context(), CtxKeyUserID, id)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}

func (wb *WebService) OauthLogin(rw http.ResponseWriter, r *http.Request) {
	// TODO: add proper state string
	url := wb.oauthConfig.AuthCodeURL("")
	http.Redirect(rw, r, url, http.StatusTemporaryRedirect)
}

func (wb *WebService) OauthCallback(rw http.ResponseWriter, r *http.Request) {
	user, err := wb.oauthConfig.GetInfoAndGenerateUser(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: check by UserID instead?
	existingUsers, err := wb.store.GetAllUsers()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	// check if the user exists already
	for _, u := range existingUsers {
		if u.ID == user.ID {
			ctx := context.WithValue(r.Context(), CtxKeyUserID, u.ID)
			r = r.WithContext(ctx)
			wb.sessions.Save(rw, r)
			http.Redirect(rw, r, wb.frontend, http.StatusTemporaryRedirect)
			return
		}
	}

	// otherwise, create the user
	wb.store.CreateUser(user)
	ctx := context.WithValue(r.Context(), CtxKeyUserID, user.ID)
	r = r.WithContext(ctx)
	wb.sessions.Save(rw, r)
	http.Redirect(rw, r, wb.frontend, http.StatusTemporaryRedirect)
	// TODO: redirect to change username page
}

// GetExpense handles a HTTP request to get an expense by ID, returning the expense.
func (wb *WebService) GetExpense(rw http.ResponseWriter, r *http.Request) {
	expenseID := r.Context().Value(CtxKeyExpenseID).(string)

	expense, err := wb.store.GetExpense(expenseID)

	// TODO: should account for different kinds of errors
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
	}

	if expense.ImageKey != "" {
		expense, err = wb.images.AddImageToExpense(expense)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	rw.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(rw).Encode(expense)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetExpensesByUsername handles a HTTP request to get all expenses of a user,
// returning a list of expenses.
func (wb *WebService) GetExpensesByUsername(rw http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(CtxKeyUsername).(string)

	expenses, err := wb.store.GetExpensesByUsername(username)

	// TODO: account for different errors
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(rw).Encode(expenses)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetAllExpenses handles a HTTP request to get all expenses, returning a list
// of expenses.
func (wb *WebService) GetAllExpenses(rw http.ResponseWriter, r *http.Request) {
	expenses, err := wb.store.GetAllExpenses()

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, e := range expenses {
		if e.ImageKey != "" {
			expenses[i], err = wb.images.AddImageToExpense(e)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	rw.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(rw).Encode(expenses)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateExpense handles a HTTP request to create a new expense.
func (wb *WebService) CreateExpense(rw http.ResponseWriter, r *http.Request) {
	// get the userID from the context
	userID, ok := r.Context().Value(CtxKeyUserID).(string)
	if !ok {
		http.Error(rw, "user id not found in context", http.StatusUnauthorized)
	}

	err := r.ParseMultipartForm(1024 * 1024 * 5)
	if err != nil {
		if err == multipart.ErrMessageTooLarge {
			http.Error(rw, "image size too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	expenseName := r.FormValue("expenseName")
	if expenseName == "" {
		http.Error(rw, "expense name not found", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	// don't error on missing file - it's ok not to have an image
	if err != nil && err != http.ErrMissingFile {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var imageKey string
	// upload the image only if one was supplied
	if file != nil {
		// check image is OK
		ok, err := wb.images.Validate(file)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(rw, "image invalid", http.StatusUnprocessableEntity)
			return
		}

		imageKey, err = wb.images.Upload(file, *header)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = wb.store.CreateExpense(ExpenseDetails{Name: expenseName, UserID: userID, ImageKey: imageKey})

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}

// CreateUser handles a request to create a new user.
func (wb *WebService) CreateUser(rw http.ResponseWriter, r *http.Request) {
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = wb.store.CreateUser(u)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}

// ListUser handles a request to get all users and return the list of users.
func (wb *WebService) ListUsers(rw http.ResponseWriter, r *http.Request) {
	users, err := wb.store.GetAllUsers()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("content-type", jsonContentType)
	// TODO: return under a "users" key in JSON
	err = json.NewEncoder(rw).Encode(users)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetUser handles a HTTP request to get a user by ID, returning the user.
func (wb *WebService) GetUser(rw http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(CtxKeyUserID).(string)

	user, err := wb.store.GetUser(userID)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	rw.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(rw).Encode(user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetSelf handles a HTTP request to return the logged in user.
func (wb *WebService) GetSelf(rw http.ResponseWriter, r *http.Request) {
	id, err := wb.sessions.GetUserID(r)

	// TODO: add case for non-existent user
	// TODO: handle non-valid session
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	user, err := wb.store.GetUser(id)

	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
	}

	rw.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(rw).Encode(user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// LogOut handles a HTTP request to log out the current user.
func (wb *WebService) LogOut(rw http.ResponseWriter, r *http.Request) {
	wb.sessions.Remove(rw, r)

	http.Redirect(rw, r, wb.frontend, http.StatusTemporaryRedirect)
}
