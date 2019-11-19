package utils

import (
    "net/http"
    "reflect"
)

var client http.Client = HTTPClient(SetProxy())

// CheckType check var type
func CheckType(i interface{}) reflect.Type {
    return reflect.TypeOf(i)
}

// CheckNil *interface
func CheckNil(any interface{}) bool {
    return reflect.ValueOf(any).IsNil()
}

// CheckValid exclude *interface
func CheckValid(any interface{}) bool {
    return reflect.ValueOf(any).IsValid()
}

// DefaultUA Global UA
func DefaultUA() (ua string) {
    ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.87 Safari/537.36"
    return
}

// SetProxy Global Proxy
func SetProxy() (proxy string) {
    proxy = ""
    return
}
