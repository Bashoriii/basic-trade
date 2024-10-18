package middleware

import (
	"basic-trade/helpers"
	"context"
	"net/http"
)

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		verifyToken, err := helpers.VerifyToken(r)
		if err != nil {
			http.Error(w, "Unauthenticated: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Set the user data in the request context
		ctx := context.WithValue(r.Context(), "userData", verifyToken)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
