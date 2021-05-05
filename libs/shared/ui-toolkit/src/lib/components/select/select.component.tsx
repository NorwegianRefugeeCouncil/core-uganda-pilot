import React, { FunctionComponent, MouseEventHandler, ReactNode } from 'react';
import Select, { components } from 'react-select';
import type { Ref } from 'react';
import classNames from 'classnames';

// a thin wrapper over react-select applying our own bootstrap styles
/* 
The following components are customisable and switchable:
    ClearIndicator
    Control
    DropdownIndicator
    DownChevron
    CrossIcon
    Group
    GroupHeading
    IndicatorsContainer
    IndicatorSeparator
    Input
    LoadingIndicator
    Menu
    MenuList
    MenuPortal
    LoadingMessage
    NoOptionsMessage
    MultiValue
    MultiValueContainer
    MultiValueLabel
    MultiValueRemove
    Option
    Placeholder
    SelectContainer
    SingleValue
    ValueContainer
 */

const MultiValueLabel: FunctionComponent<MultiValueGenericProps> = ({
  children,
  innerProps,
  ...props
}) => {
  const classes = classNames('badge bg-primary', props.className, {});
  return (
    <div {...innerProps} className={classes}>
      {children}
    </div>
  );
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

const SingleSelect: FunctionComponent<SelectProps> = (props) => {
  return <Select {...props} components={{ Control, Option }} />;
};

const MultiSelect: FunctionComponent<SelectProps> = (props) => {
  return (
    <Select
      {...props}
      components={{ Control, Option, MultiValueLabel }}
      isMulti
    />
  );
};

export { SingleSelect, MultiSelect };

/* TYPE DEFINITIONS */

type SelectProps = {
  autoFocus?: boolean;
  className?: string;
  classNamePrefix?: string;
  isDisabled?: boolean;
  isMulti?: boolean;
  isSearchable?: boolean;
  name?: string;
  onChange?: (
    val:
      | Record<string, unknown>
      | Array<Record<string, unknown>>
      | null
      | undefined,
    opt: {
      action: ActionTypes;
      option?: OptionType;
      removedValue?: Record<string, unknown>;
      name?: string;
    }
  ) => undefined;
  options?: OptionsType;
  placeholder?: ReactNode;
  noOptionsMessage?: () => undefined;
  value?: ValueType;
};

export type MultiValueGenericProps = {
  children?: Node;
  data?: any;
  innerProps?: { className?: string };
  selectProps?: any;
};

type OptionType = { [key: string]: string };
type OptionsType = Array<OptionType>;

type GroupType = {
  [key: string]: any; // group label
  options: OptionsType;
};

type ValueType = OptionType | OptionsType | null | void;

type CommonProps = {
  clearValue?: () => void;
  getStyles?: (string, any) => Record<string, unknown>;
  getValue?: () => ValueType;
  hasValue?: boolean;
  isMulti?: boolean;
  options?: OptionsType;
  selectOption?: (val: OptionType) => void;
  selectProps?: any;
  setValue?: (ValueType, ActionTypes) => void;
  emotion?: any;
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
