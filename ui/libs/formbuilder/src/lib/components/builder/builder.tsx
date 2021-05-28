import * as React from 'react';
import {
  addField,
  patchField,
  removeField,
  removeTranslation,
  setTranslation,
  StateSlice,
  TranslationType
} from '../../reducers';
import { useDispatch, useSelector } from 'react-redux';
import { FieldType, FormElement, TranslatedString, TranslatedStrings } from '@core/api-client';
import { Button, Card, FormGroup, FormLabel, FormSelect } from '@core/ui-toolkit';
import { ChangeEvent, useEffect, useRef, useState } from 'react';
import { CardBody } from '@core/ui-toolkit';

export const Builder: React.FC = (props, context) => {
  return <div></div>;
};


type BuilderContainerProps = {
  root: FormElement
}

export const BuilderContainer: React.FC<BuilderContainerProps> = (props, context) => {

  const { root } = props;

  // const root = useSelector<StateSlice, FormElement>(state => state.formBuilder.root);
  const dispatch = useDispatch();

  const doAddField = () => {
    dispatch(addField({ path: '/root', field: {} }));
  };

  return <div>
    {root?.children?.map((c, idx) => {
      const path = 'root/children/' + idx + '/';
      return <FormElementContainer
        key={idx}
        element={c}
        path={path}
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


export const FormElementContainer: React.FC<FormContainerProps> = (props, context) => {

  const dispatch = useDispatch();

  const setKey = (key: string) => {
    dispatch(patchField({
      path: props.path, field: { key }
    }));
  };

  const setType = (type: FieldType) => {
    dispatch(patchField({
      path: props.path, field: { type }
    }));
  };

  const setRequired = (required: boolean) => {
    dispatch(patchField({
      path: props.path, field: { required }
    }));
  };

  const doSetTranslation = (type: TranslationType, locale: string, value: string) => {
    dispatch(setTranslation({
      path: props.path,
      value: value,
      locale: locale,
      type: type
    }));
  };

  const doRemoveTranslation = (type: TranslationType, locale: string) => {
    dispatch(removeTranslation({
      path: props.path,
      locale: locale,
      type: type
    }));
  };


  const doRemoveElement = () => {
    dispatch(removeField({ path: props.path }));
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

export const OptionsView: React.FC<FormElementProps> = props => {

  const { required } = props.element;
  const { setRequired } = props;

  const { tooltip } = props.element;
  const {
    selectedTranslationType,
    setSelectedTranslationType,
    setTranslation,
    removeTranslation
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
    <select className={'form-select'} value={type}
            onChange={(ev) => setType(ev.target.value as FieldType)}>
      <option value={'text'}>Text</option>
      <option value={'checkbox'}>Checkbox</option>
    </select>
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
