import React, { useState } from 'react';
import { Button, Text } from 'react-native-paper';
import { Platform } from 'react-native';
import { DatePickerModal } from 'react-native-paper-dates';

import { darkTheme } from '../../constants/theme';

import { InputProps } from './FormControl';

const DateInput: React.FC<InputProps> = ({
  fieldDefinition,
  style,
  value,
  onChange,
  onBlur,
  error,
  invalid,
  isTouched,
  isDirty,
  isMultiple,
  isQuantity,
}) => {
  // export function DateModal(props: { date: Date; setDate: React.Dispatch<React.SetStateAction<Date>>; }) {
  // const [date, setDate] = props.date != null && props.setDate != null ? [props.date, props.setDate] : useState(new Date(Date.now()));
  // const [date, setDate] = useState(new Date(Date.now()));
  const [isOpen, setIsOpen] = useState(false);

  const show = () => setIsOpen(true);
  const hide = () => setIsOpen(false);

  let datePicker;

  if (Platform.OS === 'web') {
    datePicker = <input type="date" value={toDateString(value)} onChange={(event) => onChange(toDate(event.target.value))} />;
  } else {
    datePicker = (
      <>
        {fieldDefinition.name && <Text theme={darkTheme}>{fieldDefinition.name}</Text>}
        {fieldDefinition.description && (
          <Text theme={darkTheme} style={{ fontSize: 10 }}>
            {fieldDefinition.description}
          </Text>
        )}
        <Button onPress={show}>{toDateString(value)}</Button>
        <DatePickerModal
          locale="en"
          mode="single"
          onConfirm={(p) => {
            hide();
            onChange(p.date as Date);
          }}
          onDismiss={hide}
          visible={isOpen}
        />
      </>
    );
  }

  return <>{datePicker}</>;
};

function toDateString(date: Date): string {
  console.log('TO DATE STRING', date);
  if (!date) {
    return 'yyyy-mm-dd';
  }
  const y = date.getFullYear();
  const zeroPad = (n: number) => (n < 9 ? `0${n + 1}` : n + 1);
  const m = zeroPad(date.getMonth());
  const d = date.getDate();
  return `${y}-${m}-${d}`;
}

function toDate(yyyymmdd: string): Date {
  const [yyyy, mm, dd] = yyyymmdd.split('-').map((s) => (s[0] === '0' ? +s[1] : +s));
  return new Date(yyyy, mm - 1, dd);
}

export default DateInput;
