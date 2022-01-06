import { FC } from 'react';
import { Link } from 'react-router-dom';

import {
  useFormOrSubForm,
  useOwnerRecord,
  useRecordFromPath,
  useSubRecords,
} from '../../app/hooks';

import { RecordFieldList } from './RecordFieldList';

export const RecordBrowser: FC = () => {
  const record = useRecordFromPath('recordId');
  const form = useFormOrSubForm(record?.formId);
  const subRecords = useSubRecords(record?.id);
  const ownerRecord = useOwnerRecord(record?.id);

  if (!record) {
    return <div>Record not found</div>;
  }
  if (!form) {
    return <div>Form not found</div>;
  }
  return (
    <div className="flex-grow-1 bg-light py-3">
      <div className="container">
        <div className="row">
          <div className="col">
            <div className="card shadow">
              <div className="card-body">
                {ownerRecord ? (
                  <div className="mb-2">
                    <Link to={`/browse/records/${ownerRecord.id}`}>
                      Back to Parent Record
                    </Link>
                  </div>
                ) : (
                  <></>
                )}
                {form?.fields.map((f) => (
                  <RecordFieldList
                    key={f.id}
                    record={record}
                    field={f}
                    subRecords={subRecords?.byFieldId[f.id]}
                  />
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
