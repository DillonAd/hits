package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/honeycombio/beeline-go"
	"github.com/honeycombio/beeline-go/wrappers/hnygorilla"
	"github.com/honeycombio/beeline-go/wrappers/hnynethttp"
)

var connStr string
var salt string

func main() {
	connStr = os.Getenv("CONN_STR")
	salt = os.Getenv("SALT")
	writeKey := os.Getenv("HONEYCOMB_WRITEKEY")

	beeline.Init(beeline.Config{
		WriteKey: writeKey,
		Dataset:  "hits",
	})

	defer beeline.Close()

	r := mux.NewRouter()
	r.Use(hnygorilla.Middleware)

	r.HandleFunc("/", homeHandler).
		Methods("GET").
		Schemes("http")

	r.HandleFunc("/hit/{tenantId}/{pageName}", hitHandler).
		Methods("POST").
		Schemes("http")

	listenerPort := ":8080"

	log.Printf("Listening on port %s\n", listenerPort)
	log.Fatalln(http.ListenAndServe(listenerPort, hnynethttp.WrapHandler(r)))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	_, span := beeline.StartSpan(r.Context(), "Get the teapot!")
	defer span.Send()
	w.WriteHeader(http.StatusTeapot)
}

func hitHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenantId"]
	pageName := vars["paCountHiteName"]

	beeline.AddField(r.Context(), "tenantID", tenantID)
	beeline.AddField(r.Context(), "pageName", pageName)

	storage, err := NewStorage(connStr)

	if err != nil {
		log.Fatalln(err)
		beeline.AddField(r.Context(), "NewStorage Error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer storage.Disconnect()

	hitService := NewHitService(storage, salt)

	_, span := beeline.StartSpan(r.Context(), "CountHit")
	defer span.Send()
	err = hitService.CountHit(tenantID, pageName, r.RemoteAddr)

	if err != nil {
		log.Fatalln(err)
		beeline.AddField(r.Context(), "CountHit Error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
