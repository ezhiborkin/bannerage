package handler

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log/slog"
	"net/http"
	"strings"
)

type Handler struct {
	log            *slog.Logger
	bannerProvider BannerProvider
	userProvider   UserProvider
	authProvider   AuthProvider
	context        context.Context
}

func New(log *slog.Logger,
	userProvider UserProvider,
	bannerProvider BannerProvider,
	authProvider AuthProvider,
	context context.Context,
) (*Handler, error) {
	return &Handler{
		log:            log,
		userProvider:   userProvider,
		bannerProvider: bannerProvider,
		authProvider:   authProvider,
		context:        context,
	}, nil
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", h.loginUser)
	mux.HandleFunc("POST /create/user", h.createUser)

	mux.HandleFunc("POST /banner", adminMiddleware(http.HandlerFunc(h.postBanner)))
	mux.HandleFunc("GET /banner", authMiddleware(http.HandlerFunc(h.listBanners)))

	mux.HandleFunc("GET /user_banner", authMiddleware(http.HandlerFunc(h.getUserBanner)))

	mux.HandleFunc("POST /choose_revision", adminMiddleware(http.HandlerFunc(h.chooseBanner)))

	mux.HandleFunc("GET /banner_revisions/{banner_id}", adminMiddleware(http.HandlerFunc(h.listRevisions)))

	mux.HandleFunc("DELETE /banner/{id}", adminMiddleware(http.HandlerFunc(h.deleteBanner)))
	mux.HandleFunc("PATCH /banner/{id}", adminMiddleware(http.HandlerFunc(h.patchBanner)))

	mux.HandleFunc("DELETE /banner_deferred", adminMiddleware(h.deleteBannerFeatureTag(h.context)))

	return mux
}

func authMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			handleUnauthorized(w, "Authorization header missing")
			return
		}

		tokenString := extractTokenString(authHeader)
		token, err := parseToken(tokenString)
		if err != nil {
			handleUnauthorized(w, "Error parsing token: "+err.Error())
			return
		}

		if !token.Valid {
			handleUnauthorized(w, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			handleUnauthorized(w, "Token claims are not in the expected format")
			return
		}

		role, ok := claims["Role"].(string)
		if !ok {
			handleUnauthorized(w, "Role claim not found or not a string")
			return
		}

		ctx := context.WithValue(r.Context(), "role", role)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func adminMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			handleUnauthorized(w, "Authorization header missing")
			return
		}

		tokenString := extractTokenString(authHeader)
		token, err := parseToken(tokenString)
		if err != nil {
			handleUnauthorized(w, "Error parsing token: "+err.Error())
			return
		}

		if !token.Valid {
			handleUnauthorized(w, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			handleUnauthorized(w, "Token claims are not in the expected format")
			return
		}

		role, ok := claims["Role"].(string)
		if !ok {
			handleUnauthorized(w, "Role claim not found or not a string")
			return
		}

		if role != "admin" {
			handleUnauthorized(w, "Wrong role")
			return
		}

		ctx := context.WithValue(r.Context(), "role", role)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func handleUnauthorized(w http.ResponseWriter, message string) {
	fmt.Println(message)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(message))
}

func extractTokenString(authHeader string) string {
	return strings.Replace(authHeader, "Bearer ", "", 1)
}

func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
}
