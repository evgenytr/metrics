// Package middleware contains custom middleware for metrics server service.
package middleware

import (
	"fmt"
	"net"
	"net/http"
)

// WithSubnetCheck wraps handler with middleware to check whether request comes from allowed subnet.
func WithSubnetCheck(ipNet *net.IPNet) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		checkFn := func(res http.ResponseWriter, req *http.Request) {
			reqIP := req.Header["X-Real-Ip"]
			if len(reqIP) == 0 {
				res.WriteHeader(http.StatusForbidden)
				fmt.Println("IP not allowed")
				return
			}

			parsedIP := net.ParseIP(reqIP[0])
			if parsedIP == nil || !ipNet.Contains(parsedIP) {
				fmt.Println("IP not allowed")
				res.WriteHeader(http.StatusForbidden)
				return
			}
			fmt.Println("IP allowed")
			next.ServeHTTP(res, req)
		}
		return http.HandlerFunc(checkFn)
	}
}
