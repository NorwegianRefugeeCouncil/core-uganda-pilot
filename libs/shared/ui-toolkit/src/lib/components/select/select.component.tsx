import { FunctionComponent, MouseEventHandler } from 'react';
import Select, { components } from 'react-select';
import type { Ref } from 'react';
import classNames from 'classnames';

// a thin wrapper over react-select applying our own bootstrap styles

// const customStyles = {
//   clearIndicator: (provided, state) => ({
//     ...provided,
//   }),
//   container: (provided, state) => ({
//     ...provided,
//   }),
//   control: (provided, state) => ({
//     ...provided,
//   }),
//   dropdownIndicator: (provided, state) => ({
//     ...provided,
//   }),
//   group: (provided, state) => ({
//     ...provided,
//   }),
//   groupHeading: (provided, state) => ({
//     ...provided,
//   }),
//   indicatorsContainer: (provided, state) => ({
//     ...provided,
//   }),
//   indicatorSeparator: (provided, state) => ({
//     ...provided,
//   }),
//   input: (provided, state) => ({
//     ...provided,
//   }),
//   loadingIndicator: (provided, state) => ({
//     ...provided,
//   }),
//   loadingMessage: (provided, state) => ({
//     ...provided,
//   }),
//   menu: (provided, state) => ({
//     ...provided,
//   }),
//   menuList: (provided, state) => ({
//     ...provided,
//   }),
//   menuPortal: (provided, state) => ({
//     ...provided,
//   }),
//   multiValue: (provided, state) => ({
//     ...provided,
//   }),
//   multiValueLabel: (provided, state) => ({
//     ...provided,
//   }),
//   multiValueRemove: (provided, state) => ({
//     ...provided,
//   }),
//   noOptionsMessage: (provided, state) => ({
//     ...provided,
//   }),
//   option: (provided, state) => ({
//     ...provided,
//     backgroundColor: state.isSelected ? 'var()' : '$secondary',
//   }),
//   placeholder: (provided, state) => ({
//     ...provided,
//   }),
//   singleValue: (provided, state) => ({
//     ...provided,
//   }),
//   valueContainer: (provided, state) => ({
//     ...provided,
//   }),
// };

type OptionType = { [key: string]: string };
type OptionsType = Array<OptionType>;

type GroupType = {
  [key: string]: any; // group label
  options: OptionsType;
};

type ValueType = OptionType | OptionsType | null | void;

type CommonProps = HTMLElement & {
  clearValue: () => void;
  getStyles: (string, any) => Record<string, unknown>;
  getValue: () => ValueType;
  hasValue: boolean;
  isMulti: boolean;
  options: OptionsType;
  selectOption: (val: OptionType) => void;
  selectProps: any;
  setValue: (ValueType, ActionTypes) => void;
  emotion: any;
};

// passed as the second argument to `onChange`
type ActionTypes =
  | 'clear'
  | 'create-option'
  | 'deselect-option'
  | 'pop-value'
  | 'remove-value'
  | 'select-option'
  | 'set-value';

type PropsWithInnerRef = {
  /** The inner reference. */
  innerRef: Ref<any>;
};

type OptionProps = CommonProps &
  PropsWithInnerRef & {
    data: any;
    id: number;
    index: number;
    isDisabled: boolean;
    isFocused: boolean;
    isSelected: boolean;
    label: string;
    onClick: MouseEventHandler;
    onMouseOver: MouseEventHandler;
    value: any;
  };

const Control: FunctionComponent<CommonProps> = ({ children, ...props }) => {
  const classes = classNames('dropdown', props.className, {});
  return (
    <components.Control {...props} className={classes}>
      {children}
    </components.Control>
  );
};

const Option: FunctionComponent<OptionProps> = ({
  innerRef,
  innerProps,
  ...props
}) => {
  const classes = classNames('list-group-item', props.className, {
    disabled: props.isDisabled,
    active: props.isSelected,
  });
  return (
    <div ref={innerRef} {...innerProps} className={classes}>
      {props.label}
    </div>
  );
};

const SingleSelect: FunctionComponent<CommonProps> = (props) => {
  return <Select {...props} components={{ Control, Option }} />;
};

const MultiSelect: FunctionComponent<CommonProps> = (props) => {
  return <Select {...props} components={{ Control, Option }} isMulti />;
};

export { SingleSelect, MultiSelect };
