// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package report

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/google/pprof/internal/binutils"
	"github.com/google/pprof/internal/proftest"
	"github.com/google/pprof/profile"
)

type testcase struct {
	rpt  *Report
	want string
}

func TestSource(t *testing.T) {
	const path = "testdata/"

	sampleValue1 := func(v []int64) int64 {
		return v[1]
	}

	for _, tc := range []testcase{
		{
			rpt: New(
				testProfile.Copy(),
				&Options{
					OutputFormat: List,
					Symbol:       regexp.MustCompile(`.`),
					Title:        filepath.Base(testProfile.Mapping[0].File),

					SampleValue: sampleValue1,
					SampleUnit:  testProfile.SampleType[1].Unit,
				},
			),
			want: path + "source.rpt",
		},
		{
			rpt: New(
				testProfile.Copy(),
				&Options{
					OutputFormat: Dot,
					CallTree:     true,
					Symbol:       regexp.MustCompile(`.`),
					Title:        filepath.Base(testProfile.Mapping[0].File),

					SampleValue: sampleValue1,
					SampleUnit:  testProfile.SampleType[1].Unit,
				},
			),
			want: path + "source.dot",
		},
	} {
		b := bytes.NewBuffer(nil)
		if err := Generate(b, tc.rpt, &binutils.Binutils{}); err != nil {
			t.Fatalf("%s: %v", tc.want, err)
		}

		gold, err := ioutil.ReadFile(tc.want)
		if err != nil {
			t.Fatalf("%s: %v", tc.want, err)
		}
		if string(b.String()) != string(gold) {
			d, err := proftest.Diff(gold, b.Bytes())
			if err != nil {
				t.Fatalf("%s: %v", "source", err)
			}
			t.Error("source" + "\n" + string(d) + "\n" + "gold:\n" + tc.want)
		}
	}
}

var testM = []*profile.Mapping{
	{
		ID:              1,
		HasFunctions:    true,
		HasFilenames:    true,
		HasLineNumbers:  true,
		HasInlineFrames: true,
	},
}

var testF = []*profile.Function{
	{
		ID:       1,
		Name:     "main",
		Filename: "testdata/source1",
	},
	{
		ID:       2,
		Name:     "foo",
		Filename: "testdata/source1",
	},
	{
		ID:       3,
		Name:     "bar",
		Filename: "testdata/source1",
	},
	{
		ID:       4,
		Name:     "tee",
		Filename: "testdata/source2",
	},
}

var testL = []*profile.Location{
	{
		ID:      1,
		Mapping: testM[0],
		Line: []profile.Line{
			{
				Function: testF[0],
				Line:     2,
			},
		},
	},
	{
		ID:      2,
		Mapping: testM[0],
		Line: []profile.Line{
			{
				Function: testF[1],
				Line:     4,
			},
		},
	},
	{
		ID:      3,
		Mapping: testM[0],
		Line: []profile.Line{
			{
				Function: testF[2],
				Line:     10,
			},
		},
	},
	{
		ID:      4,
		Mapping: testM[0],
		Line: []profile.Line{
			{
				Function: testF[3],
				Line:     2,
			},
		},
	},
	{
		ID:      5,
		Mapping: testM[0],
		Line: []profile.Line{
			{
				Function: testF[3],
				Line:     8,
			},
		},
	},
}

var testProfile = &profile.Profile{
	PeriodType:    &profile.ValueType{Type: "cpu", Unit: "millisecond"},
	Period:        10,
	DurationNanos: 10e9,
	SampleType: []*profile.ValueType{
		{Type: "samples", Unit: "count"},
		{Type: "cpu", Unit: "cycles"},
	},
	Sample: []*profile.Sample{
		{
			Location: []*profile.Location{testL[0]},
			Value:    []int64{1, 1},
		},
		{
			Location: []*profile.Location{testL[2], testL[1], testL[0]},
			Value:    []int64{1, 10},
		},
		{
			Location: []*profile.Location{testL[4], testL[2], testL[0]},
			Value:    []int64{1, 100},
		},
		{
			Location: []*profile.Location{testL[3], testL[0]},
			Value:    []int64{1, 1000},
		},
		{
			Location: []*profile.Location{testL[4], testL[3], testL[0]},
			Value:    []int64{1, 10000},
		},
	},
	Location: testL,
	Function: testF,
	Mapping:  testM,
}