// Package engine provides the core database engine
//
// This is a *headless* database engine that can be used to store and retrieve
// records. It is headless because it does not communicate with the database
// directly. Instead, it generates SQL statements that are to be executed.
//
// The engine is designed to be used in standalone mode or as a part of a
// peer-to-peer network.
//
// # Introduction
// The engine is designed to store and retrieve records in a table.
// A table is a collection of records.
// A record is a collection of fields.
// A record has a unique id.
// A field is a named value.
// A value is a string, a number, or a boolean.
//
// The engine is designed to be used in standalone mode or as a part of a
// peer-to-peer network.
//
// For each table, the engine internally creates two distinct tables:
// * <table>
// * <table>_history
//
// The <table> table is used to store the reconciled state of the records.
// The <table>_history table is used to store the history of the records.
//
// When a new version of a record is created, the engine will append this
// version inside the <table>_history table. The <table> table will be
// updated to reflect the new version.
//
// When multiple versions of a record are in a conflict, the engine will
// choose the version with the highest revision number. This is the version
// of the record that will be reflected in the <table> table.
//
// The revision of a record is a string in the format <num>-<hash>.
// The <num> is an incrementing number.
// The <hash> is a hash of the record.
//
// Two records might have the same number but different hashes.
// In this case, the record with the alphabetically higher hash is the
// winner.
//
// The engine also stores the previous revision of a record.
// Which allows to build the history of a record, where the history
// is a branching tree.
package engine
