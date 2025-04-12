package middlewares

import (
	"context"
	"github.com/go-chi/render"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"net/http"
	"os"
	"strings"
	"time"
	"tpu-practice-searcher/internal/utils"
)

type contextKey string

const (
	_initDataKey contextKey = "init-data"
)

func withInitData(ctx context.Context, initData initdata.InitData) context.Context {
	return context.WithValue(ctx, _initDataKey, initData)
}

func CtxInitData(ctx context.Context) (initdata.InitData, bool) {
	initData, ok := ctx.Value(_initDataKey).(initdata.InitData)
	return initData, ok
}

func CtxClaims(ctx context.Context) (*utils.Claims, bool) {
	claims, ok := ctx.Value("claims").(*utils.Claims)
	return claims, ok
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := os.Getenv("TELEGRAM_BOT_TOKEN")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("Authorization header required"))

			return
		}

		authParts := strings.Split(authHeader, " ")

		authType := authParts[0]
		authData := authParts[1]

		switch authType {
		case "tma":
			if err := initdata.Validate(authData, token, time.Hour*72); err != nil {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, utils.NewErrorResponse("Unauthorized"))

				return
			}

			initData, err := initdata.Parse(authData)
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, utils.NewErrorResponse("Internal server error"))

				return
			}

			ctx := withInitData(r.Context(), initData)
			next.ServeHTTP(w, r.WithContext(ctx))

		case "Bearer":
			claims, err := utils.ValidateAccessToken(authData)
			if err != nil {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, utils.NewErrorResponse("Invalid access token"))

				return
			}

			ctx := context.WithValue(r.Context(), "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
