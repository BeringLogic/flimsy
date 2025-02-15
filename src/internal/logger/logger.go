package logger


import (
  "os"
  "log"
)


type FlimsyLogger struct {
  logger *log.Logger
}


func CreateNew() *FlimsyLogger {
  flimsyLogger := &FlimsyLogger{}

  flimsyLogger.logger = log.New(os.Stderr, "", log.LstdFlags)
  flimsyLogger.logger.Print(`
  ______  ____    ____  ____    __  ______ __    _ 
 |   ___||    |  |    ||    \  /  ||   ___|\ \  // 
 |   ___||    |_ |    ||     \/   | \-.\-.  \ \//  
 |___|   |______||____||__/\__/|__||______| /__/   
`)

  return flimsyLogger
}

func (log *FlimsyLogger) Print(v ...interface{}) {
  log.logger.Print(v...)
}

func (log *FlimsyLogger) Printf(format string, v ...interface{}) {
  log.logger.Printf(format, v...)
}

func (log *FlimsyLogger) Fatal(v ...interface{}) {
  log.logger.Fatal(v...)
}

func (log *FlimsyLogger) Fatalf(format string, v ...interface{}) {
  log.logger.Fatalf(format, v...)
}
