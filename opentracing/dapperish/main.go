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
		span := opentracing.Global().StartTrace("getInput")
		ctx := opentracing.BackgroundContextWithSpan(span)
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
		opentracing.AddTraceContextToHttpHeader(span.TraceContext(), httpReq.Header)
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
		reqCtx, err := opentracing.GetTraceContextFromHttpHeader(
			req.Header, opentracing.Global())
		if err != nil {
			panic(err)
		}

		serverSpan := opentracing.Global().JoinTrace(
			"serverSpan", reqCtx,
			"component", "server",
		)
		defer serverSpan.Finish()
		fullBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			serverSpan.Error("body read error", err)
		}
		serverSpan.Info("got request with body: " + string(fullBody))
		fmt.Fprintf(w, "Hello: %v / %q", reqCtx.SerializeString(), html.EscapeString(req.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	opentracing.InitGlobalTracer(
		NewTrivialRecorder("dapperish_tester"),
		NewDapperishTraceContextIDSource())

	go server()
	go client()

	runtime.Goexit()
}
