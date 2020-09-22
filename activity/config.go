package activity

const (
	noneType = iota
	syslogType
	fileType
)

// Config - activity configuration
type Config struct {
	// Type: 0 - none, 1 - syslog, 2 - file
	Type int
	// Name: if type is syslog then name is priority, else name is filename
	Name string
	// Tag: syslog tag name
	Tag string
}
