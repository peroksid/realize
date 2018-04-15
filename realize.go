package core

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"fmt"
	"time"
	"go/build"
	"log"
	"bytes"
)

// Custom log
type Log struct{}

// Realize main struct
type Realize struct {
	Options Options `yaml:"options,omitempty" json:"options,omitempty"`
	// Server  Server      `yaml:"server,omitempty" json:"server,omitempty"`
	Schema []Activity  `yaml:"schema,inline,omitempty" json:"schema,inline,omitempty"`
	Sync   chan string `yaml:"-" json:"-"`
	Exit   chan bool   `yaml:"-" json:"-"`
}

// initial set up
func init() {
	// custom log
	log.SetFlags(0)
	log.SetOutput(Log{})
	if build.Default.GOPATH == "" {
		log.Fatal("GOPATH isn't set properly")
	}
	path := filepath.SplitList(build.Default.GOPATH)
	if err := os.Setenv("GOBIN", filepath.Join(path[len(path)-1], "bin")); err != nil {
		log.Fatal("GOBIN impossible to set", err)
	}
}

// Dot check an hidden file
func Dot(path string) bool {
	if runtime.GOOS != "windows" {
		if filepath.Base(path)[0:1] == "." {
			return true
		}
	}
	return false
	// need a way to check on windows
}

// Ext return file extension
func Ext(path string) string {
	var ext string
	for i := len(path) - 1; i >= 0 && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			ext = path[i:]
			if index := strings.LastIndex(ext, "."); index > 0 {
				ext = ext[index:]
			}
		}
	}
	if ext != "" {
		return ext[1:]
	}
	return ""
}

func Print(msg ...interface{}) string{
	var buffer bytes.Buffer
	for i := 0; i < len(msg); i++ {
			buffer.WriteString(msg[i].(string) + " ")
	}
	return buffer.String()
}

// Event print a new message on cli and stream on server
func Record(prefix string, msg interface{}) {
 // switch type err string
	switch m := msg.(type) {
	case string:
		log.Println(prefix, m)
	case error:
		log.Println(prefix, m)
	}
}

// Rewrite timestamp log layout
func (l Log) Write(bytes []byte) (int, error) {
	if len(bytes) > 0 {
		return fmt.Fprint(Output, Yellow.Regular("["), time.Now().Format("15:04:05"), Yellow.Regular("]"), string(bytes))
	}
	return 0, nil
}
