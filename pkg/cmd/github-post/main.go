// Copyright 2016 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

// Command github-post parses the JSON-formatted output from a Go test session,
// as generated by either 'go test -json' or './pkg.test | go tool test2json -t',
// and posts issues for any failed tests to GitHub. If there are no failed
// tests, it assumes that there was a build error and posts the entire log to
// GitHub.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/cockroach/pkg/cmd/internal/issues"
	"github.com/pkg/errors"
)

const (
	pkgEnv = "PKG"
)

func main() {
	ctx := context.Background()

	f := func(ctx context.Context, title, packageName, testName, testMessage, authorEmail string) error {
		log.Printf("filing issue with title: %s", title)
		return issues.Post(ctx, title, packageName, testName, testMessage, authorEmail, nil)
	}

	if err := listFailures(ctx, os.Stdin, f); err != nil {
		log.Fatal(err)
	}
}

// This struct is described in the test2json documentation.
// https://golang.org/cmd/test2json/
type testEvent struct {
	Action  string
	Test    string
	Output  string
	Time    time.Time // encodes as an RFC3339-format string
	Elapsed float64   // seconds
}

func listFailures(
	ctx context.Context,
	input io.Reader,
	f func(ctx context.Context, title, packageName, testName, testMessage, authorEmail string) error,
) error {
	// Tests that took less than this are not even considered for slow test
	// reporting. This is so that we protect against large number of
	// programatically-generated subtests.
	const shortTestFilterSecs float64 = 0.5
	var timeoutMsg = "panic: test timed out after"

	dec := json.NewDecoder(input)

	packageName, ok := os.LookupEnv(pkgEnv)
	if !ok {
		return errors.Errorf("package name environment variable %s is not set", pkgEnv)
	}
	trimmedPkgName := strings.TrimPrefix(packageName, issues.CockroachPkgPrefix)

	var packageOutput bytes.Buffer

	// map  from test name to list of events (each log line is an event, plus
	// start and pass/fail events).
	// Tests/events are "outstanding" until we see a final pass/fail event.
	// Because of the way the go tet runner prints output, in case a subtest times
	// out or panics, we don't get a  pass/fail event for sibling and ancestor
	// tests. Those tests will remain "outstanding" and will be ignored for the
	// purpose of issue reporting.
	outstandingOutput := make(map[string][]testEvent)
	failures := make(map[string][]testEvent)
	var slowPassingTests []testEvent
	var slowFailingTests []testEvent

	// init is true for the preamble of the input before the first "run" test
	// event.
	init := true
	// trustTimestamps will be set if we don't find a marker suggesting that the
	// input comes from a stress run. In that case, stress prints all its output
	// at once (for a captured failed test run), so the test2json timestamps are
	// meaningless.
	trustTimestamps := true
	// elapsedTotalSec accumulates the time spent in all tests, passing or
	// failing. In case the input comes from a stress run, this will be used to
	// deduce the duration of a timed out test.
	var elapsedTotalSec float64
	// Will be set if the last test timed out.
	var timedOutTestName string
	var timedOutEvent testEvent
	var curTestStart time.Time
	var lastTestName string
	var lastEvent testEvent
	for {
		var te testEvent
		if err := dec.Decode(&te); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		lastEvent = te

		if te.Test != "" {
			init = false
		}
		if init && strings.Contains(te.Output, "-exec 'stress '") {
			trustTimestamps = false
		}
		if timedOutTestName == "" && te.Elapsed > 0 {
			// We don't count subtests as those are counted in the parent.
			if split := strings.SplitN(te.Test, "/", 2); len(split) == 1 {
				elapsedTotalSec += te.Elapsed
			}
		}

		if timedOutTestName == te.Test && te.Elapsed != 0 {
			te.Elapsed = timedOutEvent.Elapsed
		}

		// Events for the overall package test do not set Test.
		if len(te.Test) > 0 {
			switch te.Action {
			case "run":
				lastTestName = te.Test
				if trustTimestamps {
					curTestStart = te.Time
				}
			case "output":
				outstandingOutput[te.Test] = append(outstandingOutput[te.Test], te)
				if strings.Contains(te.Output, timeoutMsg) {
					timedOutTestName = te.Test

					// Fill in the Elapsed field for a timeout event.
					// As of go1.11, the Elapsed field is bogus for fail events for timed
					// out tests, so we do our own computation.
					// See https://github.com/golang/go/issues/27568
					//
					// Also, if the input is coming from stress, there will not even be a
					// fail event for the test, so the Elapsed field computed here will be
					// useful.
					if trustTimestamps {
						te.Elapsed = te.Time.Sub(curTestStart).Seconds()
					} else {
						// If we don't trust the timestamps, then we compute the test's
						// duration by subtracting all the durations that we've seen so far
						// (which we do trust to some extent). Note that this is not
						// entirely accurate, since there's no information about the
						// duration about sibling subtests which may have run. And further
						// note that it doesn't work well at all for small timeouts because
						// the resolution that the test durations have is just tens of
						// milliseconds, so many quick tests are rounded of to a duration of
						// 0.
						re := regexp.MustCompile(`panic: test timed out after (\d*(?:\.\d*)?)(.)`)
						matches := re.FindStringSubmatch(te.Output)
						if matches == nil {
							log.Printf("failed to parse timeout message: %s", te.Output)
							te.Elapsed = -1
						} else {
							dur, err := strconv.ParseFloat(matches[1], 64)
							if err != nil {
								log.Fatal(err)
							}
							if matches[2] == "m" {
								// minutes to seconds
								dur *= 60
							} else if matches[2] != "s" {
								log.Fatalf("unexpected time unit in: %s", te.Output)
							}
							te.Elapsed = dur - elapsedTotalSec
						}
					}
					timedOutEvent = te
				}
			case "pass", "skip":
				if timedOutTestName != "" {
					panic(fmt.Sprintf("detected test timeout but test seems to have passed (%+v)", te))
				}
				delete(outstandingOutput, te.Test)
				if te.Elapsed > shortTestFilterSecs {
					// We ignore subtests; their time contributes to the parent's.
					if !strings.Contains(te.Test, "/") {
						slowPassingTests = append(slowPassingTests, te)
					}
				}
			case "fail":
				// Record slow tests. We ignore subtests; their time contributes to the
				// parent's. Except the timed out (sub)test, for which the parent (if
				// any) is not going to appear in the report because there's not going
				// to be a pass/fail event for it.
				if !strings.Contains(te.Test, "/") || timedOutTestName == te.Test {
					slowFailingTests = append(slowFailingTests, te)
				}
				// Move the test to the failures collection unless the test timed out.
				// We have special reporting for timeouts below.
				if timedOutTestName != te.Test {
					failures[te.Test] = outstandingOutput[te.Test]
				}
				delete(outstandingOutput, te.Test)
			}
		} else if te.Action == "output" {
			// Output was outside the context of a test. This consists mostly of the
			// preamble and epilogue that Make outputs, but also any log messages that
			// are printed by a test binary's main function.
			packageOutput.WriteString(te.Output)
		}
	}

	// On timeout, we might or might not have gotten a fail event for the timed
	// out test (we seem to get one when processing output from a test binary run,
	// but not when processing the output of `stress`, which adds some lines at
	// the end). If we haven't gotten a fail event, the test's output is still
	// outstanding and the test is not registered in the slowFailingTests
	// collection. The timeout handling code below relies on slowFailingTests not
	// being empty though, so we'll process the test here.
	if timedOutTestName != "" {
		if _, ok := outstandingOutput[timedOutTestName]; ok {
			slowFailingTests = append(slowFailingTests, timedOutEvent)
			delete(outstandingOutput, timedOutTestName)
		}
	} else {
		// If we haven't received a final event for the last test, then a
		// panic/log.Fatal must have happened. Consider it failed.
		// Note that because of https://github.com/golang/go/issues/27582 there
		// might be other outstanding tests; we ignore those.
		if _, ok := outstandingOutput[lastTestName]; ok {
			log.Printf("found outstanding output. Considering last test failed: %s", lastTestName)
			failures[lastTestName] = outstandingOutput[lastTestName]
		}
	}

	// test2json always puts a fail event last unless it sees a big pass message
	// from the test output.
	if lastEvent.Action == "fail" && len(failures) == 0 && timedOutTestName == "" {
		// If we couldn't find a failing Go test, assume that a failure occurred
		// before running Go and post an issue about that.
		const unknown = "(unknown)"
		title := fmt.Sprintf("%s: package failed under stress", trimmedPkgName)
		if err := f(
			ctx, title, packageName, unknown, packageOutput.String(), "", /* authorEmail */
		); err != nil {
			return errors.Wrap(err, "failed to post issue")
		}
	} else {
		for test, testEvents := range failures {
			if split := strings.SplitN(test, "/", 2); len(split) == 2 {
				parentTest, subTest := split[0], split[1]
				log.Printf("consolidating failed subtest %q into parent test %q", subTest, parentTest)
				failures[parentTest] = append(failures[parentTest], testEvents...)
				delete(failures, test)
			} else {
				log.Printf("failed parent test %q", test)
			}
		}
		// Sort the failed tests to make the unit tests for this script deterministic.
		var failedTestNames []string
		for name := range failures {
			failedTestNames = append(failedTestNames, name)
		}
		sort.Strings(failedTestNames)
		for _, test := range failedTestNames {
			testEvents := failures[test]
			authorEmail, err := getAuthorEmail(ctx, packageName, test)
			if err != nil {
				log.Printf("unable to determine test author email: %s\n", err)
			}
			var outputs []string
			for _, testEvent := range testEvents {
				outputs = append(outputs, testEvent.Output)
			}
			message := strings.Join(outputs, "")
			title := fmt.Sprintf("%s: %s failed under stress", trimmedPkgName, test)
			if err := f(ctx, title, packageName, test, message, authorEmail); err != nil {
				return errors.Wrap(err, "failed to post issue")
			}
		}
	}

	// Sort slow tests descendingly by duration.
	sort.Slice(slowPassingTests, func(i, j int) bool {
		return slowPassingTests[i].Elapsed > slowPassingTests[j].Elapsed
	})
	sort.Slice(slowFailingTests, func(i, j int) bool {
		return slowFailingTests[i].Elapsed > slowFailingTests[j].Elapsed
	})

	report := genSlowTestsReport(slowPassingTests, slowFailingTests)
	if err := writeSlowTestsReport(report); err != nil {
		log.Printf("failed to create slow tests report: %s", err)
	}

	// If the run timed out, file an issue. A couple of cases:
	// 1) If the test that was running when the package timed out is the longest
	// test, then we blame it. The common case is the test deadlocking - it would
	// have run forever.
	// 2) Otherwise, we don't blame anybody in particular. We file a generic issue
	// listing the package name containing the report of long-running tests.
	if timedOutTestName != "" {
		slowest := slowFailingTests[0]
		if len(slowPassingTests) > 0 && slowPassingTests[0].Elapsed > slowest.Elapsed {
			slowest = slowPassingTests[0]
		}
		if timedOutTestName == slowest.Test {
			// The test that was running when the timeout hit is the one that ran for
			// the longest time.
			authorEmail, err := getAuthorEmail(ctx, packageName, timedOutTestName)
			if err != nil {
				log.Printf("unable to determine test author email: %s\n", err)
			}
			title := fmt.Sprintf("%s: %s timed out under stress", trimmedPkgName, timedOutTestName)
			log.Printf("timeout culprit found: %s\n", timedOutTestName)
			if err := f(ctx, title, packageName, timedOutTestName, report, authorEmail); err != nil {
				return errors.Wrap(err, "failed to post issue")
			}
		} else {
			title := fmt.Sprintf("%s: package timed out under stress", trimmedPkgName)
			// Andrei gets these reports for now, but don't think I'll fix anything
			// you fools.
			// TODO(andrei): Figure out how to assign to the on-call engineer. Maybe
			// get their name from the Slack channel?
			log.Printf("timeout culprit not found\n")
			if err := f(
				ctx, title, packageName, "(unknown)" /* testName */, report, "andreimatei1@gmail.com",
			); err != nil {
				return errors.Wrap(err, "failed to post issue")
			}
		}
	}

	return nil
}

func genSlowTestsReport(slowPassingTests, slowFailingTests []testEvent) string {
	var b strings.Builder
	b.WriteString("Slow failing tests:\n")
	for i, te := range slowFailingTests {
		if i == 20 {
			break
		}
		fmt.Fprintf(&b, "%s - %.2fs\n", te.Test, te.Elapsed)
	}
	if len(slowFailingTests) == 0 {
		fmt.Fprint(&b, "<none>\n")
	}

	b.WriteString("\nSlow passing tests:\n")
	for i, te := range slowPassingTests {
		if i == 20 {
			break
		}
		fmt.Fprintf(&b, "%s - %.2fs\n", te.Test, te.Elapsed)
	}
	if len(slowPassingTests) == 0 {
		fmt.Fprint(&b, "<none>\n")
	}
	return b.String()
}

func writeSlowTestsReport(report string) error {
	return ioutil.WriteFile("artifacts/slow-tests-report.txt", []byte(report), 0644)
}

func getAuthorEmail(ctx context.Context, packageName, testName string) (string, error) {
	// Search the source code for the email address of the last committer to touch
	// the first line of the source code that contains testName. Then, ask GitHub
	// for the GitHub username of the user with that email address by searching
	// commits in cockroachdb/cockroach for commits authored by the address.
	subtests := strings.Split(testName, "/")
	testName = subtests[0]
	packageName = strings.TrimPrefix(packageName, "github.com/cockroachdb/cockroach/")
	cmd := exec.Command(`/bin/bash`, `-c`,
		fmt.Sprintf(`git grep -n "func %s" $(git rev-parse --show-toplevel)/%s/*_test.go`,
			testName, packageName))
	// This command returns output such as:
	// ../ccl/storageccl/export_test.go:31:func TestExportCmd(t *testing.T) {
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Errorf("couldn't find test %s in %s: %s %s",
			testName, packageName, err, string(out))
	}
	re := regexp.MustCompile(`(.*):(.*):`)
	// The first 2 :-delimited fields are the filename and line number.
	matches := re.FindSubmatch(out)
	if matches == nil {
		return "", errors.Errorf("couldn't find filename/line number for test %s in %s: %s",
			testName, packageName, string(out))
	}
	filename := matches[1]
	linenum := matches[2]

	// Now run git blame.
	cmd = exec.Command(`/bin/bash`, `-c`,
		fmt.Sprintf(`git blame --porcelain -L%s,+1 %s | grep author-mail`,
			linenum, filename))
	// This command returns output such as:
	// author-mail <jordan@cockroachlabs.com>
	out, err = cmd.CombinedOutput()
	if err != nil {
		return "", errors.Errorf("couldn't find author of test %s in %s: %s %s",
			testName, packageName, err, string(out))
	}
	re = regexp.MustCompile("author-mail <(.*)>")
	matches = re.FindSubmatch(out)
	if matches == nil {
		return "", errors.Errorf("couldn't find author email of test %s in %s: %s",
			testName, packageName, string(out))
	}
	return string(matches[1]), nil
}
