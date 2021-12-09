package seeder

import "github.com/nrc-no/core/pkg/api/types"

// General options
var (
	yesNoChoice = []*types.SelectOption{
		{"b22014f7-a7aa-41ec-91c6-6612f874a471", "Yes"},
		{"755947d2-0130-46bc-897c-20da5e521fbd", "No"},
	}
	wgShortSet = []*types.SelectOption{
		{"d3c94761-4974-484e-ab46-e23de5ff13bd", "Moderate Impairment"},
		{"d3c94761-4974-484e-ab46-e23de5ff13bd", "Severe Impairment"},
	}
	globalDisplacementStatuses = []*types.SelectOption{
		{"23e8e352-b9fd-4140-8c4f-aee6984e89c3", "Refugee"},
		{"39999a52-aac3-4530-9bda-089077ec5324", "Internally Displaced Person (DP)"},
		{"93df4b6f-a815-4b42-b043-4facc3aa28c7", "Host Community"},
		{"8d2ae254-24b6-458e-95bd-704b58972a3e", "Other"},
	}
	globalGenders = []*types.SelectOption{
		{"8ef07173-01d7-406b-8d63-2e84d0f3fbc6", "Male"},
		{"23ba6a25-f059-4005-837d-6198f14e9c23", "Female"},
		{"9bd11faa-9d44-4cb2-a710-222dedab7463", "Non-Binary"},
		{"434d6b3e-8c60-4d0c-9a31-08cea585817e", "Other"},
	}
) // Uganda Specific Options
var (
	ugCC = []*types.SelectOption{
		{Name: "ICLA"},
		{Name: "Education"},
		{Name: "LFS"},
		{Name: "Protection"},
		{Name: "WASH"},
		{Name: "Shelter"},
	}
	ugNationality = []*types.SelectOption{
		{"", "Uganda"},
		{"", "South Sudan"},
		{"", "Eritrea"},
		{"", "Somalia"},
		{"", "Democratic Republic of Congo"},
		{"", "Rwanda"},
		{"", "Burundi"},
		{"", "Sudan"},
	}
	ugRelationshipToHH = []*types.SelectOption{
		// TBD
	}
	ugTypeOfIdentification = []*types.SelectOption{
		{"", "Family Attestation (OPM)"},
		{"", "Refugee ID (OPM)"},
		{"", "Asylum Certificate (OPM)"},
		{"", "Resident ID"},
		{"", "National ID"},
		{"", "Passport"},
		{"", "Conventional Travel Document"},
	}
	ugLocationOfIdentification = []*types.SelectOption{
		{"", "Kabusu Access Center"},
		{"", "Nsambya Access Center"},
		{"", "Kisenyi ICLA Center"},
		{"", "Lukuli ICLA Center"},
		{"", "Kawempe ICLA Center"},
		{"", "Ndejje ICLA Center"},
		{"", "Mengo Field Office"},
		{"", "Community (Specify location)"},
		{"", "Home Visit"},
		{"", "Phone"},
		{"", "Other (Specify)"},
	}
	ugSourceOfIdentification = []*types.SelectOption{
		{"", "Walk-in Center"},
		{"", "FFRM Referral"},
		{"", "Internal Referral (Other - Specify)"},
		{"", "ICLA Outreach Team"},
		{"", "External Referral (Community Leader/Contact"},
		{"", "External Referral (INGO/LNGO)"},
		{"", "External Referral (UN Agency)"},
		{"", "External Referral (Government)"},
		{"", "External Referral (Other - Specify)"},
		{"", "Self (Telephone)"},
		{"", "Self (Email)"},
		{"", "Internal Referral (Other NRC Sector – Specify)"},
		{"", "CBP Outreach Team "},
		{"", "Other NRC Outreach Team (Specify)"},
		{"", "Other - Specify"},
	}
	ugAdmin2 = []*types.SelectOption{
		{"", "Abim"},
		{"", "Adjumani"},
		{"", "Alebtong"},
		{"", "Amolatar"},
		{"", "Amudat"},
		{"", "Amuria"},
		{"", "Amuru"},
		{"", "Apac"},
		{"", "Budaka"},
		{"", "Bugiri"},
		{"", "Buikwe"},
		{"", "Bukomansimbi"},
		{"", "Bukwo"},
		{"", "Bulambuli"},
		{"", "Buliisa"},
		{"", "Bundibugyo"},
		{"", "Bushenyi"},
		{"", "Buyende"},
		{"", "Dokolo"},
		{"", "Butambala"},
		{"", "Hoima"},
		{"", "Iganga"},
		{"", "Kaabong"},
		{"", "Kabale"},
		{"", "Kabarole"},
		{"", "Kalangala"},
		{"", "Kaliro"},
		{"", "Kalungu"},
		{"", "Kamuli"},
		{"", "Kanungu"},
		{"", "Kapchorwa"},
		{"", "Katakwi"},
		{"", "Kayunga"},
		{"", "Sheema"},
		{"", "Kitgum"},
		{"", "Koboko"},
		{"", "Kole"},
		{"", "Kotido"},
		{"", "Kisoro"},
		{"", "Kween"},
		{"", "Lamwo"},
		{"", "Lira"},
		{"", "Luuka"},
		{"", "Lyantonde"},
		{"", "Manafwa"},
		{"", "Masaka"},
		{"", "Masindi"},
		{"", "Mayuge"},
		{"", "Mbale"},
		{"", "Mbarara"},
		{"", "Moroto"},
		{"", "Moyo"},
		{"", "Nakapiripirit"},
		{"", "Nakaseke"},
		{"", "Nakasongola"},
		{"", "Namutumba"},
		{"", "Napak"},
		{"", "Nebbi"},
		{"", "Ngora"},
		{"", "Buhweju"},
		{"", "Ntoroko"},
		{"", "Maracha"},
		{"", "Otuke"},
		{"", "Oyam"},
		{"", "Pader"},
		{"", "Rubirizi"},
		{"", "Sironko"},
		{"", "Soroti"},
		{"", "Wakiso"},
		{"", "Yumbe"},
		{"", "Zombo"},
		{"", "Isingiro"},
		{"", "Mitooma"},
		{"", "Kyegegwa"},
		{"", "Ntungamo"},
		{"", "Rukungiri"},
		{"", "Kamwenge"},
		{"", "Ibanda"},
		{"", "Kasese"},
		{"", "Kiruhura"},
		{"", "Kyenjojo"},
		{"", "Mubende"},
		{"", "Gomba"},
		{"", "Kiboga"},
		{"", "Mpigi"},
		{"", "Kyankwanzi"},
		{"", "Kakumiro"},
		{"", "Nwoya"},
		{"", "Kiryandongo"},
		{"", "Serere"},
		{"", "Omoro"},
		{"", "Arua"},
		{"", "Lwengo"},
		{"", "Sembabule"},
		{"", "Rakai"},
		{"", "Mityana"},
		{"", "Luwero"},
		{"", "Mukono"},
		{"", "Kampala"},
		{"", "Buvuma"},
		{"", "Jinja"},
		{"", "Namayingo"},
		{"", "Busia"},
		{"", "Bududa"},
		{"", "Tororo"},
		{"", "Butaleja"},
		{"", "Bukedea"},
		{"", "Kumi"},
		{"", "Pallisa"},
		{"", "Kibuku"},
		{"", "Kaberamaido"},
		{"", "Agago"},
		{"", "Kagadi"},
		{"", "Kibaale"},
		{"", "Gulu"},
		{"", "Rubanda"},
	}

	ugMeansOfContact = []*types.SelectOption{
		{"", "Phone Call"},
		{"", "Text Message"},
		{"", "WhatsApp"},
		{"", "Signal"},
		{"", "Telegram"},
		{"", "Email"},
		{"", "Home Visit"},
		{"", "Other - Specify"},
	}
	ugServicesRequested = []*types.SelectOption{
		{"", "Health care (including medication)"},
		{"", "Legal assistance"},
		{"", "Education"},
		{"", "Mental Health"},
		{"", "Transportation"},
		{"", "Food"},
		{"", "Non-food items (including hygiene items)"},
		{"", "Disability"},
		{"", "MPC"},
		{"", "Shelter/Housing"},
		{"", "Shelter construction/repair"},
		{"", "Youth Livelihoods (e.g. vocational training)"},
		{"", "Small/Medium Business Grants"},
		{"", "Other livelihood activities"},
	}
	ugPerceivedPriorityResponseLevel = []*types.SelectOption{
		{"", "High - Immediate contact within 24 hrs"},
		{"", "Medium - Contact within 72hrs"},
		{"", "Low – Contact within five days "},
	}
	ugSpecificNeed = []*types.SelectOption{
		{"", "Pregnant woman"},
		{"", "Elderly taking care of minors alone"},
		{"", "Single parent"},
		{"", "Chronic illness"},
		{"", "Legal protection needs"},
	}
	ugProtectionRisk = []*types.SelectOption{
		{Name: "Violence"},
		{Name: "Coercion"},
		{Name: "Discrimination"},
		{Name: "Deprivation"},
		{Name: "Other"},
	}
	ugProtectionConcern = []*types.SelectOption{
		{Name: "Physical violence"},
		{Name: "Neglect"},
		{Name: "Family separation"},
		{Name: "Arrest"},
		{Name: "Denial of resources"},
		{Name: "Psychosocial violence"},
	}
	ugProtectionVulnerability = []*types.SelectOption{
		{Name: "Child at risk"},
		{Name: "Elder at risk"},
		{Name: "Single parent"},
		{Name: "Separated child"},
		{Name: "Disability"},
		{Name: "Woman at risk"},
		{Name: "Legal & physical protection"},
		{Name: "Legal & physical protection"},
		{Name: "Medical condition"},
		{Name: "Pregnant/Lactating woman"},
	}
	ugProtectionSpecificNeeds = []*types.SelectOption{
		{Name: "Disability"},
		{Name: "Pregnant woman"},
		{Name: "Elder taking care of minors alone"},
		{Name: "Elder living alone"},
		{Name: "Single parent"},
		{Name: "Chronic illness"},
		{Name: "Legal protection needs"},
		{Name: "Child"},
		{Name: "Other"},
	}
	ugProtectionHomeSituation = []*types.SelectOption{
		{Name: "Lives alone"},
		{Name: "Lives with family"},
		{Name: "Hosted by relative"},
		{Name: "Other"},
	}
	ugProtectionDisability = []*types.SelectOption{
		{Name: "None"},
		{Name: "Moderate physical impairment"},
		{Name: "Severe physical impairment"},
		{Name: "Moderate sensory impairment"},
		{Name: "Severe sensory impairment"},
		{Name: "Moderate mental disability"},
		{Name: "Severe mental disability"},
	}
	ugProtectionHHAbilityToMeetNeeds = []*types.SelectOption{
		{Name: "We can meet all our needs without worry"},
		{Name: "We can meet our needs"},
		{Name: "We can barely meet our needs"},
		{Name: "We are unable to meet our needs"},
		{Name: "We are totally unable to meet our needs"},
	}
	ugProtectionFoodNeedsObstacles = []*types.SelectOption{
		{Name: "Insufficient funds"},
		{Name: "Distance issues"},
		{Name: "Insecurity"},
		{Name: "Social discrimination"},
		{Name: "Insufficient quantity of goods"},
		{Name: "Inadequate quality of goods"},
		{Name: "Insufficient service providers"},
		{Name: "Other"},
	}
	ugProtectionAccommodationNeedsObstacles = []*types.SelectOption{
		{Name: "Insufficient funds"},
		{Name: "Distance issues"},
		{Name: "Insecurity"},
		{Name: "Social discrimination"},
		{Name: "Insufficient quantity of goods"},
		{Name: "Inadequate quality of goods/services"},
		{Name: "Insufficient capabilities/competencies"},
		{Name: "Other"},
	}
	ugProtectionWASHNeedsObstacles = []*types.SelectOption{
		{Name: "Insufficient funds"},
		{Name: "Distance issues"},
		{Name: "Insecurity"},
		{Name: "Social discrimination"},
		{Name: "Insufficient quantity of goods/services"},
		{Name: "Inadequate quality of goods"},
		{Name: "Other"},
	}
	ugProtectionPriorityOpts = []*types.SelectOption{
		{"", "High (follow-up requested in 24 hours)"},
		{"", "Medium (follow-up in 3 days)"},
		{"", "Low (follow-up in 7 days)"}}
	ugICLALegalIssues = []*types.SelectOption{
		{"", "RSD"},
		{"", "ELP"},
		{"", "HLP"},
		{"", "IDP registration"},
		{"", "Other"},
	}
	ugPriority = []*types.SelectOption{
		{"", "High"},
		{"", "Medium"},
		{"", "Low"},
	}
	ugICLADisplacementStatus = []*types.SelectOption{
		{"", "Unregistered asylum seeker"},
		{"", "Registered asylum seeker"},
		{"", "Refugee"},
	}
	ugICLARSDDocuments = []*types.SelectOption{
		{"", "Family Attesttion"},
		{"", "Refugee ID"},
		{"", "Asylum certificate"},
		{"", "Rejection decision"},
		{"", "Other"},
	}
	ugICLASpecificHLPConcern = []*types.SelectOption{
		{"", "Housing"},
		{"", "Land"},
		{"", "Property"},
	}
	ugICLAHomeOwnership = []*types.SelectOption{
		{"", "Own house"},
		{"", "Rent"},
		{"", "Other"},
	}
	ugICLAEvictionDocuments = []*types.SelectOption{
		{"", "Eviction Notice"},
		{"", "Other"},
	}
	ugICLANatureLandTenure = []*types.SelectOption{
		{"", "Joint ownership"},
		{"", "Co-ownership"},
		{"", "Individual ownership"},
		{"", "Other"},
	}
	ugICLANatureTenure = []*types.SelectOption{
		{"", "Mailo"},
		{"", "Lease"},
		{"", "Freehold"},
		{"", "Sustomary"},
	}
	ugICLATypeOfDocumentation = []*types.SelectOption{
		{"", "Legal"},
		{"", "Civil"},
	}
	ugICLATypeOfAgreement = []*types.SelectOption{
		{"", "Oral"},
		{"", "Written"},
	}
	ugICLATypeOfChallenge = []*types.SelectOption{
		{"", "Employment"},
		{"", "Business"},
	}
	ugICLATypeOfActions = []*types.SelectOption{
		{"", "(1) Discussion with supervisor/team leader"},
		{"", "(2) Conducting legal analysis, including the study of judicial practice"},
		{"", "(3) Preparing letters, inquiries to various authorities"},
		{"", "(4) Drafting of other legal documents (such leases or contracts)"},
		{"", "(5) Lodging of a court application"},
		{"", "(6) Attending of court session/hearing"},
		{"", "(7) Review of the decision/appeal"},
		{"", "(8) Negotiation"},
		{"", "(9) Follow up with relevant administrative authority or other entities"},
		{"", "(10) Accompaniment"},
		{"", "(11) Other"},
	}
	ugICLATypesOfServices = []*types.SelectOption{
		{"", "Legal counselling"},
		{"", "Referral"},
		{"", "Relocation"},
		{"", "Livelihood"},
		{"", "Business support"},
	}
	ugICLACaseClosureReason = []*types.SelectOption{
		{"", "All action plan objectives achieved"},
		{"", "Other priorities by beneficiary"},
		{"", "Beneficiary unreachable"},
		{"", "Other"},
	}
	ugAgreedFollowupMeans = []*types.SelectOption{
		{"", "Schedule in-person meeting with beneficiary"},
		{"", "Schedule phone call"},
		{"", "Other"},
	}
	ugEILFSIncomeActivity = []*types.SelectOption{
		{Name: "Employment"},
		{Name: "Business"},
		{Name: "Skill training/Education"},
		{Name: "Farming"},
		{Name: "Bodaboda"},
		{Name: "Causal laborer"},
		{Name: "Other"},
	}
	ugEILFSStartupCapitalSource = []*types.SelectOption{
		{Name: "Donation (NGO, well-wishers, family & friends)"},
		{Name: "Member savings"},
		{Name: "Remittances"},
		{Name: "Loan/Credit"},
		{Name: "Other"},
	}
	ugEILFSProfitUse = []*types.SelectOption{
		{Name: "Reinvested"},
		{Name: "Shared to members"},
		{Name: "Loanable fund"},
		{Name: "Other"},
	}
	ugEILFSTrainingKind = []*types.SelectOption{
		{Name: "Record keeping/Tracking"},
		{Name: "Savings and loans management"},
		{Name: "Financial literacy"},
		{Name: "Business management"},
		{Name: "Agronomic practices"},
		{Name: "Farming as a business"},
		{Name: "Post-harvest handling"},
		{Name: "Skills training"},
		{Name: "Basic numeracy training"},
		{Name: "Lifeskills and entrepreneurship"},
		{Name: "Other"},
	}
	ugEILFSSkillsTraining = []*types.SelectOption{
		{Name: "Tailoring and garment cutting"},
		{Name: "Sewing"},
		{Name: "Arts & crafts"},
		{Name: "Agriculture"},
		{Name: "Briquettes making"},
		{Name: "Reusable packing production"},
		{Name: "Energy stoves"},
		{Name: "Hair dressing"},
		{Name: "Bakery and confectioneries"},
		{Name: "Motor mechanics and driving"},
		{Name: "Other"},
	}
	ugEILFSMarketChallenges = []*types.SelectOption{
		{Name: "Business registration"},
		{Name: "Taxes"},
		{Name: "Distance to markets"},
		{Name: "Product quality"},
		{Name: "Product certification"},
	}
	ugEILFSGroupBusinessArea = []*types.SelectOption{
		{Name: "Savings"},
		{Name: "Credit & loans"},
		{Name: "Arts & crafts"},
		{Name: "Skills training"},
		{Name: "Farming/Agriculture"},
		{Name: "Tailoring"},
		{Name: "Soap making"},
		{Name: "Bakery"},
		{Name: "Salon"},
		{Name: "Other"},
	}
	ugEILFSFinancialServices = []*types.SelectOption{
		{Name: "Credit"},
		{Name: "Savings"},
		{Name: "Insurance"},
		{Name: "Training"},
	}
	ugEILFSFinancialServiceProviders = []*types.SelectOption{
		{Name: "Banks"},
		{Name: "Micro-finance/micro-deposit institutions"},
		{Name: "SACCOs"},
		{Name: "VSLAs"},
		{Name: "Mobile money"},
		{Name: "Other"},
	}
	ugEducationVulnerability = []*types.SelectOption{
		{Name: "(1) Person with disability"},
		{Name: "(2) Living in child-headed household"},
		{Name: "(3) Child mother/father"},
		{Name: "(4) Orphan who's rights are not fulfilled"},
		{Name: "(5) Affected with chronic illness"},
		{Name: "(6) Experiencing forms of abuse"},
		{Name: "(7) Living on the street / Abandoned / Neglected"},
		{Name: "(8) In contact with the law"},
		{Name: "(9) Living in a (poverty-stricken) impoverished household"},
		{Name: "(10) Living with / elderly caregiver / guardian"},
		{Name: "(11) Affected by armed conflict"},
		{Name: "(12) Living in worst forms of labour"},
	}
	ugAction = []*types.SelectOption{
		{Name: "Open"},
		{Name: "Closed"},
		{Name: "Referred"},
		{Name: "Follow-up"},
	}
	ugReferralType = []*types.SelectOption{
		{Name: "External"},
		{Name: "Internal: Shelter/NFI"},
		{Name: "Internal: Livelihood/Food security"},
		{Name: "Internal: Education"},
		{Name: "Internal: WASH"},
		{Name: "Internal: Camp Management/UDOC"},
	}
	ugReferralReason = []*types.SelectOption{
		{Name: ""}, // TBD
	}
	ugReferralMeans = []*types.SelectOption{
		{Name: "(1) Phone"},
		{Name: "(2) E-mail"},
		{Name: "(3) Personal meeting"},
		{Name: "(4) Other"},
	}
	ugReferralPriority = []*types.SelectOption{
		{Name: "High (follow-up requested in 24 hours)"},
		{Name: "Medium (follow-up in 3 days)"},
		{Name: "Low (follow-up in 7 days)"},
	}
)
