package dmutils

import (
	crand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

func SlugifyBase(in string) string {
	ot := strings.ToLower(in)
	rep := strings.NewReplacer(" ", "-", "+", "-", "_", "-", ".", "-", "/", "-")
	return rep.Replace(ot)
}

func SetIds(separator string, opts ...string) string {
	ids := ""
	for _, val := range opts {
		if ids == "" {
			ids = val
		} else {
			ids = ids + separator + val
		}
	}
	return ids
}

func GenerateRandomString(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

func GenerateRandomANString(length int) (string, error) {
	// Precompute the charset length as it won't change
	charsetLength := big.NewInt(int64(len(letterBytes)))

	// Create a byte slice to store the result
	result := make([]byte, length)

	// Generate random characters
	for i := range result {
		// Get a random number from the charset
		num, err := crand.Int(crand.Reader, charsetLength)
		if err != nil {
			return "", err
		}

		// Use the random number to index into the charset
		result[i] = letterBytes[num.Int64()]
	}

	return string(result), nil
}

func GenerateUniqueID() (string, error) {
	timestamp := time.Now().UnixNano()

	randomPart, err := GenerateRandomANString(6)
	if err != nil {
		return "", err
	}

	uniquePart := fmt.Sprintf("%x", timestamp)[:3] // first 3 characters from timestamp in hex

	// Return final unique ID (3 chars from time and 6 random)
	return uniquePart + randomPart, nil
}

func IsBase64Encoded(in string) bool {
	// Check if the length is valid (divisible by 4)
	if len(in) == 0 || len(in)%4 != 0 {
		return false
	}

	// Regular expression to validate Base64 content
	validChars := "^[A-Za-z0-9+/]*={0,2}$"
	matched, err := regexp.MatchString(validChars, in)
	if err != nil || !matched {
		return false
	}

	// Additional check: padding should only be at the end
	// and follow correct patterns
	paddingIndex := strings.IndexByte(in, '=')
	if paddingIndex >= 0 {
		for i := paddingIndex; i < len(in); i++ {
			if in[i] != '=' {
				return false
			}
		}
	}

	return true
}

func Base64Detectors(in string) (string, bool) {
	isBs := false
	if IsBase64Encoded(in) {
		isBs = true
	} else {
		return in, isBs
	}
	bytes, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return in, isBs
	}
	return string(bytes), isBs
}

func BoolToString(wp bool) string {
	if wp {
		return "true"
	}
	return "false"
}

func StringToBool(s string) bool {
	re, err := strconv.ParseBool(s)
	if err != nil {
		if slices.Contains([]string{"true", "True", "TRUE", "t", "T", "yes", "Yes", "YES", "y", "Y", "1"}, s) {
			return true
		} else {
			return false
		}
	}
	return re
}

func ConvertStringToType(value string) any {
	if value == "" {
		return nil
	}

	lowerVal := strings.ToLower(value)
	if lowerVal == "true" || lowerVal == "false" {
		boolValue, err := strconv.ParseBool(lowerVal)
		if err == nil {
			return boolValue
		}
	}

	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}

	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		return floatValue
	}

	var arrayResult []interface{}
	if err := json.Unmarshal([]byte(value), &arrayResult); err == nil {
		return arrayResult
	}

	var mapResult map[string]any
	if err := json.Unmarshal([]byte(value), &mapResult); err == nil {
		return mapResult
	}

	return value
}

func StringSlicer(s string, limit int) string {
	if len(s) > limit {
		return s[:limit]
	}
	return s
}

func ExtractBetweenDelimiters(input, startDelimiter, endDelimiter string) (string, error) {
	startIndex := strings.Index(input, startDelimiter)
	if startIndex == -1 {
		return input, nil
	}
	log.Println("Start index", startIndex)
	startIndex += len(startDelimiter)
	endIndex := strings.Index(input[startIndex:], endDelimiter)
	if endIndex == -1 {
		return strings.ReplaceAll(input, startDelimiter, ""), nil
	}
	log.Println("End index", endIndex)
	endIndex += startIndex
	return input[startIndex:endIndex], nil
}

func TextToMD[V any](heading V, items []V) string {
	var builder strings.Builder

	// Add heading
	builder.WriteString(fmt.Sprintf("## %v\n\n", heading))

	// Add list items
	for _, item := range items {
		builder.WriteString(fmt.Sprintf("- %v\n", item))
	}

	return builder.String()
}

func AnyToStr(v any) string {
	switch v := v.(type) {
	case []byte:
		return string(v)
	case nil:
		return ""
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.Itoa(int(v))
	case int16:
		return strconv.Itoa(int(v))
	case int32:
		return strconv.Itoa(int(v))
	case int64:
		return strconv.Itoa(int(v))
	case bool:
		return strconv.FormatBool(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case time.Time:
		return v.Format(time.RFC3339)
	case map[string]any:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(jsonBytes)
	case []any:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(jsonBytes)
	case struct{}:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(jsonBytes)
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}
