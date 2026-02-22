package middleware

import (
	"context"
	"net/http"
)

const UserIDKey = "user_id"

func GetUserIDFromContext(r *http.Request) string {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}

func SetUserIDInContext(r *http.Request, userID string) *http.Request {
	ctx := context.WithValue(r.Context(), UserIDKey, userID)
	return r.WithContext(ctx)
}
