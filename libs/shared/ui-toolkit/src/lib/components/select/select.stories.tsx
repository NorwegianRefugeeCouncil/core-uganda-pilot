import { storiesOf } from '@storybook/react';
import Select from './select.component';
import VanillaSelect from 'react-select';

const options = [
  { value: 'A', label: 'Dingle' },
  { value: 'B', label: 'Dongle' },
  { value: 'C', label: 'Dangle' },
  { value: 'D', label: 'Doogle' },
  { value: 'E', label: 'Diggle' },
  { value: 'F', label: 'Bibble' },
];

storiesOf('Select', module).add('default', () => {
  return (
    <>
      <div className={'container'}>
        <div className={'col-12 mb-4'}>
          <h3>Single Select</h3>
          <Select options={options} />
        </div>
      </div>
      <div className={'container'}>
        <div className={'col-12 mb-4'}>
          <h3>Disabled Select</h3>
          <Select options={options} isDisabled />
        </div>
      </div>
      <div className={'container'}>
        <div className={'col-12 mb-4'}>
          <h3>Multi Select</h3>
          <Select options={options} isMulti />
        </div>
      </div>
      <div className={'container'}>
        <div className={'col-12 mb-4'}>
          <h3>Multi Select</h3>
          <VanillaSelect options={options} isMulti />
        </div>
      </div>
    </>
  );
});
