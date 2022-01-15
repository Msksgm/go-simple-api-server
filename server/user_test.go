package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/msksgm/go-simple-api-server/mock"
	"github.com/msksgm/go-simple-api-server/model"
)

func Test_createUser(t *testing.T) {
	userStore := &mock.UserService{}
	srv := testServer()
	srv.userService = userStore

	input := `{
		"user": {
			"username": "username",
			"password": "password"
		}
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(input))
	w := httptest.NewRecorder()

	var user model.User
	userStore.CreateUserFn = func(u *model.User) error {
		user = *u
		return nil
	}

	srv.router.ServeHTTP(w, req)
	expectedResp := userResponse(&user)
	gotResp := M{}
	extractResponseUserBody(w.Body, &gotResp)

	if code := w.Code; code != http.StatusCreated {
		t.Errorf("expected status code of 201, but got %d", code)
	}

	if !reflect.DeepEqual(expectedResp, gotResp) {
		t.Errorf("expected response %v, but got %v", expectedResp, gotResp)
	}
}

func Test_loginUser(t *testing.T) {
	userStore := &mock.UserService{}
	srv := testServer()
	srv.userService = userStore

	userStore.AuthenticateFn = func() *model.User {
		user := &model.User{
			Username: "username",
		}
		return user
	}

	input := `{
		"user": {
			"username": "username",
			"password": "password"
		}
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", strings.NewReader(input))
	w := httptest.NewRecorder()

	srv.router.ServeHTTP(w, req)

	if code := w.Code; code != http.StatusOK {
		t.Errorf("expected status code of 200, but got %d", code)
	}
}

func extractResponseUserBody(body io.Reader, v interface{}) {
	mm := M{}
	_ = readJSON(body, &mm)
	byt, err := json.Marshal(mm["user"])
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(byt, v)
	if err != nil {
		panic(err)
	}
}
