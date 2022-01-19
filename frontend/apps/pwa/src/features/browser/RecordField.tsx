import { FieldDefinition, Record } from 'core-api-client';
import { FC } from 'react';
import { Link } from 'react-router-dom';
import format from 'date-fns/format';

type RecordFieldProps = {
  field: FieldDefinition;
  value: string;
  subRecords: Record[] | undefined;
};

export const RecordField: FC<RecordFieldProps> = ({
  field,
  value,
  subRecords,
}) => {
  const renderField = (f: FieldDefinition, v: any) => {
    if (!v) return <div className="text-muted">-</div>;

    if (f.fieldType.reference) {
      return (
        <div>
          <Link to={`/browse/records/${v}`}>View</Link>
        </div>
      );
    }
    if (f.fieldType.month) {
      return <div className="fw-bold">{format(new Date(v), 'yyyy-MM')}</div>;
    }
    if (f.fieldType.week) {
      return <div className="fw-bold">{format(new Date(v), "yyyy-'W'ww")}</div>;
    }
    return (
      <div className="fw-bold" style={{ whiteSpace: 'pre-wrap' }}>
        {v}
      </div>
    );
  };

  return (
    <div className="form-group mb-3">
      <span className="form-label">{field.name}</span>
      {renderField(field, value)}
      {subRecords?.map((r) => (
        <Link key={r.id} to={`/browse/records/${r.id}`}>
          Sub Record
        </Link>
      ))}
    </div>
  );
};
