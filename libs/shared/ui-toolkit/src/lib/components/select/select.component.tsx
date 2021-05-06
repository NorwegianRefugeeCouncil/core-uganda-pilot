import * as  React from 'react';
import Select, { components, ControlProps, MultiValueProps, OptionProps, Props } from 'react-select';
import classNames from 'classnames';
import { XIcon } from '@nrc.no/ui-toolkit';

// a thin wrapper over react-select applying our own bootstrap styles


const MultiValue: React.FC<MultiValueProps<any>> = (props) => {
  const classes = classNames('bg-primary text-light', props.className, {});
  return (
    <components.MultiValue {...props} className={classes}>
      {props.children}
    </components.MultiValue>
  );
};

const MultiValueLabel: React.FC<MultiValueProps<any>> = ({
                                                           children,
                                                           innerProps,
                                                           ...props
                                                         }) => {
  const classes = classNames('bg-primary', props.className, {});
  return (
    <components.MultiValueLabel {...innerProps} className={classes}>
      {children}
    </components.MultiValueLabel>
  );
};

const MultiValueRemove: React.FC<MultiValueProps<any>> = ({
                                                            innerProps,
                                                            ...props
                                                          }) => {
  const classes = classNames(props.className, '', {});
  return (
    <components.MultiValueRemove {...innerProps}>
      <XIcon />
    </components.MultiValueRemove>
  );
};

const Control: React.FC<ControlProps<any, any>> = ({ children, ...props }) => {
  const classes = classNames('dropdown', props.className, {});
  return (
    <components.Control {...props} className={classes}>
      {children}
    </components.Control>
  );
};

const Option: React.FC<OptionProps<any, any>> = ({
                                                   innerRef,
                                                   innerProps,
                                                   ...props
                                                 }) => {
  const classes = classNames('list-group-item', props.className, {
    disabled: props.isDisabled,
    active: props.isSelected
  });
  return (
    <div ref={innerRef} {...innerProps} className={classes}>
      {props.label}
    </div>
  );
};

type SelectProps = Props

const CustomSelect: React.FC<SelectProps> = (props) => {
  return <Select {...props} components={{ Control, Option, MultiValue, MultiValueLabel, MultiValueRemove }} />;
};

export { Control, Option, MultiValue, MultiValueLabel, MultiValueRemove };
export default CustomSelect;
