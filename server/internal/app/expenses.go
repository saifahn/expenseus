package app

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
)

const CtxKeyExpenseID contextKey = iota

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

// GetExpense handles a HTTP request to get an expense by ID, returning the expense.
func (a *App) GetExpense(rw http.ResponseWriter, r *http.Request) {
	expenseID := r.Context().Value(CtxKeyExpenseID).(string)

	expense, err := a.store.GetExpense(expenseID)

	// TODO: should account for different kinds of errors
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
	}

	if expense.ImageKey != "" {
		expense, err = a.images.AddImageToExpense(expense)
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
func (a *App) GetExpensesByUsername(rw http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(CtxKeyUsername).(string)

	expenses, err := a.store.GetExpensesByUsername(username)

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
func (a *App) GetAllExpenses(rw http.ResponseWriter, r *http.Request) {
	expenses, err := a.store.GetAllExpenses()

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, e := range expenses {
		if e.ImageKey != "" {
			expenses[i], err = a.images.AddImageToExpense(e)
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
func (a *App) CreateExpense(rw http.ResponseWriter, r *http.Request) {
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
		ok, err := a.images.Validate(file)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(rw, "image invalid", http.StatusUnprocessableEntity)
			return
		}

		imageKey, err = a.images.Upload(file, *header)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = a.store.CreateExpense(ExpenseDetails{Name: expenseName, UserID: userID, ImageKey: imageKey})

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}
