package seed

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/apps/cms"
	"github.com/nrc-no/core-kafka/pkg/apps/iam"
	"github.com/sirupsen/logrus"
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

	// Create people
	logger := logrus.WithContext(ctx).WithField("logger", "staffmock.Init")
	for _, person := range people {
		logger.Infof("initializing mock party %#v", person)
		if _, err := iamClient.Parties().Create(ctx, person); err != nil {
			logger.WithError(err).Errorf("failed to create mock party: %#v", person)
			return err
		}
	}

	// Create staff
	for _, s := range staffs {
		logger.Infof("initializing mock staff %#v", s)
		if _, err := iamClient.Staff().Create(ctx, s); err != nil {
			logger.WithError(err).Errorf("failed to initialize mock staff: %#v", s)
			return err
		}
	}

	// Create memberships
	for _, membership := range mbships {
		logger.Infof("initializing mock membership %#v", membership)
		if _, err := iamClient.Memberships().Create(ctx, membership); err != nil {
			logger.WithError(err).Errorf("failed to initialize mock membership: %#v", membership)
			return err
		}
	}

	return nil

}
