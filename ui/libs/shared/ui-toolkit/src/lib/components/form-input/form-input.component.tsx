import * as React from 'react';
import { FormLabel, FormLabelProps } from '../form-label/label.component';
import classNames from 'classnames';

export enum InputType {
  Text = 'text',
  Number = 'number',
  Date = 'date',
  DateTime = 'datetimelocal',
}

export interface FormInputProps extends React.ComponentPropsWithRef<'input'> {
  plaintext?: boolean;
  colorInput?: boolean;
  onValueChanged?: (value: any) => void;
}

const transformInputValue = (props: FormInputProps, value: any) => {
  if (props.type === 'number') {
    if (value) {
      value = +value;
    }
    if (!isNaN(value)) {
      if (typeof +props.min === 'number' && !isNaN(+props.min)) {
        if (value < props.min) {
          value = +props.min;
        }
      }
      if (typeof +props.max === 'number' && !isNaN(+props.max)) {
        if (value > props.max) {
          value = +props.max;
        }
      }
      if (typeof +props.step === 'number' && !isNaN(+props.step)) {
        if (+props.step !== 0) {
          value = Math.round(value / +props.step) * +props.step;
        }
      }
    }
  }
  if (props.type === 'text') {
    if (value) {
      value = '' + value;
    }
  }
  return value;
};

export const FormInput: React.FC<FormInputProps> = (props) => {
  const {
    plaintext,
    colorInput,
    className,
    onValueChanged,
    ...otherProps
  } = props;
  const classes: string[] = [];
  if (plaintext) {
    classes.push('form-control-plaintext');
  } else {
    classes.push('form-control');
    if (colorInput) {
      classes.push('form-control-color');
    }
  }

  const handleOnChange = (ev: React.ChangeEvent<HTMLInputElement>) => {
    const value = transformInputValue(props, ev?.target?.value);
    if (onValueChanged) {
      onValueChanged(value);
    }
    if (props.onChange) {
      props.onChange(ev);
    }
  };

  return (
    <input
      {...otherProps}
      onChange={handleOnChange}
      className={classNames(props.className, ...classes)}
    />
  );
};

export interface Inliner {
  inline?: boolean;
}

export const FormCheck: React.FC<React.ComponentPropsWithRef<'div'> & Inliner> = ({ className, inline, ...props }) => {
  const classes = classNames(className, 'form-check', { inline: inline });
  return (
    <div className={classes} {...props}>
      {props.children}
    </div>
  );
};

export const FormSwitch: React.FC<React.ComponentPropsWithRef<'div'>> = ({
  className,
  ...props
}) => {
  return (
    <div
      className={classNames(className, 'form-check', 'form-switch')}
      {...props}
    >
      {props.children}
    </div>
  );
};

export const FormCheckInput: React.FC<React.ComponentPropsWithRef<'input'>
> = (props) => {
  return (
    <input
      {...props}
      className={classNames(props.className, 'form-check-input')}
      type={'checkbox'}
    />
  );
};

export const FormRadioInput: React.FC<React.ComponentPropsWithRef<'input'> & { checked?: boolean }> = (props) => {
  return (
    <input
      {...props}
      className={classNames(props.className, 'form-check-input')}
      type={'radio'}
    />
  );
};

export const FormCheckLabel: React.FC<
  React.ComponentPropsWithRef<'label'>
> = (props) => {
  return (
    <label
      {...props}
      className={classNames(props.className, 'form-check-label')}
    />
  );
};


type ValidFeedbackProps = React.ComponentPropsWithRef<'div'> & { show: boolean };

export const ValidFeedback: React.FC<ValidFeedbackProps> = ({
  className,
  children,
  show,
  ...props
}) => {
  return (
    <div
      {...props}
      className={classNames(className, 'valid-feedback')}
      style={{ display: show ? 'block' : 'none' }}
    >
      {children}
    </div>
  );
};


type InvalidFeedbackProps = React.ComponentPropsWithRef<'div'> & { show: boolean };

export const InvalidFeedback: React.FC<InvalidFeedbackProps> = (
  { className, children, show, ...props },
  ref
) => {
  return (
    <div
      {...props}
      ref={ref}
      className={classNames(className, 'invalid-feedback')}
      style={{ display: show ? 'block' : 'none' }}
    >
      {children}
    </div>
  );
};


type FormHelpProps = React.ComponentPropsWithRef<'div'>;

export const FormHelp: React.FC<FormHelpProps> = (
  { className, children, ...props },
  ref
) => {
  return (
    <small {...props} className={classNames(className, 'text-muted')}>
      {children}
    </small>
  );
};

type FormControlProps = {
  label?: string;
  description?: string;
  containerClassName?: string;
  containerStyle?: React.CSSProperties;
  containerProps?: React.ComponentPropsWithRef<'div'>;
  labelProps?: FormLabelProps;
  descriptionProps?: React.ComponentPropsWithRef<'div'>;
  validFeedback?: string;
  invalidFeedback?: string;
} & FormInputProps;

export const FormControl: React.FC<FormControlProps> = (
  {
    label,
    description,
    containerProps,
    labelProps,
    descriptionProps,
    ...props
  },
  ref
) => {
  return (
    <div {...containerProps}>
      {label && <FormLabel {...labelProps}>{label}</FormLabel>}
      <FormInput {...props} />
      {description && <FormHelp {...descriptionProps}>{description}</FormHelp>}
    </div>
  );
};
