package gotest2rdf_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/aereal/gotest2rdf"
)

func TestTransform(t *testing.T) {
	testCases := []struct {
		name    string
		want    []gotest2rdf.RDFDiagnostic
		wantErr error
	}{
		{
			name: "ok",
			want: []gotest2rdf.RDFDiagnostic{
				{Severity: gotest2rdf.RDFSeverityError, Location: &gotest2rdf.RDFLocation{Path: "test_test.go", Range: &gotest2rdf.RDFRange{Start: &gotest2rdf.RDFPosition{Line: 10}}}, Message: "=== RUN   Test_ng\n    test_test.go:10: failing\n--- FAIL: Test_ng (0.00s)\n"},
				{Severity: gotest2rdf.RDFSeverityError, Location: &gotest2rdf.RDFLocation{Path: "test_test.go", Range: &gotest2rdf.RDFRange{Start: &gotest2rdf.RDFPosition{Line: 14}}}, Message: "    test_test.go:14: failing line1\n        line2\n--- FAIL: Test_ng_multiline (0.00s)\n"},
				{Severity: gotest2rdf.RDFSeverityInfo, Location: &gotest2rdf.RDFLocation{Path: "test_test.go", Range: &gotest2rdf.RDFRange{Start: &gotest2rdf.RDFPosition{Line: 22}}}, Message: "=== RUN   Test_skip\n    test_test.go:22: skipped\n--- SKIP: Test_skip (0.00s)\n"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", "inputs", tc.name))
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			out := new(bytes.Buffer)
			gotErr := gotest2rdf.Transform(f, out)
			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("error:\n\twant: %v\n\t got: %v", tc.wantErr, gotErr)
			}
			if gotErr != nil {
				return
			}
			got, err := readLines(out)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("diagnostics:\n\twant: %#v\n\t got: %#v", tc.want, got)
			}
		})
	}
}

func readLines(out io.Reader) ([]gotest2rdf.RDFDiagnostic, error) {
	r := bufio.NewReader(out)
	var diags []gotest2rdf.RDFDiagnostic
	for {
		line, _, err := r.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		var diag gotest2rdf.RDFDiagnostic
		if err := json.Unmarshal(line, &diag); err != nil {
			return nil, err
		}
		diags = append(diags, diag)
	}
	return diags, nil
}
