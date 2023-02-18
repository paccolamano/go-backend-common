package domain

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ErrorSuite struct {
	suite.Suite
	message string
}

func (suite *ErrorSuite) SetupTest() {
	suite.message = "foo bar"
}

func TestErrorSuite(t *testing.T) {
	suite.Run(t, new(ErrorSuite))
}

func (suite *ErrorSuite) TestConstant() {
	suite.Equal("one or more of the services is null", ServiceNilError)

	suite.Equal("sending http %s request to %s with body %v", HttpRequestInfoMessage)
	suite.Equal("received http reply with status code %d and body %v", HttpResponseInfoMessage)
	suite.Equal("error while sending http %s request to %s caused by %s", HttpErrorMessage)

	suite.Equal("cannot create %s caused by %s", EntityCreateError)
	suite.Equal("cannot get %s caused by %s", EntityGetError)
	suite.Equal("cannot update %s caused by %s", EntityUpdateError)
	suite.Equal("cannot delete %s caused by %s", EntityDeleteError)
	suite.Equal("cannot find %s caused by %s", EntityFindError)
	suite.Equal("%s with uuid %s not found", EntityNotFoundError)
}

func (suite *ErrorSuite) TestNewRestError() {
	tests := []struct {
		data        *RestError
		expectedErr *RestError
	}{
		{
			NewRestError(500, "foo %s", "bar"),
			&RestError{suite.message, 500},
		},
		{
			NewBadRequestError("foo %s", "bar"),
			&RestError{suite.message, 400},
		},
		{
			NewUnauthorizedError("foo %s", "bar"),
			&RestError{suite.message, 401},
		},
		{
			NewNotFoundError("foo %s", "bar"),
			&RestError{suite.message, 404},
		},
		{
			NewInternalServerError("foo %s", "bar"),
			&RestError{suite.message, 500},
		},
	}

	for _, test := range tests {
		suite.Equal(test.expectedErr, test.data)
	}
}
