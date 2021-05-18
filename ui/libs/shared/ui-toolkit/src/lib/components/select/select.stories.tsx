import Select from './select.component';

const options = [
  { value: 'A', label: 'Dingle' },
  { value: 'B', label: 'Dongle' },
  { value: 'C', label: 'Dangle' },
  { value: 'D', label: 'Doogle' },
  { value: 'E', label: 'Diggle' },
  { value: 'F', label: 'Bibble' },
];

export default {
  title: 'Select',
  component: Select,
};

export const Basic = () => (
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
  </>
);
