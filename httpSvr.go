package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type HttpSvr struct {
	config  Config
	session string
}

type HttpSvrHandler interface {
	Router() string
	PostHandle(jsonReq string) (jsonRes string)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var randStrBytes = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
var randStrByteLen = len(randStrBytes)

func randStr(n int) string {
	if n <= 0 {
		return ""
	}

	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = randStrBytes[rand.Intn(randStrByteLen)]
	}
	return string(b)
}

func (svr *HttpSvr) StartServe(handles []HttpSvrHandler) {
	router := mux.NewRouter()

	rootPath := "/" + svr.config.MainURLHash + "/"

	// static files directory
	fs := http.FileServer(http.Dir(svr.config.StaticFileDir + "/"))
	router.Handle(rootPath, http.StripPrefix(rootPath, fs))

	// login router
	router.HandleFunc(rootPath+"login/", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			var result string
			var err error
			if body, err := ioutil.ReadAll(req.Body); err != nil {
				jsonData := make(map[string]interface{})
				if err = json.Unmarshal(body, &jsonData); err != nil {
					if secret, exist := jsonData["secret"]; exist && secret == svr.config.SecrtKey {
						result = "OK"
						svr.session = randStr(10)
					} else {
						result = "secret key error!"
					}
				}
			}

			if err != nil {
				result = err.Error()
			}

			res.Write([]byte(fmt.Sprintf("{\"result\":\"%s\"}", result)))
		}
	})

	// handle some functions for POST method
	if handles != nil && len(handles) > 0 {
		for _, v := range handles {
			router.HandleFunc(rootPath+v.Router()+"/", func(res http.ResponseWriter, req *http.Request) {
				if req.Method == "POST" {
					if key, err := req.Cookie("session"); err == nil && key.Value == svr.session {
						if body, err := ioutil.ReadAll(req.Body); err != nil {
							res.Write([]byte(v.PostHandle(string(body))))
						}
					} else {
						// TODO return client to let it relogin
					}
				}
			})
		}
	}

	(&http.Server{
		Addr:         ":" + strconv.Itoa(svr.config.SvrPort),
		Handler:      router,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}).ListenAndServe()
}
