import NewRelationshiptypePage from '../pages/newRelationshiptype.page';

describe('CaseType Page', function () {
  describe('Create', () => {
    it('should navigate to New CaseType Form Page and save new attribute', () => {
      const newRelationshiptypePage = new NewRelationshiptypePage();
      newRelationshiptypePage
        .visitPage()
        .typeName('Test')
        .checkIsDirectional()
        .typeFirstPartyRole ('Is head of household of')
        .typeSecondPartyRole("Has for head of household")
        .selectFristPartyType("Individual")
        .selectSecondPartyType("Household")
        .save();
    });
  });
});
