package similarity

import (
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
)

const threshold = 0.95

func IsSimilar(article1, article2 string) bool {
	swg := &metrics.SmithWatermanGotoh{
		CaseSensitive: false,
		GapPenalty:    -0.2,
		Substitution: metrics.MatchMismatch{
			Match:    1,
			Mismatch: -1,
		},
	}

	sim := strutil.Similarity(article1, article2, swg)

	return sim >= threshold
}
