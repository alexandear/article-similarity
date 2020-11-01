package similarity

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/devchallenge/article-similarity/internal/util"
)

const (
	irregularVerbs = 200

	irregularForms = 3
)

type IrregularVerb struct {
	verbs map[string]irregularVerb
}

type irregularVerb struct {
	simplePast     string
	pastParticiple string
}

func (v *IrregularVerb) Load(irregularVerbFilePath string) error {
	file, err := os.Open(irregularVerbFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer util.Close(file, log.Printf)

	reader := csv.NewReader(file)

	v.verbs = make(map[string]irregularVerb, irregularVerbs)

	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return errors.Wrap(err, "failed to read")
		}

		if len(record) < irregularForms {
			continue
		}

		infinitive := record[0]
		simplePast := record[1]
		pastParticiple := record[2]

		v.verbs[infinitive] = irregularVerb{
			simplePast:     simplePast,
			pastParticiple: pastParticiple,
		}
	}

	return nil
}

func (v *IrregularVerb) ToInfinitive(verb string) string {
	for infinitive, irregular := range v.verbs {
		for _, verbForm := range []string{infinitive, irregular.simplePast, irregular.pastParticiple} {
			if strings.EqualFold(verbForm, verb) {
				return infinitive
			}
		}
	}

	return verb
}
