import { Action } from '@reduxjs/toolkit';
import { removeIndexedValue, removeValue, setIndexedValue, setValue, TranslationType } from '../../reducers';
import { FieldType, FormElement } from '@core/api-client';
import * as React from 'react';
import { useState } from 'react';
import { useDispatch } from 'react-redux';
import { FormElementComponent } from './element.builder';

export type FormElementContainerProps = {
  element: FormElement
  path: string
}

/**
 * Utility method to either set or remove a numeric field value
 * @param n the input value
 * @param path the path of the value
 */
const setOrRemoveNumber = (n: any, path: string): Action => {
  if (typeof n !== 'number') {
    return removeValue({ path });
  } else {
    return setValue({ path, value: n });
  }
};

/**
 * Renders the container for displaying a single form element builder
 * @param props
 * @constructor
 */
export const FormElementContainer: React.FC<FormElementContainerProps> = (props) => {

  const { path } = props;

  const dispatch = useDispatch();

  // Sets the 'key' property
  const setKey = (key: string) => {
    dispatch(setValue({ path: path + '.key', value: key }));
  };

  // Sets the 'type' property
  const setType = (type: FieldType) => {
    dispatch(setValue({ path: path + '.type', value: type }));
  };

  // sets the 'required' property
  const setRequired = (required: boolean) => {
    dispatch(setValue({ path: path + '.required', value: required }));
  };

  // sets the 'minimum' property
  const setMinimum = (min: number) => {
    dispatch(setOrRemoveNumber(min, path + '.min'));
  };

  // sets the 'maximum' property
  const setMaximum = (max: number) => {
    dispatch(setOrRemoveNumber(max, path + '.max'));
  };

  // sets the 'minLength' property
  const setMinLength = (minLength: number) => {
    dispatch(setOrRemoveNumber(minLength, path + '.minLength'));
  };

  // sets the 'maxLength' property
  const setMaxLength = (maxLength: number) => {
    dispatch(setOrRemoveNumber(maxLength, path + '.maxLength'));
  };

  // sets a translation given a locale
  const doSetTranslation = (type: TranslationType, locale: string, value: string) => {
    dispatch(setIndexedValue({
      path: path + '.' + type, key: 'locale', value: {
        locale: locale,
        value: value
      }
    }));
  };

  // removes a translation given a locale
  const doRemoveTranslation = (type: TranslationType, locale: string) => {
    dispatch(removeIndexedValue({
      path: path + '.' + type, key: 'locale', keyValue: locale
    }));
  };

  // action to remove the form element itself
  const doRemoveElement = () => {
    dispatch(removeValue({ path: props.path }));
  };

  // the translation editor opens a dropdown to edit the translations
  // only one dropdown can be opened at a time
  // this value keeps track of the opened translation dropdown
  // if empty string, no dropdown is opened
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
