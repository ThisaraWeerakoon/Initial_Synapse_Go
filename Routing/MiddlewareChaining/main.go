package main

import (
        "context"
        "encoding/json"
        "fmt"
        "log"
        "net/http"
        "strings"

        "github.com/golang-jwt/jwt/v5" // Using v5 for latest features
)

// Example JWT secret key (replace with a secure key in production)
var jwtSecret = []byte("your-secret-key")

// Example data structure for the API response
type ApiResponse struct {
        Message string `json:"message"`
        Data    interface{} `json:"data,omitempty"`
        Error   string `json:"error,omitempty"`
}

// Example data structure for the request
type ApiRequest struct {
        Payload string `json:"payload"`
}

// Middleware to validate JWT
func jwtAuthMiddleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                authHeader := r.Header.Get("Authorization")
                if authHeader == "" {
                        http.Error(w, "Authorization header required", http.StatusUnauthorized)
                        return
                }

                tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

                token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
                        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
                        }
                        return jwtSecret, nil
                })

                if err != nil {
                        http.Error(w, "Invalid token", http.StatusUnauthorized)
                        return
                }

                if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
                        ctx := context.WithValue(r.Context(), "claims", claims)
                        next.ServeHTTP(w, r.WithContext(ctx))
                } else {
                        http.Error(w, "Invalid token claims", http.StatusUnauthorized)
                        return
                }
        })
}

// Middleware for standard request validation
func requestValidationMiddleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                if r.Method != http.MethodPost {
                        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                        return
                }

                var req ApiRequest
                if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                        http.Error(w, "Invalid request body", http.StatusBadRequest)
                        return
                }

                if req.Payload == "" {
                        http.Error(w, "Payload is required", http.StatusBadRequest)
                        return
                }

                ctx := context.WithValue(r.Context(), "request", req)
                next.ServeHTTP(w, r.WithContext(ctx))
        })
}

// Middleware for response validation
func responseValidationMiddleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

                originalWriter := w
                wrappedWriter := &responseWriterInterceptor{ResponseWriter: w}

                next.ServeHTTP(wrappedWriter, r)

                contentType := wrappedWriter.Header().Get("Content-Type")

                if strings.Contains(contentType, "application/json") {
                        var response ApiResponse

                        if err := json.Unmarshal(wrappedWriter.body, &response); err != nil {
                                log.Printf("Response Validation Error: %v", err)
                                http.Error(originalWriter, "Internal Server Error", http.StatusInternalServerError)
                                return
                        }
                        if response.Error != "" {
                                log.Printf("API returned an error: %s", response.Error)
                        }

                        originalWriter.Write(wrappedWriter.body)

                } else {
                        originalWriter.Write(wrappedWriter.body)
                }

        })
}

type responseWriterInterceptor struct {
        http.ResponseWriter
        body []byte
}

func (w *responseWriterInterceptor) Write(b []byte) (int, error) {
        w.body = b
        return len(b), nil

}

// Example API handler
func apiHandler(w http.ResponseWriter, r *http.Request) {
        req := r.Context().Value("request").(ApiRequest)

        response := ApiResponse{
                Message: "API Success",
                Data:    fmt.Sprintf("Processed payload: %s", req.Payload),
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
}

func main() {
        handler := http.HandlerFunc(apiHandler)
        http.Handle("/api", jwtAuthMiddleware(requestValidationMiddleware(responseValidationMiddleware(handler))))

        fmt.Println("Server listening on :8080")
        if err := http.ListenAndServe(":8080", nil); err != nil {
                log.Fatal(err)
        }
}