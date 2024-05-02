package logger

import (
	"io"
	"log"
	"os"
)

var l *log.Logger

func InitLogger() {
	l = log.Default()
}

func GetLogger() *log.Logger {
	return l
}

func NewLogger(path string, prefix string) *log.Logger {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Panic(err)
	}
	newLogger := log.New(io.MultiWriter(os.Stdout, logFile), prefix, log.Lshortfile|log.Ldate|log.Ltime)
	return newLogger
}

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of
func Print(v ...any) {
	l.Print(v...)
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of [fmt.Printf].
func Printf(format string, v ...any) {
	l.Printf(format, v...)
}

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of [fmt.Println].
func Println(v ...any) {
	l.Println(v...)
}

// Fatal is equivalent to l.Print() followed by a call to [os.Exit](1).
func Fatal(v ...any) {
	l.Fatal(v...)
}

// Fatalf is equivalent to l.Printf() followed by a call to [os.Exit](1).
func Fatalf(format string, v ...any) {
	l.Fatalf(format, v...)
}

// Fatalln is equivalent to l.Println() followed by a call to [os.Exit](1).
func Fatalln(v ...any) {
	l.Fatalln(v...)
}

// Panic is equivalent to l.Print() followed by a call to panic().
func Panic(v ...any) {
	l.Panic(v...)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func Panicf(format string, v ...any) {
	l.Panicf(format, v...)
}

// Panicln is equivalent to l.Println() followed by a call to panic().
func Panicln(v ...any) {
	l.Panicln(v...)
}
