package server

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/msksgm/go-simple-api-server/model"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fid reflect.StructField) string {
		name := strings.SplitN(fid.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			name = ""
		}
		return name
	})
}

func userResponse(user *model.User, _token ...string) M {
	if user == nil {
		return nil
	}
	return M{
		"username": user.Username,
	}
}

func (s *Server) createUser() http.HandlerFunc {
	type Input struct {
		User struct {
			Username string `json:"username" validate:"required,min=2"`
			Password string `json:"password" validate:"required,min=8,max=72"`
		} `json:"user" validate:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		input := &Input{}

		if err := readJSON(r.Body, &input); err != nil {
			errorResponse(w, http.StatusUnprocessableEntity, err)
			return
		}

		if err := validate.Struct(input.User); err != nil {
			validationError(w, err)
			return
		}

		user := model.User{
			Username: input.User.Username,
		}

		if err := user.SetPassword(input.User.Password); err != nil {
			panic(err)
		}

		if err := s.userService.CreateUser(r.Context(), &user); err != nil {
			switch {
			case errors.Is(err, model.ErrDuplicateUsername):
				err = ErrorM{"username": []string{"this username is already in use"}}
				errorResponse(w, http.StatusConflict, err)
			default:
				serverError(w, err)
			}
			return
		}

		writeJSON(w, http.StatusCreated, M{"user": user})
	}
}
