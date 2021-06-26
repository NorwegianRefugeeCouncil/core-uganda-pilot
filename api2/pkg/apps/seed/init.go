package seed

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/apps/cms"
	"github.com/nrc-no/core-kafka/pkg/apps/iam"
)

func Init(ctx context.Context, iamClient iam.Interface, cmsClient cms.Interface) error {

	// Create Individuals
	for _, individual := range mockIndividuals {
		if _, err := iamClient.Individuals().Create(ctx, individual); err != nil {
			return err
		}
	}

	// Create case types
	for _, caseType := range []cms.CaseType{
		GenderViolence,
		Childcare,
		HousingRights,
		FinancialAssistInd,
		FinancialAssistHH,
	} {
		_, err := cmsClient.CaseTypes().Create(ctx, &caseType)
		if err != nil {
			return err
		}
	}

	// Create cases
	for _, kase := range []cms.Case{
		DomesticAbuse,
		MonthlyAllowance,
		ChildCare,
	} {
		_, err := cmsClient.Cases().Create(ctx, &kase)
		if err != nil {
			return err
		}
	}

	return nil

}
