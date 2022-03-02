package main

import (
	"context"
	"fmt"
	"go-learn-kvs/kvstorage"
	"log"
	"net/http"
)

var httpErrorStatusCodeMessages = map[int]string{
	http.StatusNotFound:            "404 There is no record in the storage for key '%v'.\n",
	http.StatusInternalServerError: "500 Internal storage error.\n",
	http.StatusMethodNotAllowed:    "Sorry, only GET, POST and DELETE methods are allowed",
}

func main() {
	storage := kvstorage.NewStorage()
	ctx := context.Background()

	if storage == nil {
		log.Fatal("Cannot initialize storage!")
		return
	}

	http.HandleFunc("/", process(ctx, storage))
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func process(ctx context.Context, storage kvstorage.KVStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path[1:]

		switch r.Method {
		case "GET":
			val, err := storage.Get(ctx, key)

			if err != nil {
				http.Error(w, httpErrorStatusCodeMessages[http.StatusInternalServerError], http.StatusInternalServerError)
				return
			}
			if val == nil {
				http.Error(w, fmt.Sprintf(httpErrorStatusCodeMessages[http.StatusNotFound], key), http.StatusNotFound)
				return
			}
			fmt.Fprintf(w, "%s is a %s\n", key, val)

		case "POST":
			val := r.FormValue("value")
			err := storage.Put(ctx, key, val)

			if err != nil {
				http.Error(w, httpErrorStatusCodeMessages[http.StatusInternalServerError], http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "%s is set to %s\n", key, val)

		case "DELETE":
			err := storage.Delete(ctx, key)

			if err != nil {
				http.Error(w, httpErrorStatusCodeMessages[http.StatusInternalServerError], http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "%s is deleted\n", key)

		default:
			http.Error(w, httpErrorStatusCodeMessages[http.StatusMethodNotAllowed], http.StatusMethodNotAllowed)
		}
	}
}
