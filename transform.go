package gotest2rdf

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type Option func(c *config)

func WithBacklogSize(size int) Option { return func(c *config) { c.backlogSize = size } }

type config struct {
	backlogSize int
}

func Transform(input io.Reader, output io.Writer, opts ...Option) error {
	cfg := new(config)
	for _, o := range opts {
		o(cfg)
	}
	if cfg.backlogSize == 0 {
		cfg.backlogSize = 3
	}

	outputs := make([]TestEvent, 0, cfg.backlogSize+1)
	dec := json.NewDecoder(input)
	enc := json.NewEncoder(output)
	for {
		var ev TestEvent
		err := dec.Decode(&ev)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		switch ev.Action {
		case "output":
			outputs = append(outputs, ev)
			if len(outputs) > cfg.backlogSize {
				outputs = outputs[1:]
			}
		case "fail":
			msg, loc := accum(outputs)
			if loc == nil {
				slog.Warn("parsed test event but the corresponding location is not found")
				continue
			}
			diag := RDFDiagnostic{Severity: RDFSeverityError, Location: loc, Message: msg}
			if err := enc.Encode(diag); err != nil {
				return err
			}
			outputs = outputs[:]
		case "skip":
			msg, loc := accum(outputs)
			if loc == nil {
				slog.Warn("parsed test event but the corresponding location is not found")
				continue
			}
			diag := RDFDiagnostic{Severity: RDFSeverityInfo, Location: loc, Message: msg}
			if err := enc.Encode(diag); err != nil {
				return err
			}
			outputs = outputs[:]
		case "run":
			// no-op
		default:
			outputs = outputs[:] // otherwise reset outputs
		}
	}
	return nil
}

func accum(outputs []TestEvent) (string, *RDFLocation) {
	msg := new(strings.Builder)
	var loc *RDFLocation
	var foundLoc bool
	for _, oe := range outputs {
		msg.Grow(len(oe.Output))
		msg.WriteString(oe.Output)
		if !foundLoc {
			loc, foundLoc = extractLocationFromOutput(oe.Output)
		}
	}
	return msg.String(), loc
}

func extractLocationFromOutput(s string) (*RDFLocation, bool) {
	parts := strings.SplitN(s, ":", 3)
	if len(parts) < 3 {
		return nil, false
	}
	loc := &RDFLocation{Path: strings.TrimSpace(parts[0])}
	lineNum, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, false
	}
	loc.Range = &RDFRange{Start: &RDFPosition{Line: lineNum}}
	return loc, true
}

type TestEvent struct {
	Time    time.Time
	Action  string
	Package string
	Test    string
	Elapsed float64
	Output  string
}

type RDFSeverity int

const (
	RDFSeverityUnknown RDFSeverity = iota
	RDFSeverityError
	RDFSeverityWarning
	RDFSeverityInfo
)

type RDFPosition struct {
	Line   int `json:"line,omitempty"`
	Column int `json:"column,omitempty"`
}

type RDFRange struct {
	Start *RDFPosition `json:"start,omitempty"`
	End   *RDFPosition `json:"end,omitempty"`
}

type RDFLocation struct {
	Path  string    `json:"path"`
	Range *RDFRange `json:"range"`
}

// RDFDiagnostic is an diagnostic representation of Reviewdog Diagnostic Format.
//
// refs. https://github.com/reviewdog/reviewdog/blob/master/proto/rdf/jsonschema/Diagnostic.jsonschema
type RDFDiagnostic struct {
	Message  string       `json:"message"`
	Severity RDFSeverity  `json:"severity"`
	Location *RDFLocation `json:"location"`
}
