package middlewares

import (
	"fmt"
	"net/http"

	"github.com/swethabhageerath/logservice/internal/models"
)

type RequestHandlerFunc func(w http.ResponseWriter, r *http.Request) models.StandardResponse

func RequestHandler(h RequestHandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := h(w, r)
		if response.Err != nil {
			fmt.Println(response.Err.Error())
		} else {
			fmt.Println(response.Data)
		}
	})
}
