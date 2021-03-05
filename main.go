package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/http"
	"os"
	"strings"
)

var redisAddr = fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     redisAddr,
	Username: os.Getenv("REDIS_USERNAME"),
	Password: os.Getenv("REDIS_PASSWORD"),
	DB:       0,
	TLSConfig: &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         os.Getenv("REDIS_HOST"),
		ClientAuth:         tls.RequireAnyClientCert,
	},
})

func main() {
	defer rdb.Close()
	http.HandleFunc("/", GetAccountsServer)
	http.ListenAndServe(os.Getenv("GO_PORT"), nil)
}

func GetAccountsServer(w http.ResponseWriter, r *http.Request) {
	var accountName = strings.ToLower(r.URL.Path[1:])

	accountID, err := rdb.Get(ctx, accountName).Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}

	if err == redis.Nil {
		fmt.Fprintf(w, "{\"%s\":null}", accountName)
	} else {
		fmt.Fprintf(w, "{\"%s\":%s}", accountName, accountID)
	}
}
