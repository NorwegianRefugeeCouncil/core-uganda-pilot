import { FunctionComponent } from 'react';
import ReactSelect from 'react-select';

// a thin wrapper over react-select applying our own bootstrap styles

const customStyles = {
  clearIndicator: (provided, state) => ({
    ...provided,
  }),
  container: (provided, state) => ({
    ...provided,
  }),
  control: (provided, state) => ({
    ...provided,
  }),
  dropdownIndicator: (provided, state) => ({
    ...provided,
  }),
  group: (provided, state) => ({
    ...provided,
  }),
  groupHeading: (provided, state) => ({
    ...provided,
  }),
  indicatorsContainer: (provided, state) => ({
    ...provided,
  }),
  indicatorSeparator: (provided, state) => ({
    ...provided,
  }),
  input: (provided, state) => ({
    ...provided,
  }),
  loadingIndicator: (provided, state) => ({
    ...provided,
  }),
  loadingMessage: (provided, state) => ({
    ...provided,
  }),
  menu: (provided, state) => ({
    ...provided,
  }),
  menuList: (provided, state) => ({
    ...provided,
  }),
  menuPortal: (provided, state) => ({
    ...provided,
  }),
  multiValue: (provided, state) => ({
    ...provided,
  }),
  multiValueLabel: (provided, state) => ({
    ...provided,
  }),
  multiValueRemove: (provided, state) => ({
    ...provided,
  }),
  noOptionsMessage: (provided, state) => ({
    ...provided,
  }),
  option: (provided, state) => ({
    ...provided,
    backgroundColor: state.isSelected ? 'var()' : '$secondary',
  }),
  placeholder: (provided, state) => ({
    ...provided,
  }),
  singleValue: (provided, state) => ({
    ...provided,
  }),
  valueContainer: (provided, state) => ({
    ...provided,
  }),
};

export const Select: FunctionComponent = (props) => {
  return <ReactSelect {...props} styles={customStyles} />;
};
