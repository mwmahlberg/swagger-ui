/* 
 *  initializer_test.go is part of github.com/mwmahlberg/swagger-ui project.
 *  
 *  Copyright 2023 Markus W Mahlberg
 *  
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *  
 *      http://www.apache.org/licenses/LICENSE-2.0
 *  
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

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
