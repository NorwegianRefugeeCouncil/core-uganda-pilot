import { storiesOf } from '@storybook/react';
import { Select } from './select.component';

const options = [
  { value: 'A', label: 'Dingle' },
  { value: 'B', label: 'Dongle' },
  { value: 'C', label: 'Dangle' },
];

storiesOf('Select', module).add('default', () => {
  return (
    <>
      <div className={'container'}>
        <div className={'col-12 mb-4'}>
          <h3>Empty Select</h3>
          <Select />
        </div>
      </div>
      <div className={'container'}>
        <div className={'col-12 mb-4'}>
          <h3>Single Select</h3>
          <Select options={options} />
        </div>
      </div>
    </>
  );
});
