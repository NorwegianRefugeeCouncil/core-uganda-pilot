import { storiesOf } from '@storybook/react';
import { SingleSelect, MultiSelect } from './select.component';

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
          <SingleSelect options={options} />
        </div>
      </div>
      <div className={'container'}>
        <div className={'col-12 mb-4'}>
          <h3>Multi Select</h3>
          <MultiSelect options={options} disabled />
        </div>
      </div>
    </>
  );
});
