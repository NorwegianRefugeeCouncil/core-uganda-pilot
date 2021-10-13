package iam

import (
	"context"
	"github.com/nrc-no/core/pkg/generic/pagination"
	"github.com/nrc-no/core/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IndividualStore struct {
	getCollection utils.MongoCollectionFn
}

func NewIndividualStore(mongoClientFn utils.MongoClientFn, database string) *IndividualStore {
	store := &IndividualStore{
		getCollection: utils.GetCollectionFn(database, "parties", mongoClientFn),
	}
	return store
}

func (s *IndividualStore) create(ctx context.Context, individual *Individual) (*Individual, error) {
	individual.AddPartyType(IndividualPartyType.ID)
	collection, done, err := s.getCollection(ctx)
	if err != nil {
		return nil, err
	}
	defer done()

	_, err = collection.InsertOne(ctx, individual)
	if err != nil {
		return nil, err
	}
	return individual, nil
}

func (s *IndividualStore) get(ctx context.Context, id string) (*Individual, error) {
	collection, done, err := s.getCollection(ctx)
	if err != nil {
		return nil, err
	}
	defer done()

	findResult := collection.FindOne(ctx, bson.M{"id": id})
	if findResult.Err() != nil {
		return nil, findResult.Err()
	}
	var individual *Individual
	if err := findResult.Decode(&individual); err != nil {
		return nil, err
	}
	return individual, nil
}

func (s *IndividualStore) upsert(ctx context.Context, individual *Individual) (*Individual, error) {
	collection, done, err := s.getCollection(ctx)
	if err != nil {
		return nil, err
	}
	defer done()

	individual.AddPartyType(IndividualPartyType.ID)
	_, err = collection.UpdateOne(ctx, bson.M{
		"id": individual.ID,
	}, bson.M{
		"$set": bson.M{
			"attributes":   individual.Attributes,
			"partyTypeIds": individual.PartyTypeIDs,
		},
	}, options.Update().SetUpsert(true))

	if err != nil {
		return nil, err
	}
	return individual, nil
}

func (s *IndividualStore) list(ctx context.Context, listOptions IndividualListOptions) (*IndividualList, error) {
	includesIndividualPartyType := false

	for _, partyTypeID := range listOptions.PartyTypeIDs {
		if partyTypeID == IndividualPartyType.ID {
			includesIndividualPartyType = true
		}
	}

	if !includesIndividualPartyType {
		listOptions.PartyTypeIDs = append(listOptions.PartyTypeIDs, IndividualPartyType.ID)
	}

	filter := bson.M{
		"partyTypeIds": bson.M{
			"$all": listOptions.PartyTypeIDs,
		},
	}

	for key, value := range listOptions.Attributes {
		filter["attributes."+key] = value
	}

	if len(listOptions.SearchParam) != 0 {
		filter["$text"] = bson.M{"$search": listOptions.SearchParam}
	}

	collection, done, err := s.getCollection(ctx)
	if err != nil {
		return nil, err
	}
	defer done()

	maxPerPage := pagination.GetMaxPerPage(listOptions.PerPage)
	currentPage := pagination.GetCurrentPage(listOptions.Page)
	totalCount, _ := collection.CountDocuments(ctx, filter)

	cursor, err := collection.Find(ctx, filter, getFindOptions(currentPage, maxPerPage, listOptions))
	if err != nil {
		return nil, err
	}
	var items []*Individual
	for {
		if !cursor.Next(ctx) {
			break
		}
		var b Individual
		if err := cursor.Decode(&b); err != nil {
			return nil, err
		}
		items = append(items, &b)
	}
	if cursor.Err() != nil {
		return nil, cursor.Err()
	}
	if items == nil {
		items = []*Individual{}
	}
	return &IndividualList{
		Items:    items,
		Metadata: pagination.GetPaginationMetaData(int(totalCount), currentPage, maxPerPage, listOptions.Sort, listOptions.SearchParam),
	}, nil
}

func getFindOptions(currentPage int, maxPerPage int, listOptions IndividualListOptions) *options.FindOptions {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "attributes." + DisplayNameAttribute.ID, Value: pagination.GetSortOptionType(listOptions.Sort)}})
	findOptions.SetSkip(int64((currentPage - 1) * maxPerPage))
	findOptions.SetLimit(int64(maxPerPage))
	return findOptions
}
