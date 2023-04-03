package demo

import "github.com/stretchr/testify/suite"

type MyTestSuite struct {
	suite.Suite
}

func (m *MyTestSuite) SetupTest() {
	// TODO implement me
	panic("implement me")
}

func (m *MyTestSuite) SetupSuite() {
	// TODO implement me
	panic("implement me")
}

func (m *MyTestSuite) AfterTest(suiteName, testName string) {
	// TODO implement me
	panic("implement me")
}
