package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
)

type Logger struct {
	mu       sync.Mutex
	file     *os.File
	logger   *log.Logger
	level    Level
	maxSize  int64
	filename string
}

var (
	instance *Logger
	once     sync.Once
)

func GetInstance() *Logger {
	once.Do(func() {
		var err error
		instance, err = NewLogger("logs/falcon-agent.log", INFO, 10*1024*1024) // 10MB
		if err != nil {
			log.Fatal(err)
		}
	})
	return instance
}

func NewLogger(filename string, level Level, maxSize int64) (*Logger, error) {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{
		file:     file,
		logger:   log.New(file, "", log.Ldate|log.Ltime|log.Lmicroseconds),
		level:    level,
		maxSize:  maxSize,
		filename: filename,
	}, nil
}

func (l *Logger) checkRotate() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	info, err := l.file.Stat()
	if err != nil {
		return err
	}

	if info.Size() < l.maxSize {
		return nil
	}

	// Fecha o arquivo atual
	l.file.Close()

	// Renomeia o arquivo atual com timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	newName := fmt.Sprintf("%s.%s", l.filename, timestamp)
	if err := os.Rename(l.filename, newName); err != nil {
		return err
	}

	// Cria um novo arquivo de log
	file, err := os.OpenFile(l.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	l.file = file
	l.logger = log.New(file, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	return nil
}

func (l *Logger) log(level Level, format string, v ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if err := l.checkRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao rotacionar log: %v\n", err)
	}

	levelStr := ""
	switch level {
	case DEBUG:
		levelStr = "DEBUG"
	case INFO:
		levelStr = "INFO"
	case WARNING:
		levelStr = "WARN"
	case ERROR:
		levelStr = "ERROR"
	}

	msg := fmt.Sprintf(format, v...)
	l.logger.Printf("[%s] %s", levelStr, msg)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.log(DEBUG, format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.log(INFO, format, v...)
}

func (l *Logger) Warning(format string, v ...interface{}) {
	l.log(WARNING, format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.log(ERROR, format, v...)
}

func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.SetOutput(w)
}

func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.file.Close()
}
