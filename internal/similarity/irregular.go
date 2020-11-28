package similarity

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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
		return fmt.Errorf("failed to open file=%s: %w", irregularVerbFilePath, err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("file close failed: %v", err)
		}
	}()

	reader := csv.NewReader(file)

	v.verbs = make(map[string]irregularVerb, irregularVerbs)

	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return fmt.Errorf("failed to read: %w", err)
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
