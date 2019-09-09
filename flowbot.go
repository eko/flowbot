// Flowbot - A Flowdock robot library written in Go
//
// Author: Vincent Composieux <vincent.composieux@gmail.com>

package flowbot

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "regexp"
    "time"
)

// Flowdock configuration values
var (
    FlowdockRobotName    string = "Flowbot"

    FlowdockChatUrl      string = "https://api.flowdock.com/v1/messages/chat/"
    FlowdockInboxUrl     string = "https://api.flowdock.com/v1/messages/team_inbox/"

    FlowdockStreamUrl    string
    FlowdockFlowToken    string
    FlowdockAuthUsername string
    FlowdockAuthPassword string

    commands             []Command
)

// A handler function type
type handler func(pattern Command, entry Entry)

// A command structure
// Streaming API will detect this command name on Flow and run the command handler
type Command struct {
    Pattern *regexp.Regexp
    Handler handler
}

// A Flowdock entry structure
// Each streamed response will be transformed into an event
type Entry struct {
    Event       string    `json:"event"`
    Tags        []string  `json:"tags"`
    Uuid        string    `json:"uuid"`
    Persist     bool      `json:"persist"`
    Id          int       `json:"id"`
    Flow        string    `json:"flow"`
    Content     string    `json:"content"`
    Sent        int       `json:"sent"`
    App         string    `json:"app"`
    CreatedAt   time.Time `json:"created_at"`
    Attachments []string  `json:"attachments"`
    User        string    `json:"user"`
    ThreadId    string    `json:"thread_id,omitempty"`
}

// A Flowdock chat message
type Chat struct {
    Content          string `json:"content"`
    ExternalUserName string `json:"external_user_name"`
    ThreadId         string `json:"thread_id,omitempty"`
}

// Throws a panic error in case of err is defined
func check_error(err error) {
    if err != nil { panic(err) }
}

// Adds a command
func AddCommand(pattern string, handler handler) {
    commands = append(commands, Command{Pattern: regexp.MustCompile(pattern), Handler: handler})
}

// Start streaming on Flowdock
func Stream() {
    request, err := http.NewRequest("GET", FlowdockStreamUrl, nil)
    request.SetBasicAuth(FlowdockAuthUsername, FlowdockAuthPassword)
    check_error(err)

    fmt.Printf("Connecting to Flowdock stream %s\n", FlowdockStreamUrl)

    client        := http.Client{}
    response, err := client.Do(request)
    check_error(err)

    fmt.Print("Connected! Start reading stream...\n")

    reader := bufio.NewReader(response.Body)

    for {
        line, err := reader.ReadBytes('\n')
        check_error(err)

        var entry Entry
        err = json.Unmarshal(line, &entry)

        for _, command := range commands  {
            if m, _ := regexp.MatchString(command.Pattern.String(), entry.Content); m {
                fmt.Printf("-> Command: %s\n", entry.Content)
                command.Handler(command, entry)
                break
            }
        }
    }
}

// Send a chat message
func SendChat(message string) {
    parameters := Chat{
        Content:          message,
        ExternalUserName: FlowdockRobotName,
    }

    sendChat(parameters)
}

// Send a chat message to a thread
func SendThreadChat(threadId string, message string) {
    parameters := Chat{
        Content:          message,
        ExternalUserName: FlowdockRobotName,
        ThreadId:         threadId,
    }

    sendChat(parameters)
}

// Send a chat message to configured flowdock flow token
func sendChat(parameters Chat) {
    jsonParameters, err := json.Marshal(parameters)
    check_error(err)

    response, err := http.Post(FlowdockChatUrl + FlowdockFlowToken, "application/json", bytes.NewReader(jsonParameters))

    if response.StatusCode == 200 {
        fmt.Print("-> Chat message correctly sent.\n")
    } else {
        body, _ := ioutil.ReadAll(response.Body)
        fmt.Errorf("-> Error while sending chat message: %s %s\n", response.Status, string(body))
    }
}

// Send a team inbox message to configured flowdock flow token
func SendInbox(source string, from_address string, subject string, content string) {
    parameters := map[string]string{
        "source":       source,
        "from_address": from_address,
        "subject":      subject,
        "content":      content,
    }

    jsonParameters, err := json.Marshal(parameters)
    check_error(err)

    response, err := http.Post(FlowdockInboxUrl + FlowdockFlowToken, "application/json", bytes.NewReader(jsonParameters))

    if response.StatusCode == 200 {
        fmt.Print("-> Inbox team message correctly sent.\n")
    } else {
        body, _ := ioutil.ReadAll(response.Body)
        fmt.Errorf("-> Error while sending inbox team message: %s %s\n", response.Status, string(body))
    }
}
