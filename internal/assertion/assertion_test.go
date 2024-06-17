package assertion_test

import (
	"cotton/internal/assertion"
	"cotton/internal/line"
	"testing"

	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
)

func TestParseEqualAssertionWithInteger(t *testing.T) {
	mdLine := line.Line("* `$.form.value`==`30`")

	expected := &assertion.Assertion{
		Selector: "$.form.value",
		Value:    float64(30),
		Operator: mo.NewEither3Arg3[assertion.UndefinedOperator, assertion.UnaryAssertionOperator, assertion.BinaryAssertionOperator](&assertion.EqAssertion{}),
	}

	result, ok := assertion.Try(mdLine)

	assert.True(t, ok)
	assert.True(t, expected.SimilarTo(result))
}

func TestParseGreaterThanAssertionWithInteger(t *testing.T) {
	mdLine := line.Line("* `$.form.value`>`30`")

	expected := &assertion.Assertion{
		Selector: "$.form.value",
		Value:    float64(30),
		Operator: mo.NewEither3Arg3[assertion.UndefinedOperator, assertion.UnaryAssertionOperator, assertion.BinaryAssertionOperator](&assertion.GtAssertion{}),
	}

	result, ok := assertion.Try(mdLine)

	assert.True(t, ok)
	assert.True(t, expected.SimilarTo(result))
}

func TestParseEqualAssertionWithString(t *testing.T) {
	mdLine := line.Line("* `$.form.value`==`\"30\"`")

	expected := &assertion.Assertion{
		Selector: "$.form.value",
		Value:    "30",
		Operator: mo.NewEither3Arg3[assertion.UndefinedOperator, assertion.UnaryAssertionOperator, assertion.BinaryAssertionOperator](&assertion.EqAssertion{}),
	}

	result, ok := assertion.Try(mdLine)

	assert.True(t, ok)
	assert.True(t, expected.SimilarTo(result))
}
