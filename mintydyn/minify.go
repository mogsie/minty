package mintydyn

import (
	"regexp"
	"strings"
)

// MinifyJS reduces JavaScript size by removing unnecessary whitespace and comments.
// This is a lightweight minifier suitable for the generated runtime code.
func MinifyJS(js string) string {
	// Remove single-line comments (but preserve URLs)
	// Match // comments that aren't part of http:// or https://
	singleLineComment := regexp.MustCompile(`(^|[^:])//[^\n]*`)
	js = singleLineComment.ReplaceAllString(js, "$1")
	
	// Remove multi-line comments
	multiLineComment := regexp.MustCompile(`/\*[\s\S]*?\*/`)
	js = multiLineComment.ReplaceAllString(js, "")
	
	// Normalize line endings
	js = strings.ReplaceAll(js, "\r\n", "\n")
	
	// Remove leading/trailing whitespace from lines and collapse empty lines
	lines := strings.Split(js, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	
	// Join lines - use newlines to preserve statement separation
	js = strings.Join(result, "\n")
	
	// Collapse multiple newlines to single
	multiNewline := regexp.MustCompile(`\n{2,}`)
	js = multiNewline.ReplaceAllString(js, "\n")
	
	// Remove spaces around specific operators (safe ones that don't need space)
	// Be careful not to break things like "return value" or "new Object"
	
	// Remove space before these
	js = regexp.MustCompile(`\s+([{}\[\]();,:])`).ReplaceAllString(js, "$1")
	
	// Remove space after these
	js = regexp.MustCompile(`([{}\[\](;,:])\s+`).ReplaceAllString(js, "$1")
	
	// Remove space around = but be careful with == and ===
	js = regexp.MustCompile(`\s*([^=!<>])(=)([^=])\s*`).ReplaceAllStringFunc(js, func(s string) string {
		return strings.TrimSpace(s)
	})
	
	// Collapse multiple spaces to single (but don't touch newlines yet)
	multiSpace := regexp.MustCompile(`[ \t]{2,}`)
	js = multiSpace.ReplaceAllString(js, " ")
	
	// Now convert newlines to spaces where safe
	// Keep newlines after { and before }
	js = regexp.MustCompile(`\{\n`).ReplaceAllString(js, "{")
	js = regexp.MustCompile(`\n\}`).ReplaceAllString(js, "}")
	
	// Convert remaining newlines to spaces
	js = strings.ReplaceAll(js, "\n", " ")
	
	// Clean up any double spaces that resulted
	js = regexp.MustCompile(`\s{2,}`).ReplaceAllString(js, " ")
	
	// Remove space after opening and before closing parens/brackets
	js = regexp.MustCompile(`\(\s+`).ReplaceAllString(js, "(")
	js = regexp.MustCompile(`\s+\)`).ReplaceAllString(js, ")")
	js = regexp.MustCompile(`\[\s+`).ReplaceAllString(js, "[")
	js = regexp.MustCompile(`\s+\]`).ReplaceAllString(js, "]")
	
	// Remove space around arrows
	js = regexp.MustCompile(`\s*=>\s*`).ReplaceAllString(js, "=>")
	
	// Clean up semicolons
	js = regexp.MustCompile(`\s*;\s*`).ReplaceAllString(js, ";")
	
	// Clean up commas
	js = regexp.MustCompile(`\s*,\s*`).ReplaceAllString(js, ",")
	
	// Clean up colons (but careful with ternary)
	js = regexp.MustCompile(`\s*:\s*`).ReplaceAllString(js, ":")
	
	// Ensure script tags are readable
	js = strings.ReplaceAll(js, "<script>", "<script>\n")
	js = strings.ReplaceAll(js, "</script>", "\n</script>")
	
	// Final cleanup
	js = strings.TrimSpace(js)
	
	return js
}

