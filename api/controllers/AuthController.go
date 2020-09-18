package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Muh-Sidik/BlogSimpel-Gomux/api/helpers/auth"
	"github.com/Muh-Sidik/BlogSimpel-Gomux/api/helpers/responses"
	"github.com/Muh-Sidik/BlogSimpel-Gomux/api/helpers/utils/formaterror"
	"github.com/Muh-Sidik/BlogSimpel-Gomux/api/models"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}

	err = json.Unmarshal(body, &user)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()

	err = user.Validate("login")

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	tokenUser, err := server.FindUser(user.Email, user.Password)

	if err != nil {
		formatError := formaterror.FormatError(err.Error())

		responses.ERROR(w, http.StatusBadRequest, formatError)
		return
	}

	responses.JSON(w, http.StatusOK, tokenUser)
}

func (server *Server) FindUser(email, password string) (string, error) {
	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error

	if err != nil {
		return "", err
	}

	err = models.VerifyPassword(user.Password, password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	auth.CreateToken(user.ID)
}
