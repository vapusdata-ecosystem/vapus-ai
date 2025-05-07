package dmutils

import (
	"fmt"
	"strings"
)

type MarkdownConverter struct {
	bulletChar     string
	numberFormat   string
	indentString   string
	escapeReplacer *strings.Replacer
}

type MarkdownOpts struct {
	KeysToEscape    map[string]bool
	KeysToCodeBlock map[string]bool
	KeysToIgnore    map[string]bool
	EscapeAllKeys   bool
}

func NewMarkdownConverter() *MarkdownConverter {
	conv := &MarkdownConverter{
		bulletChar:   "-",
		numberFormat: "%d.",
		indentString: "  ",
		escapeReplacer: strings.NewReplacer(
			"*", "\\*",
			"_", "\\_",
			"`", "\\`",
			"#", "\\#",
			"[", "\\[",
			"]", "\\]",
		),
	}
	return conv
}

// Option functions
func WithEscapeAllKeys() func(*MarkdownOpts) {
	return func(m *MarkdownOpts) { m.EscapeAllKeys = true }
}

func WithKeysToEscape(keys ...string) func(*MarkdownOpts) {
	return func(m *MarkdownOpts) {
		for _, key := range keys {
			m.KeysToEscape[key] = true
		}
	}
}

func WithKeysToCodeBlock(keys ...string) func(*MarkdownOpts) {
	return func(m *MarkdownOpts) {
		for _, key := range keys {
			m.KeysToCodeBlock[key] = true
		}
	}
}

func WithKeysToIgnore(keys ...string) func(*MarkdownOpts) {
	return func(m *MarkdownOpts) {
		for _, key := range keys {
			m.KeysToIgnore[key] = true
		}
	}
}

func (m *MarkdownConverter) Convert(data interface{}, options ...func(*MarkdownOpts)) string {
	opts := &MarkdownOpts{}
	for _, opt := range options {
		opt(opts)
	}
	return m.convertValue(data, 0, "", opts)
}

func (m *MarkdownConverter) convertValue(data interface{}, indentLevel int, parentKey string, opts *MarkdownOpts) string {
	var builder strings.Builder
	indent := strings.Repeat(m.indentString, indentLevel)

	switch v := data.(type) {
	case map[string]interface{}:
		for key, val := range v {
			if opts.KeysToIgnore[key] {
				continue
			}

			fullKeyPath := parentKey + "." + key
			if parentKey == "" {
				fullKeyPath = key
			}

			builder.WriteString(fmt.Sprintf("%s%s ", indent, m.bulletChar))

			// Apply escaping rules
			displayKey := key
			if m.shouldEscapeKey(fullKeyPath, opts) {
				displayKey = m.escapeReplacer.Replace(key)
			}
			builder.WriteString(fmt.Sprintf("**%s**: ", displayKey))

			if opts.KeysToCodeBlock[fullKeyPath] || opts.KeysToCodeBlock[key] {
				builder.WriteString("\n```\n")
				builder.WriteString(fmt.Sprintf("%v", val))
				builder.WriteString("\n```\n")
				continue
			}

			switch child := val.(type) {
			case map[string]interface{}, []interface{}:
				builder.WriteString("\n")
				builder.WriteString(m.convertValue(child, indentLevel+1, fullKeyPath, opts))
			default:
				builder.WriteString(fmt.Sprintf("%s\n", m.formatValue(val, fullKeyPath, opts)))
			}
		}

	case []interface{}:
		for i, item := range v {
			builder.WriteString(fmt.Sprintf("%s%s ", indent, fmt.Sprintf(m.numberFormat, i+1)))
			switch child := item.(type) {
			case map[string]interface{}, []interface{}:
				builder.WriteString("\n")
				builder.WriteString(m.convertValue(child, indentLevel+1, parentKey, opts))
			default:
				builder.WriteString(fmt.Sprintf("%s\n", m.formatValue(item, parentKey, opts)))
			}
		}

	default:
		builder.WriteString(fmt.Sprintf("%s%s\n", indent, m.formatValue(v, parentKey, opts)))
	}

	return builder.String()
}

func (m *MarkdownConverter) shouldEscapeKey(keyPath string, opts *MarkdownOpts) bool {
	return opts.EscapeAllKeys || opts.KeysToEscape[keyPath]
}

func (m *MarkdownConverter) formatValue(value interface{}, keyPath string, opts *MarkdownOpts) string {
	switch v := value.(type) {
	case string:
		if m.shouldEscapeKey(keyPath, opts) {
			return m.escapeReplacer.Replace(v)
		}
		return v
	case nil:
		return "`null`"
	case bool:
		return fmt.Sprintf("`%t`", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
