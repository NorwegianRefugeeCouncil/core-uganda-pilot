import React from 'react';
import { FormControl, Input } from 'native-base';

type Props = {
  handleChange: (v: string) => void;
  value: string;
};

export const GlobalTableFilterComponent: React.FC<Props> = ({
  handleChange,
  value,
}) => {
  return (
    <FormControl>
      <FormControl.Label>Beneficiary Name</FormControl.Label>
      <Input
        type="text"
        placeholder="Search"
        value={value || ''}
        onChangeText={handleChange}
      />
    </FormControl>
  );
};
