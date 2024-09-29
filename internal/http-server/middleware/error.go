package middleware

import (
	"blog/internal/constants/servererror"
	resp "blog/internal/lib/api/response"
	"net/http"

	"github.com/go-chi/render"
)

func ErrorHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                switch err := err.(type) {
                case error:
                    switch err {
                    case servererror.InvalidJson:
                        render.JSON(w, r, resp.Error("failed to decode JSON", http.StatusBadRequest))
                    case servererror.ResourceNotFound:
                        render.JSON(w, r, resp.Error("resource not found", http.StatusNotFound))
                    case servererror.InternalError:
                        render.JSON(w, r, resp.Error("internal server error", http.StatusInternalServerError))
                    case servererror.InvalidCrerdentials:
                        render.JSON(w, r, resp.Error("invalid password or username", http.StatusUnauthorized))
                    default:
                        render.JSON(w, r, resp.Error("unexcepted error", http.StatusInternalServerError))
                    }
                } 
            }
        }()
        next.ServeHTTP(w, r)
    })
}
