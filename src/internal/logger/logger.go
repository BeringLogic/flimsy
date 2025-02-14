package logger

import (
  "os"
  "log"
)

var logger *log.Logger

func Init() {
  logger = log.New(os.Stderr, "", log.LstdFlags)
  logger.Print(`
  ______  ____    ____  ____    __  ______ __    _ 
 |   ___||    |  |    ||    \  /  ||   ___|\ \  // 
 |   ___||    |_ |    ||     \/   | \-.\-.  \ \//  
 |___|   |______||____||__/\__/|__||______| /__/   
`)
}

func Print(v ...interface{}) {
  logger.Print(v...)
}

func Printf(format string, v ...interface{}) {
  logger.Printf(format, v...)
}

func Fatal(v ...interface{}) {
  logger.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
  logger.Fatalf(format, v...)
}

