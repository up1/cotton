package testcase

import (
	"cotton/internal/assertion"
	"cotton/internal/capture"
	"cotton/internal/console"
	"cotton/internal/executable"
	"cotton/internal/request"
	"cotton/internal/response"
	"errors"
	"fmt"
	"slices"
)

// Test cases
type TestCase struct {
	Title       string
	Description string
	Request     *request.Request

	Captures   []*capture.Capture
	Setups     []*executable.Executable
	Teardowns  []*executable.Executable
	Assertions []*assertion.Assertion
}

func (t *TestCase) Execute(logger console.Console) *TestResult {
	if logger == nil {
		logger = console.NewNilConsole()
	}

	testResult := &TestResult{
		Title:      t.Title,
		Passed:     false,
		Assertions: []AssertionResult{},
		Error:      nil,
	}

	if t.Request == nil {
		testResult.Error = errors.New("no request to be made")
		return testResult
	}

	logger.Printfln("%s", t.Title)

	for _, setup := range t.Setups {
		_, err := setup.Execute(logger)
		if err != nil {
			testResult.Error = err
			return testResult
		}
	}

	r, err := t.Request.Do()
	if err != nil {
		testResult.Error = err
		return testResult
	}
	defer r.Body.Close()

	resp, err := response.New(r)
	if err != nil {
		testResult.Error = err
		return testResult
	}

	for _, assertion := range t.Assertions {
		actual, err := resp.ValueOf(assertion.Selector)
		if err != nil {
			testResult.Error = err
			return testResult
		}
		expected := assertion.Value
		if assertion.Operator.IsArg1() {
			testResult.Error = errors.New("unexpected assertion found")
			return testResult
		}
		var result bool
		if assertion.Operator.IsArg2() {
			result, err = assertion.Operator.Arg2OrEmpty().Assert(actual)
			if err != nil {
				testResult.Error = err
				return testResult
			}
		} else {
			if assertion.Operator.IsArg3() {
				result, err = assertion.Operator.Arg3OrEmpty().Assert(expected, actual)
			}
			if err != nil {
				testResult.Error = err
				return testResult
			}
		}
		testResult.Assertions = append(testResult.Assertions, AssertionResult{
			Title:    assertion.String(),
			Passed:   result,
			Actual:   fmt.Sprintf("%v", actual),
			Expected: fmt.Sprintf("%v", expected),
		})
	}

	for _, teardown := range t.Teardowns {
		_, err := teardown.Execute(logger)
		if err != nil {
			testResult.Error = err
			return testResult
		}
	}

	testResult.Passed = true
	testResult.Error = nil
	return testResult
}

func (t *TestCase) SimilarTo(anotherTestCase *TestCase) bool {
	return t.Title == anotherTestCase.Title &&
		t.Description == anotherTestCase.Description &&
		t.Request.Similar(anotherTestCase.Request) &&
		slices.EqualFunc(t.Captures, anotherTestCase.Captures, func(c1, c2 *capture.Capture) bool {
			return c1.SimilarTo(c2)
		}) &&
		slices.EqualFunc(t.Setups, anotherTestCase.Setups, func(s1, s2 *executable.Executable) bool {
			return s1.SimilarTo(s2)
		}) &&
		slices.EqualFunc(t.Teardowns, anotherTestCase.Teardowns, func(s1, s2 *executable.Executable) bool {
			return s1.SimilarTo(s2)
		}) &&
		slices.EqualFunc(t.Assertions, anotherTestCase.Assertions, func(a1, a2 *assertion.Assertion) bool {
			return a1.SimilarTo(a2)
		})
}
