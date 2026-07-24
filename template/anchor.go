/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package template

import (
	"fmt"
)

// CreateAnchorMarkdown generates both an <a> anchor tag and a clickable link so
// that inputs/outputs can be cross-referenced from elsewhere in a document.
// The dual approach (anchor + link) ensures the item is both a jump target and
// itself clickable for navigation.
func CreateAnchorMarkdown(prefix string, value string, anchor bool, escape bool) string {
	sanitizedName := SanitizeName(value, escape)

	if anchor {
		anchorName := fmt.Sprintf("%s_%s", prefix, value)
		sanitizedAnchorName := SanitizeName(anchorName, escape)
		// the <a> link is purposely not sanitized as this breaks markdown formatting
		return fmt.Sprintf("<a name=\"%s\"></a> [%s](#%s)", anchorName, sanitizedName, sanitizedAnchorName)
	}

	return sanitizedName
}

// CreateAnchorAsciidoc is the AsciiDoc equivalent using [[id]] and <<id,label>>
// syntax. AsciiDoc uses different anchor/xref conventions than Markdown, so a
// separate function keeps format-specific concerns isolated.
func CreateAnchorAsciidoc(prefix string, value string, anchor bool, escape bool) string {
	sanitizedName := SanitizeName(value, escape)

	if anchor {
		anchorName := fmt.Sprintf("%s_%s", prefix, value)
		sanitizedAnchorName := SanitizeName(anchorName, escape)
		return fmt.Sprintf("[[%s]] <<%s,%s>>", sanitizedAnchorName, sanitizedAnchorName, sanitizedName)
	}

	return sanitizedName
}
