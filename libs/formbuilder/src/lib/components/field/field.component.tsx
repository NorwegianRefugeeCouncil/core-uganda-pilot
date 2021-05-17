import React, {
  useState,
  Fragment,
  FunctionComponent,
  FormEvent,
  Dispatch,
  SetStateAction,
} from 'react';
import {
  FormLabel,
  FormCheck,
  FormCheckLabel,
  FormCheckInput,
  FormInput,
} from '@nrc.no/ui-toolkit';

export type Translations = {
  [locale: string]: string;
};

export interface SelectOption {
  key: string;
  value: string;
}

export interface RadioOption {
  key: string;
}

export interface FieldOptions {
  name: Translations;
  description: Translations;
  tooltip: Translations;
  min?: number | string;
  max?: number | string;
  options?: SelectOption[] | RadioOption[];
  maxLength?: number;
  regex?: string;
  required?: boolean;
  disabled?: boolean;
  hidden?: boolean;
  value: any;
  default?: any;
}

export enum FieldType {
  string = 'string',
  integer = 'integer',
  float = 'float',
  checkbox = 'checkbox',
  radio = 'radio',
  select = 'select',
  multiselect = 'multiselect',
}

export interface FieldProps {
  id: string;
  key: string;
  type: FieldType;
  children: any[];
  options: FieldOptions;
}

const renderTranslatableField = (
  name: string,
  fieldTranslation: Translations
) => {
  return Object.keys(fieldTranslation).map((locale) => {
    return (
      <Fragment>
        <FormLabel>{`${locale}:`}</FormLabel>
        <FormInput
          name={`${name}-${locale}`}
          defaultValue={fieldTranslation[locale] || ''}
          aria-label={`field name for ${locale}`}
        />
      </Fragment>
    );
  });
};

const renderFieldInformationFields = (options: FieldOptions) => {
  return (
    <Fragment>
      <FormLabel>Name:</FormLabel>
      <br />
      {renderTranslatableField('name', options.name)}

      <FormLabel>Description:</FormLabel>
      <br />
      {renderTranslatableField('description', options.description)}

      <FormLabel>Tooltip:</FormLabel>
      <br />
      {renderTranslatableField('tooltip', options.tooltip)}
    </Fragment>
  );
};

const renderStringOptionFields = (options: FieldOptions) => {
  return (
    <Fragment>
      {renderFieldInformationFields(options)}

      <FormLabel>Max Length:</FormLabel>
      <FormInput
        name="maxLength"
        defaultValue={0}
        type="number"
        aria-label={`maximum length allowed for the field`}
      />
    </Fragment>
  );
};

const renderNumericOptionFields = (options: FieldOptions) => {
  return (
    <Fragment>
      {renderFieldInformationFields(options)}

      <FormLabel>Max Value:</FormLabel>
      <FormInput
        name="max"
        defaultValue={0}
        type="number"
        aria-label={`maximum number allowed for the field`}
      />

      <FormLabel>Min Value:</FormLabel>
      <FormInput
        name="min"
        defaultValue={0}
        type="number"
        aria-label={`minimum number allowed for the field`}
      />
    </Fragment>
  );
};

const renderOptions = (props: FieldProps) => {
  switch (props.type) {
    case FieldType.string:
      return renderStringOptionFields(props.options);
    case FieldType.integer:
    case FieldType.float:
      return renderNumericOptionFields(props.options);
    default:
      return <Fragment />;
  }
};

const handleChange = (
  fieldState: FieldProps,
  setFieldState: Dispatch<SetStateAction<FieldProps>>,
  event: FormEvent<HTMLFormElement>
) => {
  const newState = { ...fieldState };

  if (
    event.target['name'].startsWith('description') ||
    event.target['name'].startsWith('name') ||
    event.target['name'].startsWith('tooltip')
  ) {
    const name = event.target['name'].split('-')[0];
    const locale = event.target['name'].split('-')[1];
    newState.options[name][locale] = event.target['value'];
    setFieldState(newState);
  } else if (
    event.target['name'] === 'maxLength' ||
    event.target['name'] === 'max' ||
    event.target['name'] === 'min'
  ) {
    newState.options[event.target['name']] = event.target['value'];
    setFieldState(newState);
  } else {
    newState[event.target['name']] = event.target['value'];
    setFieldState(newState);
  }
};

export const Field: FunctionComponent<FieldProps> = (props) => {
  const [fieldState, setFieldState] = useState<FieldProps>(props);
  return (
    <Fragment>
      <div className={'container'}>
        <div className={'row'}>
          <div className={'col-1 mb-5'}></div>
          <div className={'col-5 mb-5'}>
            <h5>Field State:</h5>
            <p>
              <pre>{JSON.stringify(fieldState, null, 3)}</pre>
            </p>
          </div>
          <div className={'col-5 mb-5'}>
            <form
              onChange={(event) => {
                handleChange(fieldState, setFieldState, event);
              }}
            >
              <FormLabel>ID:</FormLabel>
              <FormInput
                name="id"
                defaultValue={fieldState.id}
                aria-label="field id"
              />

              <FormLabel>Type:</FormLabel>
              <select
                name="type"
                className="form-select"
                aria-label="select type"
              >
                <option>Select Type</option>
                {Object.keys(FieldType).map((fieldType) => {
                  return (
                    <Fragment>
                      {fieldState.type === fieldType ? (
                        <option selected value={fieldType}>
                          {FieldType[fieldType]}
                        </option>
                      ) : (
                        <option value={fieldType}>
                          {FieldType[fieldType]}
                        </option>
                      )}
                    </Fragment>
                  );
                })}
              </select>

              {renderOptions(fieldState)}
            </form>
          </div>
          <div className={'col-1 mb-5'}></div>
        </div>
      </div>
    </Fragment>
  );
};
