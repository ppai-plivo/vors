package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	redis "github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

const (
	redisEndpoint = "127.0.0.1:6379"
	streamName    = "video:redis"
	streamField   = "blob"
)

var (
	client redis.UniversalClient
)

func main() {

	client = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{redisEndpoint},
	})
	if _, err := client.Ping().Result(); err != nil {
		log.Fatalf("Failed to connect to redis at %s: %s", redisEndpoint, err.Error())
	}

	router := mux.NewRouter()
	router.PathPrefix("/videos").Handler(http.HandlerFunc(streamMedia)).Methods("GET")

	chain := alice.New(
		NewRequestLogger(log.Writer()),
	).Then(router)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: chain,
	}

	log.Fatal(srv.ListenAndServe())
}

func streamMedia(w http.ResponseWriter, r *http.Request) {

	if strings.HasSuffix(r.URL.Path, "m3u8") {
		key := strings.ReplaceAll(r.URL.Path[1:], "/", ":")
		val, err := client.Get(key).Result()
		if err != nil {
			log.Printf("redis GET failed: %s", err.Error())
			http.Error(w, "reading playlist failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(StringToBytes(val))
		return
	}

	l := strings.Split(r.URL.Path[1:], "/")
	id := l[len(l)-1]
	_, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "part name not an integer", http.StatusInternalServerError)
		return
	}

	key := strings.Join(l[:len(l)-1], ":")
	msgs, err := client.XRange(key, id, id).Result()
	if err != nil {
		log.Printf("client.XRange failed: %s", err.Error())
		http.Error(w, "getting part failed", http.StatusInternalServerError)
		return
	}

	i, ok := msgs[0].Values[streamField]
	if !ok {
		http.Error(w, "getting part failed 2", http.StatusInternalServerError)
		return
	}

	s, ok := i.(string)
	if !ok {
		http.Error(w, "getting part failed 3", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(StringToBytes(s))
}
