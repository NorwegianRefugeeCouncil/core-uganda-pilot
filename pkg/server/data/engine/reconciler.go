package engine

import (
	"context"
	"errors"

	"github.com/nrc-no/core/pkg/server/data/api"
)

const checkpointTableName = "local_reconciler_changes"
const checkpointKey = "checkpoint"

// reconciler is a type that reconciles a source and a target database
type reconciler struct{}

// initTable creates the necessary table for the reconciler to work
// the table contains the reconciled revision of the source database
// so that we don't pull tons of data from the source database every
// time we run the reconciler.
//
// The database change stream maintains a sequence number for each
// operation to the database. This incremental sequence number is used
// to determine which operations have already been reconciled.
//
// It is basically a map <source> -> <reconciled sequence>
// that exists in the destination database.
func initTable(ctx context.Context, destination api.Engine) error {
	_, err := destination.CreateTable(ctx, api.Table{
		Name: checkpointTableName,
		Columns: []api.Column{
			{
				Name: checkpointKey,
				Type: "integer",
				Constraints: []api.ColumnConstraint{
					{
						NotNull: &api.NotNullColumnConstraint{},
					},
				},
			},
		},
	})
	if api.IsError(err, api.ErrCodeTableAlreadyExists) {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

// Reconcile reconciles the source and target database
//
// It must work in a peer to peer environment, where there are no authority, but
// where both the clients and the servers are both sources of truth.
//
// Usually, the target database would be the locally running database, and the
// source database would be the remote database. So each database could technically
// run a reconciler alongside itself, and pull the data from the other databases.
//
// The logic of the reconciler is as follows:
//
// Given a source and a target database, we want to reconcile
// the source database with the target database. The goal is to pull the changes
// from the source database that are not already in the target database.
//
// For each source database, the reconciler maintains a sequence number
// corresponding to the last change that has been reconciled from that source.
// That sequence number is maintained in the destination database in a
// 'local_reconciler_changes' table. It is simply a pair of the source database ID
// and the last seen sequence number for that source.
//
// E.g.
//  source seq
//  foo    101
//  bar    200
//  baz    0
//
// The reconciler will query the source database to get the changes that
// happened since that known sequence number. And for each change, the
// reconciler will apply the change to the target database.
//
// If the destination already has the change, the reconciler will skip it.
// If the destination does not have the change, the reconciler will apply the change.
//
// Once the changes are applied, the reconciler will update the sequence number
// in the 'local_reconciler_changes' table to reflect the last change that was reconciled
// from that source database.
//
// E.g.
//  source seq
//  foo    150
//  bar    210
//  baz    5
//
// This mechanism allows multiple databases to be reconciled in between each other.
// The source and target databases can be accessed locally, or via other protocols
// such as HTTP, GRPC, etc. As long as the protocol implements the Engine interface.
//
// All the reconciler does is to synchronize the source history tables into the destination's history tables.
// The engine in the destination will take care of electing the winning records in the event
// of a conflict. This is also a deterministic process, so we can be sure that the
// results will be the same regardless of the order in which the changes are applied,
// regardless of the database that is the source and the database that is the target.
func (r *reconciler) Reconcile(ctx context.Context, source api.ReadInterface, destination api.Engine) error {

	// TODO assign unique id to destination
	var sourceId = "dest1234"
	var checkpoint int64 = -1
	var checkpointRec *api.Record
	var err error

	// create the checkpoint table if it doesn't exist
	if err = initTable(ctx, destination); err != nil {
		return err
	}

	// get the last checkpoint from the destination database for the source database
	foundCheckpointRec, err := destination.GetRecord(ctx, api.GetRecordRequest{
		TableName: checkpointTableName,
		RecordID:  sourceId,
	})
	if err != nil {
		if !errors.Is(err, api.ErrRecordNotFound) {
			return err
		}
	} else {
		checkpointRec = &foundCheckpointRec
	}

	if checkpointRec != nil {
		// already has a checkpoint for the source database
		// retrieve the checkpoint
		var (
			fieldValue interface{}
			ok         bool
		)
		if fieldValue, err = checkpointRec.GetFieldValue(checkpointKey); err != nil {
			return err
		}
		if checkpoint, ok = fieldValue.(int64); !ok {
			return errors.New("checkpoint is not an int64")
		}
	} else {
		checkpointRec = &api.Record{
			Table: checkpointTableName,
			ID:    sourceId,
			Attributes: map[string]api.Value{
				checkpointKey: api.NewIntValue(checkpoint, true),
			},
		}
	}

	changes, err := source.GetChanges(ctx, api.GetChangesRequest{Since: checkpoint})
	if err != nil {
		return err
	}
	// todo: Wrap this in a transaction somehow
	// probably have to create a WithTransaction(fn func(e Engine) error) method
	// But that would only work locally, grpc or through websockets, not through HTTP
	// Perhaps it's not so bad that reconciliation is partial, but it's still not ideal
	// todo: batch this for sure
	for _, change := range changes.Items {
		// for each change in the source change stream

		// skip if this checkpoint has already been reconciled
		if change.Sequence == checkpoint {
			continue
		}

		// check if the record revision already exist in the destination
		// if so, we don't need to do anything
		_, err = destination.GetRecord(ctx, api.GetRecordRequest{
			TableName: change.TableName,
			RecordID:  change.RecordID,
			Revision:  change.RecordRevision,
		})
		if !api.IsError(err, api.ErrCodeRecordNotFound) {
			return err
		}

		// get the revision from the source
		sourceRevisionRec, err := source.GetRecord(ctx, api.GetRecordRequest{
			TableName: change.TableName,
			RecordID:  change.RecordID,
			Revision:  change.RecordRevision,
		})
		if err != nil {
			return err
		}

		// insert the revision into the destination
		if _, err := destination.PutRecord(ctx, api.PutRecordRequest{
			Record:        sourceRevisionRec,
			IsReplication: true,
		}); err != nil {
			return err
		}
		checkpointRec.SetFieldValue(checkpointKey, api.NewIntValue(checkpoint, true))
	}

	if _, err := destination.PutRecord(ctx, api.PutRecordRequest{
		Record: *checkpointRec,
	}); err != nil {
		return err
	}

	return nil
}
