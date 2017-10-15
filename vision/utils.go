package vision

import "strings"

func urlJoin(base string, n string) string {
	if strings.HasPrefix(base, "/") {
		base = strings.TrimPrefix(base, "/")
	}
	return strings.Join([]string{base, n}, "/")
}
