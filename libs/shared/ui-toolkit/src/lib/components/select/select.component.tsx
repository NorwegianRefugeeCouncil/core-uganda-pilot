import React, { FunctionComponent } from 'react';
import Select, {
  components,
  ControlProps,
  MultiValueProps,
  OptionProps,
  Props
} from 'react-select';
import classNames from 'classnames';
import { CloseButton, XIcon } from '@nrc.no/ui-toolkit';
import StateManager from 'react-select';

// a thin wrapper over react-select applying our own bootstrap styles


const MultiValue: FunctionComponent<MultiValueProps<any>> = (props) => {
  const classes = classNames('bg-primary text-light', props.className, {})
  return (
    <components.MultiValue {...props} className={classes}>
        {props.children}
    </components.MultiValue>
  )
}

const MultiValueLabel: FunctionComponent<MultiValueProps<any>> = ({
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

const MultiValueRemove: FunctionComponent<MultiValueProps<any>> = ({
  innerProps,
  ...props
}) => {
  const classes = classNames(props.className, '', {});
  return (
    <components.MultiValueRemove {...innerProps}>
      <XIcon />
    </components.MultiValueRemove>
  )
}

const Control: FunctionComponent<ControlProps<any, any>> = ({ children, ...props }) => {
  const classes = classNames('dropdown', props.className, {});
  return (
    <components.Control {...props} className={classes}>
      {children}
    </components.Control>
  );
};

const Option: FunctionComponent<OptionProps<any, any>> = ({
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

type SelectProps = Props

const CustomSelect: FunctionComponent<SelectProps> = (props) => {
  return <Select {...props} components={{ Control, Option, MultiValue, MultiValueLabel, MultiValueRemove }} />
}

export { Control, Option, MultiValue, MultiValueLabel, MultiValueRemove, }
export default CustomSelect
