import {
  ChangeEvent,
  CSSProperties,
  FunctionComponent,
  HTMLAttributes,
  InputHTMLAttributes,
  LabelHTMLAttributes,
  useCallback,
} from 'react';
import { FormLabel, FormLabelProps } from '../form-label/label.component';
import classNames from 'classnames';

export enum InputType {
  Text = 'text',
  Number = 'number',
  Date = 'date',
  DateTime = 'datetimelocal',
}

export interface FormInputProps extends InputHTMLAttributes<HTMLInputElement> {
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

export const FormInput: FunctionComponent<FormInputProps> = (props) => {
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

  const handleOnChange = (ev: ChangeEvent<HTMLInputElement>) => {
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

export const FormCheck: FunctionComponent<HTMLAttributes<HTMLDivElement> & Inliner> = ({ className, inline, ...props }) => {
  const classes = classNames(className, 'form-check', { inline: inline });
  return (
    <div className={classes} {...props}>
      {props.children}
    </div>
  );
};

export const FormSwitch: FunctionComponent<HTMLAttributes<HTMLDivElement>> = ({
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

export const FormCheckInput: FunctionComponent<
  InputHTMLAttributes<HTMLInputElement>
> = (props) => {
  return (
    <input
      {...props}
      className={classNames(props.className, 'form-check-input')}
      type={'checkbox'}
    />
  );
};

export const FormRadioInput: FunctionComponent<InputHTMLAttributes<HTMLInputElement> & { checked?: boolean }> = (props) => {
  return (
    <input
      {...props}
      className={classNames(props.className, 'form-check-input')}
      type={'radio'}
    />
  );
};

export const FormCheckLabel: FunctionComponent<
  LabelHTMLAttributes<HTMLLabelElement>
> = (props) => {
  return (
    <label
      {...props}
      className={classNames(props.className, 'form-check-label')}
    />
  );
};


type ValidFeedbackProps = HTMLAttributes<HTMLDivElement> & { show: boolean };

export const ValidFeedback: FunctionComponent<ValidFeedbackProps> = ({
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


type InvalidFeedbackProps = HTMLAttributes<HTMLDivElement> & { show: boolean };

export const InvalidFeedback: FunctionComponent<InvalidFeedbackProps> = (
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


type FormHelpProps = HTMLAttributes<HTMLDivElement>;

export const FormHelp: FunctionComponent<FormHelpProps> = (
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
  containerStyle?: CSSProperties;
  containerProps?: HTMLAttributes<HTMLDivElement>;
  labelProps?: FormLabelProps;
  descriptionProps?: HTMLAttributes<HTMLDivElement>;
  validFeedback?: string;
  invalidFeedback?: string;
} & FormInputProps;

export const FormControl: FunctionComponent<FormControlProps> = (
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
