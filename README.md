Flowbot - A Flowdock robot written in Go
========================================

[![GoDoc](https://godoc.org/github.com/eko/flowbot?status.png)](https://godoc.org/github.com/eko/flowbot)
[![Build Status](https://travis-ci.org/eko/flowbot.png?branch=master)](https://travis-ci.org/eko/flowbot)

This is a library to create a Flowdock robot which runs in Golang.

Overview
--------

It uses the following Flowdock APIs:

* Stream API (to read content on a given flow),
* Push API (to push chat message on the flow or a team inbox notification).

Installation
------------

```bash
$ go get -u github.com/eko/flowbot
```

Run the robot
-------------

```bash
$ go run app.go
Connecting to Flowdock stream https://stream.flowdock.com/flows/[...]/[...]
Connected! Start reading stream...
-> Command found: flow uptime
-> Chat message correctly sent.
```

A robot example application
---------------------------

This sample application answers to the following commands:

* flow uptime: Renders the server uptime,
* flow image <something>: Render the first image found on Google Image corresponding to <something>

```go
package main

import (
    "github.com/eko/flowbot"
    "fmt"
    "io/ioutil"
    "net/http"
    "os/exec"
    "regexp"
)

func main() {
    flowbot.FlowdockFlowToken = "<YOUR FLOW TOKEN API KEY>"
    flowbot.FlowdockStreamUrl = "https://stream.flowdock.com/flows/<YOUR-ORGANIZATION>/<YOUR-FLOW-NAME>"
    flowbot.FlowdockAuthUsername = "<YOUR EMAIL ADDRESS>"
    flowbot.FlowdockAuthPassword = "<YOUR PASSWORD>"

    flowbot.AddCommand("flow uptime", func (command flowbot.Command, entry flowbot.Entry) {
        output, err := exec.Command("uptime").Output()
        if err != nil { panic(err) }

        flowbot.SendChat(fmt.Sprintf("My uptime: %s", output))
    })

    flowbot.AddCommand("flow image (.*)", func (command flowbot.Command, entry flowbot.Entry) {
        query := command.Pattern.FindStringSubmatch(entry.Content)[1]

        response, err := http.Get(fmt.Sprintf("https://ajax.googleapis.com/ajax/services/search/images?v=1.0&q=%s", query))
        if err != nil { panic(err) }

        body, _ := ioutil.ReadAll(response.Body)

        imagePattern := regexp.MustCompile(`(http[s]?://[a-zA-Z0-9\/\-\.\_\%]+.[gif|png|jpg|jpeg])`)
        image := imagePattern.FindStringSubmatch(string(body))[0]

        flowbot.SendChat(image)
    })

    flowbot.Stream()
}
```
