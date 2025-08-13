package config

// Default configuration values
var defaults = map[string]string{
	"PORT":            "8080",
	"DATA_FILE":       "data.tsv",
	"LOG_LEVEL":       "info",
	"LOG_FORMAT":      "text",
	"LOG_OUTPUT":      "stdout",
	"MAX_FILE_SIZE":   "1073741824", // 1GB
	"BACKUP_DIR":      "backups",
	"COMPRESSION":     "false",
	"MAX_CONNECTIONS": "100",
	"CONNECTION_TTL":  "300", // 5 minutes
	"QUERY_TIMEOUT":   "30",  // 30 seconds
	"READ_TIMEOUT":    "30",  // 30 seconds
	"WRITE_TIMEOUT":   "30",  // 30 seconds
	"IDLE_TIMEOUT":    "120", // 2 minutes
}
