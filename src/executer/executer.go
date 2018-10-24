package executer

import (
	"fmt"
	"remote"
	"term"
)

var (
	pool          *remote.Pool
	currentUser   string
	currentRaise  remote.RaiseType
	currentPasswd string
)

// ExecResult represents result of execution of a task
type ExecResult struct {
	// Codes is a map host -> statuscode
	Codes map[string]int
	// Success holds successful hosts
	Success []string
	// Error holds unsuccessful hosts
	Error []string
	// Stopped holds hosts which weren't able to complete task
	Stopped int
	// OutputMap structures hosts by different outputs
	OutputMap map[string][]string
}

// Initialize initializes executer pool and configuration
func Initialize(numThreads int, user string) {
	pool = remote.NewPool(numThreads)
	currentUser = user
	currentRaise = remote.RaiseTypeNone
	currentPasswd = ""
}

// SetUser sets current user
func SetUser(user string) {
	currentUser = user
}

// SetRaise sets current privileges raise type
func SetRaise(raise remote.RaiseType) {
	currentRaise = raise
}

// SetPasswd sets current password
func SetPasswd(passwd string) {
	currentPasswd = passwd
}

func newExecResults() *ExecResult {
	er := new(ExecResult)
	er.Codes = make(map[string]int)
	er.Success = make([]string, 0)
	er.Error = make([]string, 0)
	er.OutputMap = make(map[string][]string)
	return er
}

func printExecResults(r *ExecResult) {
	msg := fmt.Sprintf(" Hosts processed: %d, success: %d, error: %d    ",
		len(r.Success)+len(r.Error), len(r.Success), len(r.Error))
	h := term.HR(len(msg))
	fmt.Println(term.Green(h))
	fmt.Println(term.Green(msg))
	fmt.Println(term.Green(h))
}
