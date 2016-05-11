package utils

import (
    "io"
    "net/http"
    "net/url"
    "os"
)

// inspired from https://www.socketloop.com/tutorials/golang-download-file-example
// TODO: check with proxies
func Wget(rawURL string, fileName string) error {
    _, err := url.Parse(rawURL)

    if err != nil {
        return err
    }
    file, err := os.Create(fileName)

    if err != nil {
        return err
    }
    defer file.Close()

    check := http.Client{
        CheckRedirect: func(r *http.Request, via []*http.Request) error {
            r.URL.Opaque = r.URL.Path
            return nil
            },
        }

    resp, err := check.Get(rawURL) // add a filter to check redirect

    if err != nil {
        return err
    }
    defer resp.Body.Close()
    _, err = io.Copy(file, resp.Body)

    if err != nil {
        return err
    }

    return nil
}
