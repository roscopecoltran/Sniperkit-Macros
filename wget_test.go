package main

import (
    "testing"
)

func TestWget(t *testing.T) {
    err := wget("https://raw.githubusercontent.com/matthieudelaro/nutfile_go1.5/master/nut.yml", "testWget")
    if err != nil {
        t.Error(
            err,
        )
    }
}
