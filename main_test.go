package main

import (
	"log"
	"net/http"
	"testing"
)

func TestMain(t *testing.T) {

	t.Run("", func(t *testing.T) {

		go main()

		resp, err := http.Get("http://localhost:8080/")
		if err != nil {
			t.Error(err)
		} else {
			log.Println("StatusCode == 200 is ", resp.StatusCode == 200)
		}
	})

}
