package cache

type FileFormat int

const (
	Json FileFormat = iota
	Yaml
)

const (
	jsonOutput    = "json"
	yamlOutput    = "yaml"
	unknownOutput = "unknown"
)

func (s FileFormat) String() string {
	switch s {
	case Json:
		return jsonOutput
	case Yaml:
		return yamlOutput
	default:
		return unknownOutput
	}
}
