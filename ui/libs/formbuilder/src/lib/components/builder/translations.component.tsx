import * as React from 'react';
import { ChangeEvent, useEffect, useRef } from 'react';
import { TranslationType } from '../../reducers';
import { TranslatedString, TranslatedStrings } from '@core/api-client';
import { Card, CardBody } from '@core/ui-toolkit';


type EditTranslationComboProps = {
  selectedTranslationType: TranslationType | ''
  setSelectedTranslationType: (t: TranslationType | '') => void
  translationType: TranslationType
  translations: TranslatedStrings
  setTranslation: (type: TranslationType, locale: string, value: string) => void
  removeTranslation: (type: TranslationType, locale: string) => void
  label: string
}

/**
 * Renders an input control for editing a translated value
 * The editor will display a dropdown if the translationType equals to the selectedTranslationType
 * A default translation will be shown in the input (by default, english, or the first non-empty translation)
 * @param props
 * @constructor
 */
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

/**
 * Renders a dropdown that allows the user to edit the translations
 * @param props
 * @constructor
 */
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


/**
 * Renders the input control that displays the default translation
 * as well as binds the input onClick to display the dropdown
 * @param props
 * @constructor
 */
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

/**
 * Displays the translations as well as a control to add a translation
 * @param props
 * @constructor
 */
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

/**
 * Renders the component that allows the user to add a translation
 * @param props
 * @constructor
 */
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

/**
 * Renders a single translated string
 * Allows the user to edit or remove the translation
 * @param str
 * @param removeTranslation
 * @param setTranslation
 */
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

/**
 * Utility method to close the dropdown if the user clicks elsewhere
 * @param ref
 * @param onClickOutside
 */
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
