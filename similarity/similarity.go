package similarity

import (
	"reflect"
	"sort"
	"strings"
)

func IsSimilar(article1, article2 string) bool {
	const sep = " "
	words1 := strings.Split(article1, sep)
	words2 := strings.Split(article2, sep)
	sort.Strings(words1)
	sort.Strings(words2)
	return reflect.DeepEqual(words1, words2)
}
