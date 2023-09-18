package swaggerui

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SwaggerInitializerSuite struct {
	suite.Suite
}

func (suite *SwaggerInitializerSuite) TestInitializer() {
	assert.NotPanics(suite.T(), func() {
		s := getInitializer("foo.yaml", "/bar/baz")
		assert.NotEmpty(suite.T(), s)
		assert.Contains(suite.T(), string(s), "/bar/baz/foo.yaml")
	})
}

func (suite *SwaggerInitializerSuite) TestInitializerTmpl() {
	testCases := []struct {
		desc     string
		prefix   string
		filename string
		expected string
	}{
		{
			desc:     "With Prefix",
			prefix:   "/foo/bar",
			filename: "baz.yaml",
			expected: "/foo/bar/baz.yaml",
		},
		{
			desc:     "Without Prefix",
			prefix:   "",
			filename: "baz.yaml",
			expected: "./baz.yaml",
		},
	}
	for _, tC := range testCases {
		suite.T().Run(tC.desc, func(t *testing.T) {
			assert.NotPanics(t, func() {
				s := string(getInitializer(tC.filename, tC.prefix))
				assert.NotEmpty(t, s)
				assert.Contains(t, s, tC.expected)
			})
		})
	}
}

func TestInitializerSuite(t *testing.T) {
	suite.Run(t, new(SwaggerInitializerSuite))
}
