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

	dapperish "github.com/opentracing/api-golang/examples/dapperish"
	"github.com/opentracing/api-golang/opentracing"
)

func client() {
	reader := bufio.NewReader(os.Stdin)
	for {
		span := opentracing.StartTrace("getInput")
		ctx := opentracing.BackgroundGoContextWithSpan(span)
		// Make sure that global trace tag propagation works.
		span.TraceContext().SetTraceAttribute("User", os.Getenv("USER"))
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
		opentracing.AddTraceContextToHeader(
			span.TraceContext(), httpReq.Header, opentracing.DefaultTracer())
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
		reqCtx, err := opentracing.TraceContextFromHeader(
			req.Header, opentracing.DefaultTracer())
		if err != nil {
			panic(err)
		}

		serverSpan := opentracing.JoinTrace(
			"serverSpan", reqCtx,
		).SetTag("component", "server")
		defer serverSpan.Finish()
		fullBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			serverSpan.Error("body read error", err)
		}
		serverSpan.Info("got request with body: " + string(fullBody))
		contextIDMap, tagsMap := opentracing.MarshalTraceContextStringMap(reqCtx)
		fmt.Fprintf(
			w,
			"Hello: %v // %v //  %q",
			contextIDMap,
			tagsMap,
			html.EscapeString(req.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	opentracing.InitDefaultTracer(dapperish.NewTracer("dapperish_tester"))

	go server()
	go client()

	runtime.Goexit()
}
