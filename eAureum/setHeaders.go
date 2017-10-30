// Sets http headers
package eAureumFV

import (
	"net/http"
)

func setHeaders(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Set("Strict-Transport-Security", "max-age=63072000")

	w.Header().Set("Content-Security-Policy", "default-src https:; font-src 'self'; object-src 'none'; frame-src 'self'; frame-ancestors 'self'; base-uri 'none'; script-src 'self'; style-src 'self' 'unsafe-inline'")
}
