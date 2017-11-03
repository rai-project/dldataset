package vision

import "strings"

var (
	synset = map[int]string{}
)

func init() {

	lines := strings.Split(_escFSMustString(false, "/vision/support/synset.txt"), "\n")
	for ii, line := range lines {
		synset[ii] = line
	}
}
