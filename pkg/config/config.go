package config

import (
	"bytes"
	"io"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

// ConfigMap defines a configuration map generated from environment parameters read from env file.
type ConfigMap map[string]string

// LoadConfig reads a configuration map from the specified '.env' file.
// Each line in the file should be in the format: path.to.key = "value".
// After loading, environment variables defined in the OS override any
// matching keys in the map.
//
// Returns an error if the file cannot be opened or read.
func LoadConfig(filename string) (configMap ConfigMap, err error) {

	file, error := os.Open(filename)
	if error != nil {
		log.Error().
			Str("Component", "ConfigMap").
			Msgf("An error occurred during opening the required config file [%s]: %s", filename, error)
		return nil, error
	}
	defer file.Close()

	var fileContentBuffer bytes.Buffer
	_, error = io.Copy(&fileContentBuffer, file)
	if error != nil {
		log.Error().
			Str("Component", "ConfigMap").
			Msgf("An error occurred during reading the required config file [%s]: %s", filename, error)
		return nil, error
	}

	envVariables := parseFile(fileContentBuffer.Bytes())
	injectEnvsFromOS(envVariables)

	return envVariables, nil
}

// injectEnvsFromOS updates the provided environment variable map with
// values from the OS environment. Only keys that already exist in the map
// are overridden in final environment variable map.
func injectEnvsFromOS(envVariables map[string]string) {

	for _, envVariable := range os.Environ() {

		// os.Environ() returns key=value, split the value into two tokens
		parts := bytes.SplitN([]byte(envVariable), []byte("="), 2)
		if len(parts) != 2 {
			continue
		}
		key := string(parts[0])
		value := string(parts[1])

		// Override the value in the config map, if found in env map
		_, exists := envVariables[key]
		if exists {
			log.Trace().Msgf("Injecting OS environment variable [%s] with value [%s]", key, value)
			envVariables[key] = value
		}
	}
}

// parseFile parses a byte array representing a configuration file into a map of key-value pairs.
// Lines are only parsed if they:
//   - are not empty
//   - do not start with '#' or '//'
//   - contain an '=' character
func parseFile(fileContent []byte) map[string]string {

	extractedEnvVariables := make(map[string]string)
	cursor := 0
	fileContentLength := len(fileContent)

	for cursor < fileContentLength {

		// Calculate the line size until newline character
		lineSize := bytes.IndexByte(fileContent[cursor:], '\n')
		if lineSize == -1 {
			lineSize = fileContentLength - cursor
		}

		indexToCurrentLineEnd := cursor + lineSize
		fileLineContent := fileContent[cursor:indexToCurrentLineEnd]
		cursor = indexToCurrentLineEnd + 1

		// If the trimmed line is empty or is a comment line, skip the line
		fileLineContent = bytes.TrimSpace(fileLineContent)
		if len(fileLineContent) == 0 || fileLineContent[0] == '#' || string(fileLineContent[0:2]) == "//" {
			continue
		}

		keyValueSeparatorInLineIndex := bytes.IndexByte(fileLineContent, '=')
		if keyValueSeparatorInLineIndex == -1 {
			continue
		}

		envVarKey := bytes.TrimSpace(fileLineContent[:keyValueSeparatorInLineIndex])
		if len(envVarKey) == 0 {
			continue
		}

		envVarValue := bytes.TrimSpace(fileLineContent[keyValueSeparatorInLineIndex+1:])
		if bytes.HasPrefix(envVarValue, []byte{'"'}) && bytes.HasSuffix(envVarValue, []byte{'"'}) {
			envVarValue = envVarValue[1 : len(envVarValue)-1]
		}

		extractedEnvVariables[string(envVarKey)] = string(envVarValue)
	}

	return extractedEnvVariables
}

// ReadString returns the string value associated with the given key
// or the provided fallback value if the key does not exist.
func (configMap ConfigMap) ReadString(key string, fallback string) string {

	valueFromMap := configMap[key]
	if valueFromMap == "" {
		return fallback
	}
	return valueFromMap
}

// ReadBoolean returns the boolean value associated with the given key.
// If the key does not exist, contains an empty string or cannot be
// converted to a boolean, the provided fallback value is returned.
func (configMap ConfigMap) ReadBoolean(key string, fallback bool) bool {

	valueFromMap := configMap[key]
	if valueFromMap == "" {
		return fallback
	}

	convertedValue, err := strconv.ParseBool(valueFromMap)
	if err != nil {
		return fallback
	}
	return convertedValue
}

// ReadInt returns the integer value associated with the given key.
// If the key does not exist, contains an empty string or cannot be
// converted to an integer, the provided fallback value is returned.
func (configMap ConfigMap) ReadInt(key string, fallback int) int {

	valueFromMap := configMap[key]
	if valueFromMap == "" {
		return fallback
	}

	convertedValue, err := strconv.Atoi(valueFromMap)
	if err != nil {
		return fallback
	}
	return convertedValue
}

// ReadFloat retrieves the float64 value associated with the given key.
// If the key does not exist or the value cannot be parsed as a float,
// the provided fallback is returned.
func (configMap ConfigMap) ReadFloat(key string, fallback float64) float64 {

	valueFromMap := configMap[key]
	if valueFromMap == "" {
		return fallback
	}

	convertedValue, err := strconv.ParseFloat(valueFromMap, 64)
	if err != nil {
		return fallback
	}
	return convertedValue
}
