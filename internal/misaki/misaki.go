// Package misaki holds what the normalizer and the markdown pass need to know
// about Misaki, the phonemizer Kokoro uses.
package misaki

import "regexp"

// Link matches Misaki's inline link syntax, \[([^\]]+)\]\(([^\)]*)\), with
// the text and destination as submatches.
var Link = regexp.MustCompile(`\[([^\]]+)\]\(([^)]*)\)`)

// override matches the link destinations Misaki's preprocess accepts: a
// stress level, a phoneme string between slashes or a feature between
// hashes. Other links, such as URLs, are not pronunciation overrides.
var override = regexp.MustCompile(`^(?:[+-]?\d+|[+-]?0\.5|/.+/|#.+#)$`)

// IsOverride reports whether a link destination makes the link a
// pronunciation override.
func IsOverride(destination string) bool { return override.MatchString(destination) }
