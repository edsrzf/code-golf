package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/code-golf/code-golf/golfer"
)

var golferKey = key("golfer")
var hmacKey []byte

func init() {
	var err error
	if hmacKey, err = base64.RawURLEncoding.DecodeString(os.Getenv("HMAC_KEY")); err != nil {
		panic(err)
	}
}

func GolferCookie(data string) string {
	mac := hmac.New(sha256.New, hmacKey)
	mac.Write([]byte(data))

	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// GolferHandler adds the golfer to the context if logged in.
func GolferHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie("__Host-user"); err == nil {
			if i := strings.LastIndexByte(cookie.Value, ':'); i != -1 {
				mac := hmac.New(sha256.New, hmacKey)
				mac.Write([]byte(cookie.Value[:i]))

				if subtle.ConstantTimeCompare(
					[]byte(cookie.Value[i+1:]),
					[]byte(base64.RawURLEncoding.EncodeToString(mac.Sum(nil))),
				) == 1 {
					var golfer golfer.Golfer

					j := strings.IndexByte(cookie.Value, ':')
					golfer.ID, _ = strconv.Atoi(cookie.Value[:j])
					golfer.Name = cookie.Value[j+1 : i]

					// TODO
					golfer.Admin = golfer.Name == "JRaspass"

					r = r.WithContext(
						context.WithValue(r.Context(), golferKey, &golfer))
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

// Golfer gets the golfer from the request context.
func Golfer(r *http.Request) *golfer.Golfer {
	golfer, _ := r.Context().Value(golferKey).(*golfer.Golfer)
	return golfer
}
