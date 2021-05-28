import * as React from 'react';
import {
  addFormElement, addValue,
  patchFormElement,
  removeFormElement, removeIndexedValue,
  removeTranslation, removeValue, setFormDefinition, setIndexedValue,
  setTranslation, setValue,
  StateSlice,
  TranslationType
} from '../../reducers';
import { useDispatch, useSelector } from 'react-redux';
import {
  FieldType,
  FormDefinition,
  FormDefinitionVersion,
  FormElement,
  TranslatedString,
  TranslatedStrings
} from '@core/api-client';
import { Button, Card, FormGroup, FormLabel, FormSelect } from '@core/ui-toolkit';
import { ChangeEvent, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { CardBody } from '@core/ui-toolkit';
import { Action } from '@reduxjs/toolkit';
import Select from 'react-select';


type BuilderProps = {
  formDefinition: FormDefinition
}

export const Builder: React.FC<BuilderProps> = (props, context) => {

  const givenFormDefinition = props.formDefinition;

  if (!givenFormDefinition) {
    return <div />;
  }

  const dispatch = useDispatch();

  const formDefinition = useSelector<StateSlice, FormDefinition>(state => state.formBuilder.formDefinition);

  const [currentVersionName, setCurrentVersionName] = useState(formDefinition?.spec?.versions[0].name);

  useEffect(() => {
    dispatch(setFormDefinition({ formDefinition: givenFormDefinition }));
    setCurrentVersionName(givenFormDefinition.spec.versions[0].name);
  }, [givenFormDefinition]);

  const currentVersion = useMemo(() => {
    return formDefinition?.spec?.versions?.find(v => v.name === currentVersionName);
  }, [formDefinition, currentVersionName]);

  useCallback(args => {
    console.log(args);
  }, [currentVersionName]);

  return <>

    <Card>
      <CardBody>

        <div className={'input-group mb-3'}>
          <span style={{ width: '5rem' }} className='input-group-text'>Group</span>
          <input disabled type='text' className='form-control bg-light' value={formDefinition.spec.group} />
        </div>

        <div className={'input-group mb-3'}>
          <span style={{ width: '5rem' }} className='input-group-text'>Kind</span>
          <input disabled type='text' className='form-control bg-light' value={formDefinition.spec.names.kind} />
        </div>

        <div className={'input-group mb-3'}>
          <span style={{ width: '5rem' }} className='input-group-text'>Plural</span>
          <input disabled type='text' className='form-control bg-light' value={formDefinition.spec.names.plural} />
        </div>

        <div className={'input-group mb-3'}>
          <span style={{ width: '5rem' }} className='input-group-text'>Singular</span>
          <input disabled type='text' className='form-control bg-light' value={formDefinition.spec.names.singular} />
        </div>

        <hr />

        <div className='input-group mb-3'>
          <span className='input-group-text' id='basic-addon1'>Version</span>
          <select
            className={'form-select'}
            value={currentVersionName}
            onChange={ev => setCurrentVersionName(ev.target.value)}>
            {formDefinition.spec.versions.map(v => {
              return <option
                key={v.name}
                value={v.name}
              >{v.name}</option>;
            })}
          </select>
        </div>

        <div className={'mt-2'}>
          {renderVersion({
            path: 'spec.versions[0]',
            version: currentVersion
          })}
        </div>

      </CardBody>
    </Card>

  </>;
};


type renderVersionProps = {
  path: string,
  version: FormDefinitionVersion
}

const renderVersion = (props: renderVersionProps) => {
  const { version, path } = props;

  if (!version) {
    return <div />;
  }
  return <div>

    <div className='form-check'>
      <input disabled className='form-check-input' type='checkbox' value='' id='flexCheckDefault'
             checked={version.served} />
      <label className='form-check-label text-dark' htmlFor='flexCheckDefault'>
        Served
      </label>
    </div>

    <div className='form-check'>
      <input disabled className='form-check-input' type='checkbox' value='' id='flexCheckDefault'
             checked={version.storage} />
      <label className='form-check-label text-dark' htmlFor='flexCheckDefault'>
        Storage
      </label>
    </div>

    <div className={'mt-3'}>
      <RootBuilderContainer
        path={path + '.schema.formSchema.root'}
        root={version.schema.formSchema.root}
      />
    </div>
  </div>;
};

type RootBuilderContainerProps = {
  root: FormElement,
  path: string
}

export const RootBuilderContainer: React.FC<RootBuilderContainerProps> = (props, context) => {

  const { root, path } = props;

  const dispatch = useDispatch();

  const doAddField = () => {
    dispatch(addValue({
      path: path + '.children', value: {
        type: 'shortText'
      }
    }));
  };

  return <div>
    {root?.children?.map((c, idx) => {
      return <FormElementContainer
        key={idx}
        element={c}
        path={path + '.children[' + idx + ']'}
      />;
    })}

    <button
      className={'btn btn-primary shadow-sm w-100'}
      onClick={() => doAddField()}
    >
      <i className={'bi bi-plus'} /> Add
    </button>

  </div>;
};

type FormContainerProps = {
  element: FormElement
  path: string
}

type FormElementProps = FormContainerProps & {
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

type SelectedView = 'edit' | 'options'

export const FormElementComponent: React.FC<FormElementProps> = (props, context) => {

  const [selectedView, setSelectedView] = useState<SelectedView>('edit');

  return <Card className={'mb-2 shadow-sm bg-light'}>
    <CardBody>
      <div className={'row'}>
        <div className={'col-auto order-1 order-md-2 px-0 pe-1'}>
          <button
            className={'btn'}
            onClick={() => props.removeElement()}
          >
            <i className={'bi bi-trash'} />
          </button>
        </div>
        <div className={'col-auto flex-grow-1 col-md-5 order-0 order-md-0'}>
          <Button
            className={'shadow-sm'}
            theme={selectedView === 'edit' ? 'primary' : 'secondary'}
            outline={selectedView !== 'edit'}
            onClick={() => setSelectedView('edit')}
          >
            Edit
          </Button>
          <Button
            className={'ms-2 shadow-sm'}
            theme={selectedView === 'options' ? 'primary' : 'secondary'}
            outline={selectedView !== 'options'}
            onClick={() => setSelectedView('options')}
          >
            Options
          </Button>
        </div>
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

const setOrRemoveNumber = (n: any, path: string): Action => {
  if (typeof n !== 'number') {
    return removeValue({ path });
  } else {
    return setValue({ path, value: n });
  }
};

export const FormElementContainer: React.FC<FormContainerProps> = (props, context) => {

  const { path } = props;

  const dispatch = useDispatch();

  const setKey = (key: string) => {
    dispatch(setValue({ path: path + '.key', value: key }));
  };

  const setType = (type: FieldType) => {
    dispatch(setValue({ path: path + '.type', value: type }));
  };

  const setRequired = (required: boolean) => {
    dispatch(setValue({ path: path + '.required', value: required }));
  };

  const setMinimum = (min: number) => {
    dispatch(setOrRemoveNumber(min, path + '.min'));
  };

  const setMaximum = (max: number) => {
    dispatch(setOrRemoveNumber(max, path + '.max'));
  };

  const setMinLength = (minLength: number) => {
    dispatch(setOrRemoveNumber(minLength, path + '.minLength'));
  };

  const setMaxLength = (maxLength: number) => {
    dispatch(setOrRemoveNumber(maxLength, path + '.maxLength'));
  };

  const doSetTranslation = (type: TranslationType, locale: string, value: string) => {
    dispatch(setIndexedValue({
      path: path + '.' + type, key: 'locale', value: {
        locale: locale,
        value: value
      }
    }));
  };

  const doRemoveTranslation = (type: TranslationType, locale: string) => {
    dispatch(removeIndexedValue({
      path: path + '.' + type, key: 'locale', keyValue: locale
    }));
  };


  const doRemoveElement = () => {
    dispatch(removeValue({ path: props.path }));
  };

  const [selectedTranslationType, setSelectedTranslationType] = useState<TranslationType | ''>('');

  return <div>
    <FormElementComponent
      element={props.element}
      path={props.path}
      removeElement={doRemoveElement}
      setKey={setKey}
      setType={setType}
      setRequired={setRequired}
      setTranslation={doSetTranslation}
      removeTranslation={doRemoveTranslation}
      selectedTranslationType={selectedTranslationType}
      setSelectedTranslationType={setSelectedTranslationType}
      setMinimum={setMinimum}
      setMaximum={setMaximum}
      setMinLength={setMinLength}
      setMaxLength={setMaxLength}
    />
  </div>;
};


export const renderSelectedView = (selectedView: SelectedView, props: FormElementProps) => {
  if (selectedView === 'edit') {
    return <EditView {...props} />;
  }
  if (selectedView === 'options') {
    return <OptionsView {...props} />;
  }
  return null;
};

const numberOrEmpty = (value: any): number | string => {
  if (typeof value === 'number') {
    return value;
  }
  return '';
};

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


type FormElementTypeSelectorProps = {
  type?: FieldType
  setType: (type: FieldType) => void
}

export const FormElementTypeSelector: React.FC<FormElementTypeSelectorProps> = props => {

  const { type, setType } = props;

  return <div className={'input-group mb-3'}>
    <span style={{ width: '8rem' }} className={'input-group-text'}>Type</span>
    <Select
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
        console.log(ev);
      }}
    />
  </div>;
};


type EditTranslationComboProps = {
  selectedTranslationType: TranslationType | ''
  setSelectedTranslationType: (t: TranslationType | '') => void
  translationType: TranslationType
  translations: TranslatedStrings
  setTranslation: (type: TranslationType, locale: string, value: string) => void
  removeTranslation: (type: TranslationType, locale: string) => void
  label: string
}


export const EditTranslationCombo: React.FC<EditTranslationComboProps> = props => {
  const {
    selectedTranslationType,
    translationType,
    setSelectedTranslationType,
    label,
    translations,
    setTranslation,
    removeTranslation
  } = props;

  return (
    <>
      <EditTranslationInput
        label={label}
        translations={translations}
        selectedTranslationType={selectedTranslationType}
        setSelectedTranslationType={setSelectedTranslationType}
        translationType={translationType}
      />

      <EditTranslationDropdown
        selectedTranslationType={selectedTranslationType}
        setSelectedTranslationType={setSelectedTranslationType}
        translationType={translationType}
        translations={translations}
        setTranslation={setTranslation}
        removeTranslation={removeTranslation}
      />
    </>
  );
};

type EditTranslationsDropdownProps = {
  selectedTranslationType: TranslationType | ''
  setSelectedTranslationType: (t: TranslationType | '') => void
  translationType: TranslationType
  translations: TranslatedStrings
  setTranslation: (type: TranslationType, locale: string, value: string) => void
  removeTranslation: (type: TranslationType, locale: string) => void
}

export const EditTranslationDropdown: React.FC<EditTranslationsDropdownProps> = props => {

  const {
    selectedTranslationType,
    translationType,
    setSelectedTranslationType,
    translations,
    setTranslation,
    removeTranslation
  } = props;

  if (selectedTranslationType !== translationType) {
    return null;
  }

  return <TranslatedStringsComponent
    translatedStrings={translations}
    close={() => setSelectedTranslationType('')}
    removeTranslation={locale => removeTranslation(translationType, locale)}
    setTranslation={(locale, value) => setTranslation(translationType, locale, value)}
  />;

};

type EditTranslationInputProps = {
  selectedTranslationType: TranslationType | ''
  setSelectedTranslationType: (t: TranslationType | '') => void
  translations: TranslatedStrings
  translationType: TranslationType
  label: string
}

export const EditTranslationInput: React.FC<EditTranslationInputProps> = props => {

  const { selectedTranslationType, setSelectedTranslationType, translationType } = props;
  let { translations, label } = props;

  let defaultTranslation: string;
  if (!translations) {
    translations = [];
  }
  for (let translation of translations) {
    if (translation.locale === 'en') {
      defaultTranslation = translation.value;
      break;
    }
  }
  if (!defaultTranslation) {
    for (let translation of translations) {
      if (translation.value) {
        defaultTranslation = translation.value;
        break;
      }
    }
  }
  if (!defaultTranslation) {
    defaultTranslation = '';
  }

  return (<div className='input-group mb-3'>
    <span style={{ width: '8rem' }} className='input-group-text'>{label}</span>
    <input type='button' className='form-control text-start' placeholder='Username' aria-label='Username'
           aria-describedby='basic-addon1'
           value={defaultTranslation}
           onMouseDown={(ev) => {
             // open the translations dropdown
             if (selectedTranslationType !== translationType) {
               setSelectedTranslationType(translationType);
             }
           }}
    />
    <span className='input-group-text'>
      <i className='bi bi-translate' />
    </span>
  </div>);
};

type TranslatedStringsProps = {
  translatedStrings: TranslatedStrings,
  setTranslation: (locale: string, value: string) => void,
  removeTranslation: (locale: string) => void,
  close: () => void
}

const TranslatedStringsComponent: React.FC<TranslatedStringsProps> = (
  props) => {
  const { translatedStrings, removeTranslation, setTranslation, close } = props;
  const ref = useRef(null);

  // close the "dropdown" on outside click
  useOutsideAlerter(ref, () => {
    close();
  });

  return <Card className='shadow m-2' ref={ref}>
    <CardBody>
      <div className={'text-end mb-3'}>
        <button
          onClick={() => close()}
          type='button'
          className='btn-close text-end'
          aria-label='Close'
        />
      </div>
      {translatedStrings?.map(translatedString => {

        const setLocaleTranslation = (translationValue) => {
          setTranslation(translatedString.locale, translationValue);
        };

        const removeLocaleTranslation = () => {
          removeTranslation(translatedString.locale);
        };

        return renderTranslatedRow(translatedString, removeLocaleTranslation, setLocaleTranslation);

      })}
      <AddTranslationRow
        setTranslation={setTranslation}
        translatedStrings={translatedStrings}
      />
    </CardBody>
  </Card>;
};


type AddTranslationRowProps = {
  translatedStrings: TranslatedStrings
  setTranslation: (locale: string, value: string) => void,
}

const AddTranslationRow: React.FC<AddTranslationRowProps> = props => {

  const { translatedStrings, setTranslation } = props;

  const allLocales = ['en', 'fr'];
  const currentLocales = translatedStrings ? translatedStrings.map(l => l.locale) : [];
  const availableLocales = allLocales.filter(l => currentLocales.indexOf(l) === -1);

  if (availableLocales?.length === 0) {
    return null;
  }

  const handleTranslationSelected = (ev: ChangeEvent<HTMLSelectElement>) => {
    if (ev.target.value) {
      setTranslation(ev.target.value, '');
    }
  };

  const renderOption = (locale: string) => {
    return <option key={locale} value={locale}>{locale}</option>;
  };

  const renderEmptyOption = () => {
    return <option />;
  };

  return (<div className={'row mt-3'}>
    <label className={'col-5 col-sm-4 col-md-auto col-form-label'}>Add translation</label>
    <div className={'col-7 col-sm-8 col-md-auto flex-md-grow-1'}>
      <select className={'form-select'} value={''} onChange={handleTranslationSelected}>
        {renderEmptyOption()}
        {availableLocales.map(a => renderOption(a))}
      </select>
    </div>
  </div>);
};


const renderTranslatedRow = (
  str: TranslatedString,
  removeTranslation: () => void,
  setTranslation: (value: string) => void
) => {

  return <div className='input-group mb-2' key={str.locale}>
    <span style={{ width: '3.5rem' }} className='input-group-text'>{str.locale}</span>
    <input
      type='text'
      className='form-control'
      value={str.value}
      onChange={ev => setTranslation(ev.target.value)}
    />
    <button
      className='btn btn-outline-secondary'
      type='button'
      id='button-addon1'
      onClick={() => removeTranslation()}
    >
      <i className={'bi bi-x'} />
    </button>
  </div>;
};

function useOutsideAlerter(ref, onClickOutside: () => void) {
  useEffect(() => {
    /**
     * Alert if clicked on outside of element
     */
    function handleClickOutside(event) {
      if (ref.current && !ref.current.contains(event.target)) {
        onClickOutside();
      }
    }

    // Bind the event listener
    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      // Unbind the event listener on clean up
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [ref]);
}
