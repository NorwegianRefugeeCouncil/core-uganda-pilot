import React from 'react';
import { FieldDefinition, SelectOption } from 'core-api-client';

import { FormValue } from '../../reducers/recorder';

function subRecord(record: FormValue, select: () => void) {
  return (
    <a
      href="/#"
      key={record.id}
      onClick={(e) => {
        e.preventDefault();
        select();
      }}
      className="list-group-item list-group-item-action bg-dark border-secondary text-secondary"
    >
      View Record
    </a>
  );
}

export function mapSubRecords(
  records: FormValue[],
  select: (id: string) => void,
) {
  return (
    <div className="list-group bg-dark mb-3">
      {records.map((r) =>
        subRecord(r, () => {
          select(r.id);
        }),
      )}
    </div>
  );
}
export const mapFieldDescription = (fd: FieldDefinition) => {
  if (fd.description) {
    return <small className="text-muted">{fd.description}</small>;
  }
  return <></>;
};
export const mapFieldLabel = (fd: FieldDefinition) => {
  return (
    <label className="form-label opacity-75" htmlFor={fd.id}>
      {fd.name}
    </label>
  );
};

export const mapSelectOptions = (
  required: boolean,
  key: boolean,
  options?: SelectOption[],
) => {
  if (!options) {
    return <></>;
  }
  return (
    <>
      <option aria-label="no value" disabled={required || key} value="" />
      {options.map((o) => (
        <option key={o.id} value={o.id}>
          {o.name}
        </option>
      ))}
    </>
  );
};
