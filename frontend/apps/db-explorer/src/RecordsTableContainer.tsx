import {useCallback, useEffect, useState} from 'react';
import {Client, Record, Table, Value, isStringValue} from './client';

type RecordsTableProps = {
  records: Record[];
  table: Table;
  selectedRecordId?: string;
  onSelectRecord: (record: Record) => void;
  recordJSON: string;
  onChangeRecordJSON: (recordJSON: string) => void;
  onPutRecord: () => void;
};

function renderValue(value: Value) {
  if (isStringValue(value)) {
    return <span>{value.string}</span>;
  }
  return <span>{JSON.stringify(value)}</span>;
}

function RecordsTable(props: RecordsTableProps) {
  const {records, recordJSON, onChangeRecordJSON, onSelectRecord, table: {columns}} = props;
  return (
    <div className={'container'}>
      <div className={'row'}>
        <div className={'col'}>
          <table className={'table table-bordered table-hover'}>
            <thead>
            <tr>
              {columns.map((column) => (
                <th key={column.name}>{column.name}</th>
              ))}
            </tr>
            </thead>
            <tbody>
            {records.map(record => (
              <tr
                style={{"cursor": "pointer"}}
                key={record.id} onClick={() => onSelectRecord(record)}>
                {columns.map((column) => {
                  if (column.name === '_id') {
                    return <td key={column.name}>{record.id}</td>;
                  }
                  if (column.name === '_rev') {
                    return <td key={column.name}>{record.revision}</td>;
                  }
                  return <td key={column.name}>{renderValue(record.attributes[column.name])}</td>;
                })}
              </tr>
            ))}
            </tbody>
          </table>
          <textarea
            rows={10}
            className={'bg-dark font-monospace text-light w-100 fw-bold'}
            onChange={(e) => onChangeRecordJSON(e.target.value)}
            value={recordJSON}
          >
          </textarea>
          <button
            onClick={props.onPutRecord}
            className={'btn btn-primary'}>
            Put Record
          </button>
        </div>
      </div>
    </div>
  );
}

type RecordsTableContainerProps = {
  tableName: string,
  client: Client
}

function RecordsTableContainer(props: RecordsTableContainerProps) {

  const {tableName, client} = props;

  const [records, setRecords] = useState<Record[]>([]);
  const [recordsLoading, setRecordsRecordsLoading] = useState(true);

  const [table, setTable] = useState<Table>();
  const [tableLoading, setTableLoading] = useState(true);

  const [recordJSON, setRecordJSON] = useState<string>('');
  const onChangeRecordJSON = useCallback((jsonValue: string) => {
    setRecordJSON(jsonValue);
  }, []);

  const [selectedRecordId, setSelectedRecordId] = useState<string>();
  const onSelectRecord = useCallback((record: Record) => {
    setSelectedRecordId(record.id);
    setRecordJSON(JSON.stringify(record, null, 2));
  }, []);

  const fetchRecords = useCallback(async () => {
    setRecordsRecordsLoading(true);
    const fetchedRecords = await client.getRecords({table: tableName});
    setRecords(fetchedRecords.items);
    setRecordsRecordsLoading(false);
  }, [tableName, client]);

  const fetchTable = useCallback(async () => {
    setTableLoading(true);
    const fetchedTable = await client.getTable({tableName});
    setTable(fetchedTable);
    setTableLoading(false);
    const mockRecord: any = {};
    mockRecord['id'] = '';
    mockRecord['attributes'] = {};
    fetchedTable.columns.forEach(column => {
      if (column.name === '_rev' || column.name === '_id') {
        return;
      }
      mockRecord.attributes[column.name] = {'string': ''};
    });
    setRecordJSON(JSON.stringify(mockRecord, null, 2));
  }, [tableName, client]);

  const putRecord = useCallback(async () => {
    const record = JSON.parse(recordJSON) as Record;
    record.table = tableName;
    delete (record as any)._id;
    delete (record as any)._rev;
    const result = await client.putRecord({record: record});
    fetchRecords();
    setRecordJSON(JSON.stringify(result, null, 2));
  }, [tableName, recordJSON, client, fetchRecords]);

  useEffect(() => {
    fetchRecords();
    fetchTable();
  }, [fetchRecords]);

  if (tableLoading || recordsLoading) {
    return <div>Loading...</div>;
  }

  if (!table) {
    return <div>Table not found</div>;
  }

  return (
    <RecordsTable
      selectedRecordId={selectedRecordId}
      records={records}
      onSelectRecord={onSelectRecord}
      table={table}
      onChangeRecordJSON={onChangeRecordJSON}
      recordJSON={recordJSON}
      onPutRecord={putRecord}
    />
  );
}

export default RecordsTableContainer;
