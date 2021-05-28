import { FormElementContainerProps } from './element.container';
import { FieldType } from '@core/api-client';
import { TranslationType } from '../../reducers';
import * as React from 'react';
import { useState } from 'react';
import { EditTranslationCombo } from './translations.component';
import { Button, Card, CardBody } from '@core/ui-toolkit';
import Select from 'react-select';

export type FormElementProps = FormElementContainerProps & {
  removeElement: () => void
  setKey: (key: string) => void
  setType: (type: FieldType) => void
  setRequired: (required: boolean) => void
  setTranslation: (type: TranslationType, locale: string, value: string) => void
  removeTranslation: (type: TranslationType, locale: string) => void
  selectedTranslationType: TranslationType | ''
  setSelectedTranslationType: (t: TranslationType | '') => void
  setMinimum: (minimum: number) => void
  setMaximum: (maximum: number) => void
  setMinLength: (minLength: number) => void
  setMaxLength: (maxLength: number) => void
}

/**
 * Represents the currently selected view for the form element
 */
type SelectedView = 'edit' | 'options'

/**
 * Renders the component that displays the user with options for
 * configuring a form control
 * @param props
 * @param context
 * @constructor
 */
export const FormElementComponent: React.FC<FormElementProps> = (props, context) => {

  const [selectedView, setSelectedView] = useState<SelectedView>('edit');

  return <Card className={'mb-2 shadow-sm bg-light'}>
    <CardBody>
      <div className={'row'}>

        {/* Button to delete the element */}
        <div className={'col-auto order-1 order-md-2 px-0 pe-1'}>
          <button
            className={'btn'}
            onClick={() => props.removeElement()}
          >
            <i className={'bi bi-trash'} />
          </button>
        </div>


        <div className={'col-auto flex-grow-1 col-md-5 order-0 order-md-0'}>
          {/* Button to select the Edit view */}
          <Button
            className={'shadow-sm'}
            theme={selectedView === 'edit' ? 'primary' : 'secondary'}
            outline={selectedView !== 'edit'}
            onClick={() => setSelectedView('edit')}
          >
            Edit
          </Button>
          {/* Button to select the Options view */}
          <Button
            className={'ms-2 shadow-sm'}
            theme={selectedView === 'options' ? 'primary' : 'secondary'}
            outline={selectedView !== 'options'}
            onClick={() => setSelectedView('options')}
          >
            Options
          </Button>
        </div>

        {/* Input to set the Key of the form control */}
        <div className={'col-12 col-md-auto mt-3 mt-md-0 order-1 order-md-1'}>
          <div className={'input-group mb-3'}>
            <span className={'input-group-text'}>Key</span>
            <input
              type={'text'}
              className={'form-control font-monospace'}
              placeholder={'key'}
              value={props.element.key}
              onChange={(ev) => props.setKey(ev.target.value)}
            />
          </div>
        </div>
      </div>

      {renderSelectedView(selectedView, props)}

      <div className={'input-group'}>

      </div>

    </CardBody>
  </Card>;
};

/**
 * Renders the form control selected view (Edit or Options)
 * @param selectedView
 * @param props
 */
export const renderSelectedView = (selectedView: SelectedView, props: FormElementProps) => {
  if (selectedView === 'edit') {
    return <EditView {...props} />;
  }
  if (selectedView === 'options') {
    return <OptionsView {...props} />;
  }
  return null;
};

/**
 * Utility method that returns the number if it is valid, or an empty string
 * @param value
 */
const numberOrEmpty = (value: any): number | string => {
  if (typeof value === 'number') {
    return value;
  }
  return '';
};

/**
 * Represents the "Options" tab view
 * @param props
 * @constructor
 */
export const OptionsView: React.FC<FormElementProps> = props => {

  const { required } = props.element;
  const { setRequired } = props;

  const { tooltip } = props.element;
  const {
    selectedTranslationType,
    setSelectedTranslationType,
    setTranslation,
    removeTranslation,
    setMinimum,
    setMaximum,
    setMinLength,
    setMaxLength
  } = props;

  return (
    <>
      <div className={'form-check mb-3'}>
        <input
          className={'form-check-input'}
          type={'checkbox'}
          checked={!!required}
          onChange={() => setRequired(!required)}
        />
        <label className={'form-check-label'}>
          Required
        </label>
      </div>

      <EditTranslationCombo
        label={'Tooltip'}
        translations={tooltip}
        selectedTranslationType={selectedTranslationType}
        setSelectedTranslationType={setSelectedTranslationType}
        translationType={'tooltip'}
        setTranslation={setTranslation}
        removeTranslation={removeTranslation}
      />

      <div className='input-group mb-3'>
        <span style={{ width: '8rem' }} className='input-group-text'>Minimum</span>
        <input type='number'
               className='form-control'
               placeholder='Minimum value'
               value={numberOrEmpty(props.element.min)}
               onChange={(ev) => setMinimum(parseFloat(ev.target.value))}
        />
        <button
          className={'btn btn-outline-secondary'}
          onClick={() => setMinimum(null)}
        ><span className={'bi bi-x'} /></button>
      </div>

      <div className='input-group mb-3'>
        <span style={{ width: '8rem' }} className='input-group-text'>Maximum</span>
        <input type='number'
               className='form-control'
               placeholder='Maximum value'
               value={numberOrEmpty(props.element.max)}
               onChange={(ev) => setMaximum(parseFloat(ev.target.value))}
        />
        <button
          className={'btn btn-outline-secondary'}
          onClick={() => setMaximum(null)}
        ><span className={'bi bi-x'} /></button>
      </div>


      <div className='input-group mb-3'>
        <span style={{ width: '8rem' }} className='input-group-text'>Minimum Length</span>
        <input type='number'
               className='form-control'
               placeholder='Minimum length of the value'
               value={numberOrEmpty(props.element.minLength)}
               onChange={(ev) => setMinLength(parseFloat(ev.target.value))}
        />
        <button
          className={'btn btn-outline-secondary'}
          onClick={() => setMinLength(null)}
        ><span className={'bi bi-x'} /></button>
      </div>


      <div className='input-group mb-3'>
        <span style={{ width: '8rem' }} className='input-group-text'>Maximum Length</span>
        <input type='number'
               className='form-control'
               placeholder='Maximum length of the value'
               value={numberOrEmpty(props.element.maxLength)}
               onChange={(ev) => setMaxLength(parseFloat(ev.target.value))}
        />
        <button
          className={'btn btn-outline-secondary'}
          onClick={() => setMaxLength(null)}
        ><span className={'bi bi-x'} /></button>
      </div>


    </>
  );
};


type FormElementTypeSelectorProps = {
  type?: FieldType
  setType: (type: FieldType) => void
}

/**
 * Component that allows a user to select what type of component it is configuring
 * ShortText, DateTime, etc.
 * @param props
 * @constructor
 */
export const FormElementTypeSelector: React.FC<FormElementTypeSelectorProps> = props => {

  const { type, setType } = props;

  const options = [{ value: 'shortText', label: 'Short Text' },
    { value: 'longText', label: 'Long Text' },
    { value: 'checkBox', label: 'Checkbox' },
    { value: 'radio', label: 'Radio Buttons' },
    { value: 'date', label: 'Date' },
    { value: 'dateTime', label: 'Date Time' },
    { value: 'time', label: 'Time' }];

  const selected = options.find((o) => o.value === type);

  return <div className={'input-group mb-3'}>
    <span style={{ width: '8rem' }} className={'input-group-text'}>Type</span>
    <Select
      value={selected}
      styles={{
        container: () => ({
          flex: 1
        }),
        control: (provided) => ({
          ...provided,
          borderTopLeftRadius: 0,
          borderBottomLeftRadius: 0
        })
      }}
      options={[
        { value: 'shortText', label: 'Short Text' },
        { value: 'longText', label: 'Long Text' },
        { value: 'checkBox', label: 'Checkbox' },
        { value: 'radio', label: 'Radio Buttons' },
        { value: 'date', label: 'Date' },
        { value: 'dateTime', label: 'Date Time' },
        { value: 'time', label: 'Time' }
      ]}
      onChange={(ev) => {
        setType(ev.value as FieldType);
      }}
    />
  </div>;
};

/**
 * Represents the Edit view of the form control editor
 * Displays basic information such as description and label
 * @param props
 * @constructor
 */
export const EditView: React.FC<FormElementProps> = props => {

  const { type, label } = props.element;
  const {
    setType,
    selectedTranslationType,
    setSelectedTranslationType,
    setTranslation,
    removeTranslation
  } = props;

  return (
    <>
      <FormElementTypeSelector
        type={type}
        setType={setType}
      />

      <EditTranslationCombo
        label={'Label'}
        translations={label}
        selectedTranslationType={selectedTranslationType}
        setSelectedTranslationType={setSelectedTranslationType}
        translationType={'label'}
        setTranslation={setTranslation}
        removeTranslation={removeTranslation}
      />

    </>
  );
};
