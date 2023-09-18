/* 
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
	"crypto/rand"
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	yamlTestFilename = "foo.yaml"
	jsonTestFilename = "foo.json"
)

var (
	validYaml = `
foo: bar
`
	validJson = `
{"this":{"is":"json"}}
`
)

type structWithYaml struct {
	Content interface{} `valid:"isYaml"`
}

type isYamlSuite struct {
	suite.Suite
}

func (suite *isYamlSuite) TestValidYamlAsByteSlice() {
	assert.True(suite.T(), isYaml([]byte(validYaml), nil))
}

func (suite *isYamlSuite) TestValidYamlAsString() {
	assert.True(suite.T(), isYaml(validYaml, nil), nil)
}

func (suite *isYamlSuite) TestIsYaml() {
	testCases := []struct {
		desc               string
		yaml               interface{}
		expectedToValidate bool
	}{
		{
			desc:               "valid string",
			yaml:               validYaml,
			expectedToValidate: true,
		},
		{
			desc:               "valid slice",
			yaml:               []byte(validYaml),
			expectedToValidate: true,
		},
		{
			desc:               "invalid type: int",
			yaml:               int(42),
			expectedToValidate: false,
		},
	}
	for _, tC := range testCases {
		suite.T().Run(tC.desc, func(t *testing.T) {
			assert.True(t, isYaml(tC.yaml, nil) == tC.expectedToValidate)
		})
	}
}

func (suite *isYamlSuite) TestIsYamlStructValidation() {
	testCases := []struct {
		desc               string
		theStruct          structWithYaml
		expectedToValidate bool
	}{
		{
			desc:               "valid data as slice",
			theStruct:          structWithYaml{Content: []byte(validYaml)},
			expectedToValidate: true,
		},
		{
			desc:               "valid data as string",
			theStruct:          structWithYaml{Content: validYaml},
			expectedToValidate: true,
		},
		{
			desc:               "invalid data as slice",
			theStruct:          structWithYaml{Content: []byte(validJson)},
			expectedToValidate: false,
		},
		{
			desc:               "invalid data as string",
			theStruct:          structWithYaml{Content: validJson},
			expectedToValidate: false,
		},
		{
			desc:               "int as yaml",
			theStruct:          structWithYaml{Content: 3},
			expectedToValidate: false,
		},
	}
	for _, tC := range testCases {
		suite.T().Run(tC.desc, func(t *testing.T) {
			valid, err := govalidator.ValidateStruct(tC.theStruct)
			assert.True(t, valid == tC.expectedToValidate)
			if tC.expectedToValidate {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

type isAcceptedFileNameSuite struct {
	suite.Suite
}

func (suite *isAcceptedFileNameSuite) TestFilenameValidation() {
	testCases := []struct {
		desc     string
		filename interface{}
	}{
		{
			desc:     "YAML extension, lowercase",
			filename: yamlTestFilename,
		},
		{
			desc:     "YAML extension, uppercase",
			filename: "foo.YAML",
		},
		{
			desc:     "YML extension, lowercase",
			filename: "foo.yml",
		},
		{
			desc:     "YML extension, uppercase",
			filename: "foo.YML",
		},
		{
			desc:     "JSON extension, lowercase",
			filename: jsonTestFilename,
		},
		{
			desc:     "JSON extension, uppercase",
			filename: jsonTestFilename,
		},
	}
	for _, tC := range testCases {
		suite.T().Run(tC.desc, func(t *testing.T) {
			assert.True(t, isAcceptedFileName(tC.filename, nil))
		})
	}
}

func (suite *isAcceptedFileNameSuite) TestInvalidFilenameValidation() {
	testCases := []struct {
		desc     string
		filename interface{}
	}{
		{
			desc:     "No extension",
			filename: "foo",
		},
		{
			desc:     "Invalid extension",
			filename: "foo.bar",
		},
		{
			desc:     "Invalid extension, uppercase",
			filename: "foo.BAR",
		},
		{
			desc:     "Not a string",
			filename: 42,
		},
	}

	for _, tC := range testCases {
		suite.T().Run(tC.desc, func(t *testing.T) {
			assert.False(t, isAcceptedFileName(tC.filename, nil))
		})
	}
}

func TestValidation(t *testing.T) {
	// suite.Run(t, new(ValidationTestSuite))
	suite.Run(t, new(isYamlSuite))
	suite.Run(t, new(CorrectContentTestSuite))
	suite.Run(t, new(isAcceptedFileNameSuite))
}

type CorrectContentTestSuite struct {
	suite.Suite
}

func (suite *CorrectContentTestSuite) TestWithoutHandler() {
	assert.False(suite.T(), isCorrectContent(validYaml, struct{}{}))
}

func (suite *CorrectContentTestSuite) TestStringFunction() {
	testCases := []struct {
		desc               string
		content            string
		filename           string
		expectedToValidate bool
	}{
		{
			desc:               "YAML data with .yaml filename",
			content:            validYaml,
			filename:           yamlTestFilename,
			expectedToValidate: true,
		},
		{
			desc:               "YAML data with .json filename",
			content:            validYaml,
			filename:           jsonTestFilename,
			expectedToValidate: false,
		},
		{
			desc:               "JSON data with .json filename",
			content:            validJson,
			filename:           jsonTestFilename,
			expectedToValidate: true,
		},
		{
			desc:               "JSON data with .yaml filename",
			content:            validJson,
			filename:           yamlTestFilename,
			expectedToValidate: false,
		},
		{
			desc:               "Grabage data pretending to be YAML",
			content:            randString(100),
			filename:           yamlTestFilename,
			expectedToValidate: false,
		},
		{
			desc:               "Grabage data pretending to be JSON",
			content:            randString(100),
			filename:           jsonTestFilename,
			expectedToValidate: false,
		},
	}
	for _, tC := range testCases {
		suite.T().Run(tC.desc, func(t *testing.T) {
			ui := SwaggerUi{
				specFilename: tC.filename,
				specContent:  []byte(tC.content),
			}

			assert.True(t, isCorrectContent(tC.content, ui) == tC.expectedToValidate, "%s validates", tC.filename)
		})
	}
}

func (suite *CorrectContentTestSuite) TestFunction() {
	testCases := []struct {
		desc               string
		content            []byte
		filename           string
		expectedToValidate bool
	}{
		{
			desc:               "YAML data with .yaml filename",
			content:            []byte(validYaml),
			filename:           yamlTestFilename,
			expectedToValidate: true,
		},
		{
			desc:               "YAML data with .json filename",
			content:            []byte(validYaml),
			filename:           jsonTestFilename,
			expectedToValidate: false,
		},
		{
			desc:               "JSON data with .json filename",
			content:            []byte(validJson),
			filename:           jsonTestFilename,
			expectedToValidate: true,
		},
		{
			desc:               "JSON data with .yaml filename",
			content:            []byte(validJson),
			filename:           yamlTestFilename,
			expectedToValidate: false,
		},
		{
			desc:               "Grabage data pretending to be YAML",
			content:            []byte(randString(100)),
			filename:           yamlTestFilename,
			expectedToValidate: false,
		},
		{
			desc:               "Grabage data pretending to be JSON",
			content:            []byte(randString(100)),
			filename:           jsonTestFilename,
			expectedToValidate: false,
		},
	}
	for _, tC := range testCases {
		suite.T().Run(tC.desc, func(t *testing.T) {
			ui := SwaggerUi{
				specFilename: tC.filename,
				specContent:  tC.content,
			}

			assert.True(t, isCorrectContent(tC.content, ui) == tC.expectedToValidate, "%s validates", tC.filename)
		})
	}
}

func (suite *CorrectContentTestSuite) TestInvalidType() {
	assert.False(suite.T(), isCorrectContent(42, SwaggerUi{}))
}

func (suite *CorrectContentTestSuite) TestInvalidFilename() {
	assert.False(suite.T(), isCorrectContent(validYaml, SwaggerUi{specFilename: "foo.bar"}))
}

func randString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return string(b)
}
