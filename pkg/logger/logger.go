package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Logger é uma interface para logging
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

// FileLogger implementa a interface Logger
type FileLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

// New cria um novo logger
func New(logPath string) (*FileLogger, error) {
	// Cria o diretório de logs se não existir
	if err := os.MkdirAll(logPath, 0755); err != nil {
		return nil, fmt.Errorf("erro ao criar diretório de logs: %v", err)
	}

	// Abre o arquivo de log
	logFile := filepath.Join(logPath, "falcon-agent.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo de log: %v", err)
	}

	// Cria um multi-writer para escrever tanto no arquivo quanto no stdout
	multiWriter := io.MultiWriter(os.Stdout, file)

	return &FileLogger{
		infoLogger:  log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLogger: log.New(multiWriter, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}, nil
}

// Info registra mensagens de informação
func (l *FileLogger) Info(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}

// Error registra mensagens de erro
func (l *FileLogger) Error(format string, v ...interface{}) {
	l.errorLogger.Printf(format, v...)
}

// Debug registra mensagens de debug
func (l *FileLogger) Debug(format string, v ...interface{}) {
	l.debugLogger.Printf(format, v...)
}
