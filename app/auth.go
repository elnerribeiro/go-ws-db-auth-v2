package app

import (
	"context"
	"net/http"
	"strings"

	repo "github.com/elnerribeiro/go-ws-db-auth-v2/repositories"
	u "github.com/elnerribeiro/go-ws-db-auth-v2/utils"

	"github.com/golang-jwt/jwt/v5"
)

// JwtAuthentication Auth with JWT
var JwtAuthentication = func(next http.Handler) http.Handler {
	logger, logContext := u.GetLoggerAndContext()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/login"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path         //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			logger.Info().Msg("[JwtAuthentication] 403 Token Not found!")
			response = u.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(logContext, w, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			logger.Info().Msg("[JwtAuthentication] 403 Invalid/Malformed auth token!")
			response = u.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(logContext, w, response)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in

		token, err := jwt.ParseWithClaims(tokenPart, &repo.Token{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("JWTpassword123@"), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			logger.Info().Msg("[JwtAuthentication] 403 Invalid, expired or malformed token!")
			response = u.Message(false, "Invalid, expired or malformed token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(logContext, w, response)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			logger.Info().Msg("[JwtAuthentication] 403 Token not valid on this server!")
			response = u.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(logContext, w, response)
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		parsedToken := token.Claims.(*repo.Token)
		logger.Info().Msgf("User %d just logged in", parsedToken.UserID) //Useful for monitoring
		var ctx = context.WithValue(r.Context(), repo.ContextKey("user"), parsedToken.UserID)
		ctx = context.WithValue(ctx, repo.ContextKey("role"), parsedToken.Role)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}
