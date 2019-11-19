package utils

import (
    "log"
    "net"
    "net/http"
    "time"

    "golang.org/x/net/proxy"
)

func init() {
    log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func s5Proxy(proxyURL string) (transport *http.Transport) {
    dialer, err := proxy.SOCKS5("tcp", proxyURL,
        nil,
        &net.Dialer{
            Timeout:   30 * time.Second,
            KeepAlive: 30 * time.Second,
        },
    )
    if err != nil {
        log.Fatal(" [S5 Proxy Error]: ", err)
    }
    transport = &http.Transport{
        Proxy:               nil,
        Dial:                dialer.Dial,
        TLSHandshakeTimeout: 10 * time.Second,
    }
    return
}

// HTTPClient http client
func HTTPClient(proxy string) (client http.Client) {
    client = http.Client{Timeout: 30 * time.Second}
    if proxy != "" {
        transport := s5Proxy(proxy)
        client = http.Client{Timeout: 30 * time.Second, Transport: transport}
    }
    return
}
