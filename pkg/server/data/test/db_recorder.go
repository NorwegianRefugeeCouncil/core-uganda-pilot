package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// DBRecorder is a testing utility to record SQL statements and parameters
// that were issued by the engine
type DBRecorder struct {
	statements []string
	params     [][]interface{}
}

// GetStatements returns the recorded statements
func (s *DBRecorder) GetStatements() []string {
	return s.statements
}

// GetParams returns the recorded parameters
func (s *DBRecorder) GetParams() [][]interface{} {
	return s.params
}

func (s *DBRecorder) Record(stmt string, params []interface{}) {
	s.statements = append(s.statements, stmt)
	s.params = append(s.params, params)
}

// Reset resets the recorder
func (s *DBRecorder) Reset() {
	s.statements = []string{}
	s.params = [][]interface{}{}
}

type ExpectedStatement struct {
	SQL    string
	Params []interface{}
}

func (s *DBRecorder) AssertStatementsExecuted(t *testing.T, expectStatements []ExpectedStatement) {
	actualStatements := s.GetStatements()
	actualParams := s.GetParams()
	assert.Equal(t, len(expectStatements), len(actualStatements))
	for i, statement := range expectStatements {
		assert.Equal(t, statement.SQL, actualStatements[i])
		assert.Equal(t, statement.Params, actualParams[i])
	}
}
