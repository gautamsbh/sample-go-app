package shared

import (
	"regexp"
	"strconv"
)

// ParseUserIDFromPath parse user id from req url path
func ParseUserIDFromPath(rgx *regexp.Regexp, path string) (out int, err error) {
	var (
		paths = rgx.FindStringSubmatch(path)
	)

	val := paths[len(paths)-1]
	out, err = strconv.Atoi(val)
	return
}
