import React, { ChangeEvent, FocusEvent, forwardRef, ReactNode, useState } from 'react';
import classNames from 'classnames';

export interface Option {
  label: string;
  value: string;
  disabled?: boolean;
  placeholder?: string;
}

export interface FormControlProps {
  label: string;
  className?: string;
  name?: string;
  children?: ReactNode | undefined;
  onChange?: (ev: ChangeEvent<HTMLInputElement | HTMLSelectElement>) => void;
  onBlur?: (ev: FocusEvent<HTMLInputElement | Element>) => void;
  options?: Option[];
  multiple?: boolean;
  placeholder?: string;
  sensitive?: boolean;
  value?: any;
  readOnly?: boolean;
  allowCopy?: boolean;
}

export const FormControl = forwardRef<HTMLInputElement | HTMLSelectElement, FormControlProps>((props, ref) => {
  const {
    label,
    name,
    className,
    children,
    onChange,
    onBlur,
    options,
    placeholder,
    multiple,
    sensitive,
    value,
    readOnly,
    allowCopy,
  } = props;

  const [reveal, setReveal] = useState(false);

  function copyTextToClipboard(text: string) {
    if ('clipboard' in navigator) {
      return navigator.clipboard.writeText(text);
    }
    return document.execCommand('copy', true, text);
  }

  if (!options) {
    return (
      <div className="form-group pb-3 pt-2 border-secondary">
        <label className="form-label fw-bold">{label}</label>
        <div className="input-group">
          <input
            placeholder={placeholder}
            name={name}
            ref={ref as React.ForwardedRef<HTMLInputElement>}
            type={sensitive && !reveal ? 'password' : 'text'}
            onChange={onChange}
            onBlur={onBlur}
            value={value}
            readOnly={readOnly}
            className={classNames('form-control form-control-darkula', className)}
          />

          {allowCopy && (
            <button
              type="button"
              onClick={() => {
                copyTextToClipboard(value);
              }}
              className="btn btn-outline-secondary"
              title="Copy value"
            >
              <i className="bi bi-clipboard" />
            </button>
          )}

          {sensitive && (
            <button
              onClick={() => setReveal(!reveal)}
              type="button"
              className="btn btn-outline-secondary"
              title={reveal ? 'Hide' : 'Show'}
            >
              {reveal && <i className="bi bi-eye-slash" />}
              {!reveal && <i className="bi bi-eye" />}
            </button>
          )}
        </div>
        {children}
      </div>
    );
  }
  return (
    <div className="form-group pb-3 pt-2">
      <label className="form-label fw-bold">{label}</label>
      <select
        ref={ref as React.ForwardedRef<HTMLSelectElement>}
        placeholder={placeholder}
        defaultValue=""
        className={classNames('form-select bg-darkula border-secondary text-light', className)}
        name={name}
        onChange={onChange}
        multiple={multiple}
        onBlur={onBlur}
      >
        {options.map((o) => (
          <option placeholder={o.placeholder} key={o.value} disabled={o.disabled} value={o.value}>
            {o.label}
          </option>
        ))}
      </select>
      {children}
    </div>
  );
});
