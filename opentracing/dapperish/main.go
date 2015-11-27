package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/opentracing/api-golang/opentracing"
)

func client() {
	reader := bufio.NewReader(os.Stdin)
	for {
		span, ctx := opentracing.StartSpan("getInput")
		span.Info("ctx: ", ctx)
		fmt.Print("\n\nEnter text (empty string to exit): ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if len(text) == 0 {
			fmt.Println("Exiting.")
			os.Exit(0)
		}

		span.Info(text)

		httpClient := &http.Client{}
		httpReq, _ := http.NewRequest("POST", "http://localhost:8080/", bytes.NewReader([]byte(text)))
		opentracing.AddContextIDToHttpHeader(span.ContextID(), httpReq.Header)
		resp, err := httpClient.Do(httpReq)
		if err != nil {
			span.Error("error: ", err)
		} else {
			span.Info("got response: ", resp)
		}

		span.Finish()
	}
}

func server() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		reqCtxID, err := opentracing.GetContextIDFromHttpHeader(
			req.Header, opentracing.GlobalTracer())
		if err != nil {
			panic(err)
		}

		serverSpan, _ := opentracing.StartSpan("serverSpan", reqCtxID)
		defer serverSpan.Finish()
		fullBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			serverSpan.Error("body read error", err)
		}
		serverSpan.Info("got request with body: " + string(fullBody))
		fmt.Fprintf(w, "Hello: %v / %q", reqCtxID.Serialize(), html.EscapeString(req.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	opentracing.InitGlobalTracer(
		NewTrivialRecorder("dapperish_tester"),
		NewDapperishContextIDSource())

	go client()
	go server()

	runtime.Goexit()
}
