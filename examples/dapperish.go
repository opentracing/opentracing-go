package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/examples/dapperish"
)

func client() {
	reader := bufio.NewReader(os.Stdin)
	for {
		ctx, span := opentracing.BackgroundContextWithSpan(
			opentracing.StartTrace("getInput"))
		// Make sure that global trace tag propagation works.
		span.SetTraceAttribute("User", os.Getenv("USER"))
		span.LogEventWithPayload("ctx", ctx)
		fmt.Print("\n\nEnter text (empty string to exit): ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if len(text) == 0 {
			fmt.Println("Exiting.")
			os.Exit(0)
		}

		span.LogEvent(text)

		httpClient := &http.Client{}
		httpReq, _ := http.NewRequest("POST", "http://localhost:8080/", bytes.NewReader([]byte(text)))
		opentracing.PropagateSpanInHeader(
			span, httpReq.Header, opentracing.GlobalTracer())
		resp, err := httpClient.Do(httpReq)
		if err != nil {
			span.LogEventWithPayload("error", err)
		} else {
			span.LogEventWithPayload("got response", resp)
		}

		span.Finish()
	}
}

func server() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		serverSpan, err := opentracing.JoinTraceFromHeader(
			"serverSpan", req.Header, opentracing.GlobalTracer())
		if err != nil {
			panic(err)
		}
		serverSpan.SetTag("component", "server")
		defer serverSpan.Finish()

		fullBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			serverSpan.LogEventWithPayload("body read error", err)
		}
		serverSpan.LogEventWithPayload("got request with body", string(fullBody))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	opentracing.InitGlobalTracer(dapperish.NewTracer("dapperish_tester"))

	go server()
	go client()

	runtime.Goexit()
}
