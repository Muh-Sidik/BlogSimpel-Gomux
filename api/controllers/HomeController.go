package controllers

import (
	"net/http"

	"github.com/Muh-Sidik/BlogSimpel-Gomux/api/helpers/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Haloo sayang ow yaeh")
}
