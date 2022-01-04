import { FieldDefinition, SelectOption } from 'core-api-client';
import React, { FC } from 'react';

import { RecordPickerContainer } from '../../components/RecordPicker';
import { FormValue } from '../../reducers/recorder';

export type FieldEditorProps = {
  field: FieldDefinition;
  value: string | string[] | null;
  setValue: (value: string | string[] | null) => void;
  addSubRecord: () => void;
  selectSubRecord: (subRecordId: string) => void;
  subRecords: FormValue[] | undefined;
};

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

export const ReferenceFieldEditor: FC<FieldEditorProps> = (props) => {
  const { field, value, setValue } = props;
  if (Array.isArray(value)) {
    return <></>;
  }
  return (
    <div className="form-group mb-2">
      {mapFieldLabel(field)}
      <RecordPickerContainer
        formId={field.fieldType.reference?.formId}
        recordId={value}
        setRecordId={setValue}
      />
      {mapFieldDescription(field)}
    </div>
  );
};

export const TextFieldEditor: FC<FieldEditorProps> = (props) => {
  const { field, value, setValue } = props;
  return (
    <div className="form-group mb-2">
      {mapFieldLabel(field)}
      <input
        className="form-control bg-dark text-light border-secondary"
        type="text"
        id={field.id}
        value={value || ''}
        onChange={(event) => setValue(event.target.value)}
      />
      {mapFieldDescription(field)}
    </div>
  );
};

export const MultilineTextFieldEditor: FC<FieldEditorProps> = (props) => {
  const { field, value, setValue } = props;
  return (
    <div className="form-group mb-2">
      {mapFieldLabel(field)}
      <textarea
        className="form-control bg-dark text-light border-secondary"
        id={field.id}
        value={value || ''}
        onChange={(event) => setValue(event.target.value)}
      />
      {mapFieldDescription(field)}
    </div>
  );
};

export const DateFieldEditor: FC<FieldEditorProps> = (props) => {
  const { field, value, setValue } = props;
  return (
    <div className="form-group mb-2">
      {mapFieldLabel(field)}
      <input
        className="form-control bg-dark text-light border-secondary"
        type="date"
        id={field.id}
        value={value || ''}
        onChange={(event) => setValue(event.target.value)}
      />
      {mapFieldDescription(field)}
    </div>
  );
};

export const MonthFieldEditor: FC<FieldEditorProps> = (props) => {
  const { field, value, setValue } = props;
  const expectedLength = 7;

  function isValid(s: string) {
    const valid = /^(?:19|20|21)\d{2}-[01]\d$/;
    const m = +s.slice(5);
    return valid.test(s) && m > 0 && m <= 12;
  }

  if (Array.isArray(value)) {
    return <></>;
  }
  return (
    <div className="form-group mb-2">
      {mapFieldLabel(field)}
      <input
        className="form-control bg-dark text-light border-secondary"
        type="month"
        maxLength={expectedLength}
        id={field.id}
        value={value || ''}
        name={field.name}
        pattern="[0-9]{4}-[0-9]{2}"
        placeholder="YYYY-MM"
        onChange={(event) => {
          const v = event.target.value;
          if (!isValid(v)) return;
          setValue(v);
        }}
      />
      {mapFieldDescription(field)}
    </div>
  );
};

export const WeekFieldEditor: FC<FieldEditorProps> = (props) => {
  const { field, value, setValue } = props;

  function isValidWeek(weekString: string) {
    const weekRegex = /^(?:19|20|21)\d{2}-W[0-5]\d$/;
    return weekRegex.test(weekString) && +weekString.slice(6) <= 52;
  }

  function onChangeHandler(event: React.ChangeEvent<HTMLInputElement>) {
    if (!isValidWeek(event.target.value)) return;
    setValue(event.target.value);
  }

  return (
    <div className="form-group mb-2">
      {mapFieldLabel(field)}
      <input
        className="form-control bg-dark text-light border-secondary"
        type="week"
        name={field.name}
        maxLength={8}
        placeholder="2021-W52"
        id={field.id}
        value={value || ''}
        onChange={onChangeHandler}
      />
      {mapFieldDescription(field)}
    </div>
  );
};

export const QuantityFieldEditor: FC<FieldEditorProps> = (props) => {
  const { field, value, setValue } = props;
  return (
    <div className="form-group mb-2">
      {mapFieldLabel(field)}
      <input
        className="form-control bg-dark text-light border-secondary"
        type="number"
        id={field.id}
        value={value || ''}
        onChange={(event) => setValue(event.target.value)}
      />
      {mapFieldDescription(field)}
    </div>
  );
};

export const SingleSelectFieldEditor: FC<FieldEditorProps> = (props) => {
  const { field, value, setValue } = props;
  return (
    <div className="form-group mb-2">
      {mapFieldLabel(field)}
      <select
        className="form-control bg-dark text-light border-secondary"
        id={field.id}
        value={value || ''}
        onChange={(event) => setValue(event.target.value)}
      >
        {mapSelectOptions(
          field.required,
          field.key,
          field.fieldType?.singleSelect?.options,
        )}
      </select>
      {mapFieldDescription(field)}
    </div>
  );
};

export const MultiSelectFieldEditor: FC<FieldEditorProps> = (props) => {
  const { field, value, setValue } = props;
  return (
    <div className="form-group mb-2">
      {mapFieldLabel(field)}
      <select
        className="form-control bg-dark text-light border-secondary"
        id={field.id}
        value={value || []}
        multiple
        onChange={(event) => {
          const { options } = event.target;
          const selected = Object.entries(options).filter((o) => o[1].selected);
          setValue(selected.map((s) => s[1].value));
        }}
      >
        {mapSelectOptions(
          field.required,
          field.key,
          field.fieldType?.multiSelect?.options,
        )}
      </select>
      {mapFieldDescription(field)}
    </div>
  );
};

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

function mapSubRecords(records: FormValue[], select: (id: string) => void) {
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

export const SubFormFieldEditor: FC<FieldEditorProps> = (props) => {
  const { field, addSubRecord, selectSubRecord, subRecords } = props;
  return (
    <div className="mb-2">
      <div className="bg-primary border-2" />
      <span className="form-label opacity-75">{field.name}</span>
      {subRecords ? mapSubRecords(subRecords, selectSubRecord) : <></>}
      <button
        type="button"
        onClick={addSubRecord}
        className="btn btn-sm btn-outline-primary w-100"
      >
        Add record in {field.name}
      </button>
      {mapFieldDescription(field)}
    </div>
  );
};

export const FieldEditor: FC<FieldEditorProps> = (props) => {
  const {
    field: { fieldType },
  } = props;
  if (fieldType.text) {
    return <TextFieldEditor {...props} />;
  }
  if (fieldType.week) {
    return <WeekFieldEditor {...props} />;
  }
  if (fieldType.subForm) {
    return <SubFormFieldEditor {...props} />;
  }
  if (fieldType.reference) {
    return <ReferenceFieldEditor {...props} />;
  }
  if (fieldType.multilineText) {
    return <MultilineTextFieldEditor {...props} />;
  }
  if (fieldType.date) {
    return <DateFieldEditor {...props} />;
  }
  if (fieldType.month) {
    return <MonthFieldEditor {...props} />;
  }
  if (fieldType.quantity) {
    return <QuantityFieldEditor {...props} />;
  }
  if (fieldType.singleSelect) {
    return <SingleSelectFieldEditor {...props} />;
  }
  if (fieldType.multiSelect) {
    return <MultiSelectFieldEditor {...props} />;
  }
  return <></>;
};
