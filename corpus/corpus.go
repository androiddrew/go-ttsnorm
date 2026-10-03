// Package corpus is the normalizer's test corpus: each input with the output
// of KittenTTS's Python normalize_text and, where the Go port deliberately
// differs, the reviewed override. Other projects can use it to check how
// their speech pipeline reads the same inputs.
//
// normalize_corpus.txt lists the inputs. normalize_golden.json and
// normalize_overrides.yaml are kept in step with gokittentts, whose `make
// golden` regenerates the golden file from the Python reference.
package corpus

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"

	"go.yaml.in/yaml/v3"
)

//go:embed normalize_golden.json
var golden []byte

//go:embed normalize_overrides.yaml
var overrides []byte

// Case is one corpus input.
type Case struct {
	Input string
	// Python is normalize_text's output, or the error it raised.
	Python string
	// PythonFailed reports that Python raised an error.
	PythonFailed bool
	// Expected is the default normalize.Text output: Python's, or the
	// override's.
	Expected string
	// Override is the reviewed deviation from Python, if there is one.
	Override *Override
}

// Override is a deliberate difference between the Go normalizer and Python.
type Override struct {
	Input    string `yaml:"input"`
	Python   string `yaml:"python"`
	Expected string `yaml:"expected"`
	Reason   string `yaml:"reason"`
	Approved bool   `yaml:"approved"`
}

// Cases returns the corpus in order. It fails if an override is incomplete,
// duplicated, stale or names an input that isn't in the corpus, or if an
// input that Python failed on has no override.
func Cases() ([]Case, error) {
	var entries []struct {
		Input  string  `json:"input"`
		Output *string `json:"output"`
		Error  string  `json:"error"`
	}
	if err := json.Unmarshal(golden, &entries); err != nil {
		return nil, fmt.Errorf("normalize_golden.json: %w", err)
	}
	var list []Override
	if err := yaml.Unmarshal(overrides, &list); err != nil {
		return nil, fmt.Errorf("normalize_overrides.yaml: %w", err)
	}
	byInput := make(map[string]*Override, len(list))
	var errs []error
	for i := range list {
		o := &list[i]
		switch {
		case o.Reason == "":
			errs = append(errs, fmt.Errorf("override for %q: no reason", o.Input))
		case o.Python == "":
			errs = append(errs, fmt.Errorf("override for %q: no python output", o.Input))
		case o.Expected == "":
			errs = append(errs, fmt.Errorf("override for %q: no expected output", o.Input))
		}
		if _, dup := byInput[o.Input]; dup {
			errs = append(errs, fmt.Errorf("override for %q: listed twice", o.Input))
		}
		byInput[o.Input] = o
	}
	cases := make([]Case, 0, len(entries))
	for _, e := range entries {
		c := Case{Input: e.Input, Python: e.Error, PythonFailed: e.Output == nil}
		if e.Output != nil {
			c.Python, c.Expected = *e.Output, *e.Output
		}
		if o, ok := byInput[e.Input]; ok {
			if o.Python != c.Python {
				errs = append(errs, fmt.Errorf("override for %q is stale: it says Python gives %q, the golden says %q", e.Input, o.Python, c.Python))
			}
			c.Override, c.Expected = o, o.Expected
			delete(byInput, e.Input)
		} else if c.PythonFailed {
			errs = append(errs, fmt.Errorf("the Python reference raised %s on %q; the case needs an override", e.Error, e.Input))
		}
		cases = append(cases, c)
	}
	for input := range byInput {
		errs = append(errs, fmt.Errorf("override for %q: not in the corpus", input))
	}
	return cases, errors.Join(errs...)
}
