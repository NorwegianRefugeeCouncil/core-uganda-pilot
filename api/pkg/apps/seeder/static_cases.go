package seeder

var (
	BoDiddleySituationAnalysisData = map[string][]string{
		"safeDignifiedLife":  {"Yes, I live a safe and dignified life and I am reasonably happy with my achievements and quality of life."},
		"challengesBarriers": {"Yes, I live a safe and dignified life and I am reasonably happy with my achievements and quality of life."},
		"solutions":          {"A qualified interpreter, who knows the legal context could help us to agree on contractual matters."},
		"workTogether":       {"NRC could provide a translator and a legal representative to ease contract negotiations"},
	}

	BoDiddleyResponseData = map[string][]string{
		"servicesStartingPoint": {"ICLA"},
		"commentStartingPoint":  {"The individual has requested ICLA as a starting point, we should create a referral"},
		"otherServices":         {"Protection"},
		"commentOtherServices":  {"The individual has requested additional Protection services, we should create a referral"},
		"perceivedPriority":     {"High"},
	}
	MaryPoppinsSituationAnalysisData = map[string][]string{
		"safeDignifiedLife":  {"Yes, I live a safe and dignified life and I am reasonably happy with my achievements and quality of life."},
		"challengesBarriers": {"Some of the barriers I face are communication gaps between myself and refugee tenants. We are attempting to deal with these challenges by using google translate."},
		"solutions":          {"A qualified interpreter, who knows the legal context could help us to agree on contractual matters."},
		"workTogether":       {"NRC could provide a translator and a legal representative to ease contract negotiations"},
	}
	MaryPoppinsResponseData = map[string][]string{
		"servicesStartingPoint": {"S&S"},
		"commentStartingPoint":  {"The individual has requested S&S as a starting point, we should create a referral"},
		"otherServices":         {"Protection"},
		"commentOtherServices":  {"The individual has requested additional Protection services, we should create a referral"},
		"perceivedPriority":     {"High"},
	}
	JohnDoeSituationAnalysisData = map[string][]string{
		"safeDignifiedLife":  {"Yes, I live a safe and dignified life and I am reasonably happy with my achievements and quality of life."},
		"challengesBarriers": {"Some of the barriers I face are communication gaps between myself and refugee tenants. We are attempting to deal with these challenges by using google translate."},
		"solutions":          {"A qualified interpreter, who knows the legal context could help us to agree on contractual matters."},
		"workTogether":       {"NRC could provide a translator and a legal representative to ease contract negotiations"},
	}
	JohnDoeResponseData = map[string][]string{
		"servicesStartingPoint": {"LFS"},
		"commentStartingPoint":  {"The individual has requested LFS as a starting point, we should create a referral"},
		"otherServices":         {"WASH"},
		"commentOtherServices":  {"The individual has requested additional WASH services, we should create a referral"},
		"perceivedPriority":     {"High"},
	}

	BoDiddleySituationAnalysis  = kase("dba43642-8093-4685-a197-f8848d4cbaaa", Colette.ID, BoDiddley.ID, KampalaProtectionTeam.ID, UGSituationalAnalysisCaseType, true, true, BoDiddleySituationAnalysisData)
	BoDiddleyIndividualResponse = kase("3ea8c121-bdf0-46a0-86a8-698dc4abc872", Colette.ID, BoDiddley.ID, KampalaProtectionTeam.ID, UGIndividualResponseCaseType, true, true, BoDiddleyResponseData)

	MaryPoppinsSituationAnalysis  = kase("4f7708ed-240a-423f-9bd1-839542e65833", Colette.ID, MaryPoppins.ID, KampalaProtectionTeam.ID, UGSituationalAnalysisCaseType, true, true, MaryPoppinsSituationAnalysisData)
	MaryPoppinsIndividualResponse = kase("45b4a637-c610-4ab9-afe6-4e958c36a96f", Colette.ID, MaryPoppins.ID, KampalaProtectionTeam.ID, UGIndividualResponseCaseType, true, true, MaryPoppinsResponseData)

	JohnDoesSituationAnalysis = kase("43140381-8166-4fb3-9ac5-339082920ade", Colette.ID, JohnDoe.ID, KampalaProtectionTeam.ID, UGSituationalAnalysisCaseType, true, true, JohnDoeSituationAnalysisData)
	JohnDoeIndividualResponse = kase("65e02e79-1676-4745-9890-582e3d67d13f", Colette.ID, JohnDoe.ID, KampalaProtectionTeam.ID, UGIndividualResponseCaseType, true, true, JohnDoeResponseData)
)
