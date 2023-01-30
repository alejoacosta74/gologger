/*
Package gologger defines an interface and
logger based on [logrus](https://github.com/sirupsen/logrus) library,
implementing out-of-the-box features such as:
- singleton logic
- debuging option (`WithDebugLevel(true)`)
- decorated output (`WithRunTimeContext()`)
- nologger option (`WithNullLogger()`)
- option to write loggin output to file
- convenience methods to create new logger instances (i.e. `(l) func NewLoggerWithField()`)
- option to set custom log level
*/
package gologger
