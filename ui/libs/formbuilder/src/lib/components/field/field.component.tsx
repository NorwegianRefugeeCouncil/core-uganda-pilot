import React, { Dispatch, FormEvent, Fragment, FunctionComponent, SetStateAction, useState } from 'react';
import { Button, FormCheck, FormCheckInput, FormCheckLabel, FormControl, FormGroup, FormLabel } from '@core/ui-toolkit';
import './field.css';
import { TranslatedStrings } from '@core/api-client';

export interface SelectOption {
  key: string;
  value: string;
}

export interface RadioOption {
  key: string;
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
  key: string;
  type: FieldType;
  children: any[];
  name: TranslatedStrings;
  description: TranslatedStrings;
  tooltip: TranslatedStrings;
  min?: number | string;
  max?: number | string;
  maxLength?: number;
  regex?: string;
  required?: boolean;
  disabled?: boolean;
  hidden?: boolean;
  default?: any;
}

type ListKind = 'radio' | 'select';

const renderListField = (
  items: RadioOption[] | SelectOption[],
  kind: ListKind
) => {
  return kind === 'radio' ? (
    <Fragment>
      <Button>
        <i className='bi bi-plus-circle-fill' />
      </Button>
      <ul className='list-group'>
        <li className='list-group-item'>
          <FormLabel>Label:</FormLabel>
          <FormControl
            name='option-key'
            defaultValue=''
            aria-label='a radio option key for the field'
          />
        </li>
      </ul>
    </Fragment>
  ) : (
    <Fragment>
      <ul className='list-group'>
        <li className='list-group-item'>
          <FormLabel>Label:</FormLabel>
          <FormControl
            name='option-key'
            defaultValue=''
            aria-label='a radio option key for the field'
          />
          <FormLabel>Value:</FormLabel>
          <FormControl
            name='option-value'
            defaultValue=''
            aria-label='a radio option key for the field'
          />
        </li>
      </ul>
    </Fragment>
  );
};

const renderCheckboxField = (
  checkedState: boolean,
  name: string,
  label: string
) => {
  return (
    <Fragment>
      {checkedState ? (
        <FormCheck>
          <FormCheckLabel>{label}</FormCheckLabel>
          <FormCheckInput type='checkbox' name={name} checked />
        </FormCheck>
      ) : (
        <FormCheck>
          <FormCheckLabel>{label}</FormCheckLabel>
          <FormCheckInput type='checkbox' name={name} />
        </FormCheck>
      )}
    </Fragment>
  );
};

const renderGenericValidationFields = (props: FieldProps) => {
  const { required, disabled, hidden } = props;
  return (
    <Fragment>
      {renderCheckboxField(required, 'required', 'Required')}
      {renderCheckboxField(disabled, 'disabled', 'Disabled')}
      {renderCheckboxField(hidden, 'hidden', 'Hidden')}
      <FormLabel>Regex:</FormLabel>
      <FormControl
        name='regex'
        defaultValue=''
        aria-label='regular expression used to validate the field'
      />
    </Fragment>
  );
};

const renderTranslatableField = (
  name: string,
  fieldTranslation: TranslatedStrings
) => {
  if (!fieldTranslation) {
    return <div />;
  }
  return (
    <ul>
      {fieldTranslation.map(translation => {
        const locale = translation.locale;
        const value = translation.value;
        return (
          <li>
            <i className='bi bi-translate' />
            <FormLabel
              style={{ marginLeft: 5 + 'px' }}
            >{`${locale}:`}</FormLabel>
            {locale === 'ar' ? (
              <FormControl
                style={{ direction: 'rtl' }}
                name={`${name}-${locale}`}
                defaultValue={value}
                aria-label={`field name for ${locale}`}
              />
            ) : (
              <FormControl
                name={`${name}-${locale}`}
                defaultValue={value}
                aria-label={`field name for ${locale}`}
              />
            )}
          </li>
        );
      })}
    </ul>
  );
};

const renderGenericOptionFields = (props: FieldProps) => {
  const { name, description, tooltip } = props;
  return (
    <Fragment>
      <FormLabel>Name:</FormLabel>
      <br />
      {renderTranslatableField('name', name)}

      {renderGenericValidationFields(props)}

      <FormLabel>Description:</FormLabel>
      <br />
      {renderTranslatableField('description', description)}

      <FormLabel>Tooltip:</FormLabel>
      <br />
      {renderTranslatableField('tooltip', tooltip)}
    </Fragment>
  );
};

const renderStringOptionFields = (props: FieldProps) => {
  return (
    <Fragment>
      {renderGenericOptionFields(props)}

      <FormLabel>Max Length:</FormLabel>
      <FormControl
        name='maxLength'
        defaultValue={0}
        type='number'
        aria-label={`maximum length allowed for the field`}
      />
    </Fragment>
  );
};

const renderNumericOptionFields = (props: FieldProps) => {
  return (
    <Fragment>
      {renderGenericOptionFields(props)}

      <FormLabel>Max Value:</FormLabel>
      <FormControl
        name='max'
        defaultValue={0}
        type='number'
        aria-label={`maximum number allowed for the field`}
      />

      <FormLabel>Min Value:</FormLabel>
      <FormControl
        name='min'
        defaultValue={0}
        type='number'
        aria-label={`minimum number allowed for the field`}
      />
    </Fragment>
  );
};

const renderCheckboxOptionFields = (props: FieldProps) => {
  return <Fragment>{renderGenericOptionFields(props)}</Fragment>;
};

const renderRadioOptionFields = (props: FieldProps) => {
  return (
    <Fragment>
      {renderGenericOptionFields(props)}
      {renderListField([] as RadioOption[], 'radio')}
    </Fragment>
  );
};

const renderSelectOptionFields = (props: FieldProps) => {
  return (
    <Fragment>
      {renderGenericOptionFields(props)}
      {renderListField([] as SelectOption[], 'select')}
    </Fragment>
  );
};

const renderOptions = (props: FieldProps) => {
  switch (props.type) {
    case FieldType.string:
      return renderStringOptionFields(props);
    case FieldType.integer:
    case FieldType.float:
      return renderNumericOptionFields(props);
    case FieldType.checkbox:
      return renderCheckboxOptionFields(props);
    case FieldType.radio:
      return renderRadioOptionFields(props);
    case FieldType.select:
    case FieldType.multiselect:
      return renderSelectOptionFields(props);
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

    const translatedString = newState[name] as TranslatedStrings;
    const idx = translatedString.findIndex(s => s.locale === locale);
    if (idx !== -1) {
      translatedString[idx].value = event.target['value'];
    }

  } else if (
    event.target['name'] === 'maxLength' ||
    event.target['name'] === 'max' ||
    event.target['name'] === 'min' ||
    event.target['name'] === 'regex'
  ) {
    newState[event.target['name']] = event.target['value'];
    setFieldState(newState);
  } else if (
    event.target['name'] === 'required' ||
    event.target['name'] === 'disabled' ||
    event.target['name'] === 'hidden'
  ) {
    if (event.target['checked'] === true) {
      newState[event.target['name']] = true;
    } else {
      newState[event.target['name']] = false;
    }
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
            {/* <Card>
              <Card.Body> */}
            <form
              onChange={(event) => {
                handleChange(fieldState, setFieldState, event);
              }}
            >

              <FormGroup controlId={'key'}>
                <FormLabel>Key</FormLabel>
                <FormControl
                  name='key'
                  defaultValue={fieldState.key}
                  aria-label='field key'
                />
              </FormGroup>

              Type:
              <select
                name='type'
                className='form-select'
                aria-label='select type'
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
            {/* </Card.Body>
            </Card> */}
          </div>
          <div className={'col-1 mb-5'}></div>
        </div>
      </div>
    </Fragment>
  );
};
