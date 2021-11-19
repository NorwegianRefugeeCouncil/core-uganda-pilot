package test

import (
	"strings"

	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/server/public/handlers/database"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestDatabaseValidation() {
	// test valid name length
	tooShort := "aa"
	tooLong := strings.Repeat("a", 33)
	justRight := "aaaa"

	assert.False(s.T(), database.ValidDBNameLength(tooShort))
	assert.False(s.T(), database.ValidDBNameLength(tooLong))
	assert.True(s.T(), database.ValidDBNameLength(justRight))

	// test no leading or trailing whitespace
	withSpacesAtFront := " \n\ttest"
	withSpacesAtEnd := "test\n\t "
	withSpacesAtBoth := " \n\ttest\n\t "
	withoutSpaces := "test"

	assert.False(s.T(), database.DbNameHasNoLeadingOrTrailingWhitespace(withSpacesAtFront))
	assert.False(s.T(), database.DbNameHasNoLeadingOrTrailingWhitespace(withSpacesAtEnd))
	assert.False(s.T(), database.DbNameHasNoLeadingOrTrailingWhitespace(withSpacesAtBoth))
	assert.True(s.T(), database.DbNameHasNoLeadingOrTrailingWhitespace(withoutSpaces))

	// test integrated valdation
	dbStructWithEmojis := types.Database{
		Name: "🕵️🕵️🕵️🕵️",
	}

	assert.False(s.T(), database.ValidateDBStruct(&dbStructWithEmojis))
}
