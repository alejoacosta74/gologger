# gologger

A golang logger based on [logrus](https://github.com/sirupsen/logrus) library, implementing out-of-the-box features such as:
- singleton logic
- debuging option (`WithDebugLevel(true)`)
- decorated output (`WithRunTimeContext()`)
- nologger option (`WithNullLogger()`)
- option to write loggin output to file instead of stdout
- convenience methods to create new logger instances (i.e. `(l) func NewLoggerWithField()`)

## usage

### Example: Debug level with decorated output style

```golang
package main

import (
	log "github.com/alejoacosta74/gologger"
)

func main() {
	logger, err := log.NewLogger(log.WithDebugLevel(true), log.WithField("app", "example"), log.WithRuntimeContext())
	if err != nil {
		panic(err)
	}

	logger.Debug("Hello from logger...!")

}
```

Output:
```bash
❯ go run example.go
DEBU[19-11-2022 13:53:14] example.go:13 - main.main -  Hello from logger...!                         app=example fields.file=/Users/alejoacosta/code/golang/mylogger/logger.go fields.func=github.com/alejoacosta74/gologger.createNewLogger file=/Users/alejoacosta/code/golang/mylogger/logger.go func=github.com/alejoacosta74/gologger.createNewLogger line=70
```

### Example: using singleton feature

```golang
import (
	"sync"

	log "github.com/alejoacosta74/gologger"
)

func main() {
	_, err := log.NewLogger(log.WithField("loggerid", 0))
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(3)
	for i := 1; i < 4; i++ {
		go func(i int) {
			newLogger, err := log.NewLogger(log.WithField("loggerid", i))
			if err != nil {
				panic(err)
			}

			newLogger.Info("Hello from goroutine id: ", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
```

output:
```bash
❯ go run example.go
INFO Hello from goroutine id: 3                    loggerid=0
INFO Hello from goroutine id: 1                    loggerid=0
INFO Hello from goroutine id: 2                    loggerid=0
```