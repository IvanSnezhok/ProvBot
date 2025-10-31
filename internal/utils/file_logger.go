package utils

import (
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// FileLoggerConfig holds configuration for file logging
type FileLoggerConfig struct {
	LogDir      string // Directory for log files
	LogFileName string // Base name for log file
	MaxSize     int    // Maximum size in megabytes before rotation
	MaxBackups  int    // Maximum number of old log files to retain
	MaxAge      int    // Maximum number of days to retain old log files
	Compress    bool   // Whether to compress rotated log files
}

// InitFileLogger initializes logger with file output and rotation
func InitFileLogger(logLevel string, config *FileLoggerConfig) error {
	if config == nil {
		config = &FileLoggerConfig{
			LogDir:      "log",
			LogFileName: "provbot",
			MaxSize:     100, // 100 MB
			MaxBackups:  4,   // Keep 4 backups (4 weeks)
			MaxAge:      28,  // 28 days
			Compress:    true,
		}
	}

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return err
	}

	// Set up file writer with rotation
	logFile := filepath.Join(config.LogDir, config.LogFileName+".log")
	fileWriter := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
		LocalTime:  true,
	}

	// Create multi-writer (stdout + file)
	multiWriter := NewMultiWriter(os.Stdout, fileWriter)

	// Initialize logger
	Logger = logrus.New()
	Logger.SetOutput(multiWriter)

	// Set formatter to JSON for structured logging
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Set log level
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		Logger.SetLevel(logrus.InfoLevel)
		Logger.Warnf("Invalid log level '%s', defaulting to 'info'", logLevel)
	} else {
		Logger.SetLevel(level)
	}

	Logger.WithFields(logrus.Fields{
		"log_file": logFile,
		"log_dir":  config.LogDir,
	}).Info("File logger initialized")

	return nil
}

// MultiWriter writes to multiple writers
type MultiWriter struct {
	writers []interface {
		Write([]byte) (int, error)
	}
}

// NewMultiWriter creates a new multi-writer
func NewMultiWriter(writers ...interface {
	Write([]byte) (int, error)
}) *MultiWriter {
	return &MultiWriter{writers: writers}
}

// Write writes to all writers
func (mw *MultiWriter) Write(p []byte) (n int, err error) {
	for _, w := range mw.writers {
		n, err = w.Write(p)
		if err != nil {
			return
		}
		if n != len(p) {
			return n, err
		}
	}
	return len(p), nil
}

// ArchiveOldLogs archives logs older than specified days
func ArchiveOldLogs(logDir string, archiveAfterDays int) error {
	files, err := os.ReadDir(logDir)
	if err != nil {
		return err
	}

	cutoffTime := time.Now().AddDate(0, 0, -archiveAfterDays)
	archiveDir := filepath.Join(logDir, "archive")

	// Create archive directory if it doesn't exist
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".log" {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		// Check if file is older than cutoff
		if info.ModTime().Before(cutoffTime) {
			oldPath := filepath.Join(logDir, file.Name())
			archivePath := filepath.Join(archiveDir, file.Name()+".gz")

			// Compress and move to archive
			if err := compressAndMove(oldPath, archivePath); err != nil {
				Logger.WithError(err).Errorf("Failed to archive log file: %s", oldPath)
				continue
			}

			Logger.Infof("Archived log file: %s -> %s", oldPath, archivePath)
		}
	}

	return nil
}

// compressAndMove compresses a file and moves it to archive
func compressAndMove(src, dst string) error {
	// For now, just move the file (compression can be added later with gzip)
	// Lumberjack already handles compression, so we just move rotated files
	if err := os.Rename(src, dst); err != nil {
		return err
	}

	return nil
}
