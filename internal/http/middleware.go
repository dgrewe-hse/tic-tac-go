// Copyright 2025 Esslingen University of Applied Sciences
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Dennis Grewe
// Version: 1.0.0
// Date: 2025-12-03

package http

import (
	"context"
	"net/http"
)

// contextKey is a private type to avoid collisions in context.
type contextKey string

const playerIDContextKey contextKey = "playerID"

// WithPlayerID is middleware that extracts X-Player-Id and stores it in the request context.
func WithPlayerID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		playerID := r.Header.Get("X-Player-Id")
		if playerID != "" {
			ctx := context.WithValue(r.Context(), playerIDContextKey, playerID)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

// PlayerIDFromContext retrieves the player ID from context if present.
func PlayerIDFromContext(ctx context.Context) string {
	value := ctx.Value(playerIDContextKey)
	if v, ok := value.(string); ok {
		return v
	}
	return ""
}
