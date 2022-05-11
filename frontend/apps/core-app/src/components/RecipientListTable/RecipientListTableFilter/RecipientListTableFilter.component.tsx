import React from 'react';
import { FormControl, Input } from 'native-base';

type Props = {
  onChange: (v: string) => void;
  value: string;
};

export const RecipientListTableFilterComponent: React.FC<Props> = ({
  onChange,
  value,
}) => {
  return (
    <FormControl>
      <FormControl.Label>Beneficiary Name</FormControl.Label>
      <Input
        type="text"
        placeholder="Search"
        value={value || ''}
        onChangeText={onChange}
        testID="recipient-list-table-filter"
      />
    </FormControl>
  );
};
