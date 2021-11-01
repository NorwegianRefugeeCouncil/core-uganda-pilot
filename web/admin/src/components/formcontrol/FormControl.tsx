import React, {ChangeEvent, FocusEvent, forwardRef, ReactNode} from "react";
import classNames from "classnames";

export interface Option {
    label: string
    value: string
    disabled?: boolean
    placeholder?: string
}

export interface FormControlProps {
    label: string,
    className?: string
    name: string
    children?: ReactNode | undefined
    onChange: (ev: ChangeEvent<HTMLInputElement | HTMLSelectElement>) => void
    onBlur: (ev: FocusEvent<HTMLInputElement | Element>) => void
    options?: Option[]
    multiple?: boolean
    placeholder?: string
}

function FormControlInner(
    props: FormControlProps,
    ref: React.ForwardedRef<HTMLInputElement> | React.ForwardedRef<HTMLSelectElement>,
) {
    const {label, name, className, children, onChange, onBlur, options, placeholder, multiple} = props
    if (!options) {
        return <div className={"form-group mb-2"}>
            <label className={"form-label"}>{label}</label>
            <input
                placeholder={placeholder}
                name={name}
                ref={ref as React.ForwardedRef<HTMLInputElement>}
                type={"text"}
                onChange={onChange}
                onBlur={onBlur}
                className={classNames("form-control form-control-darkula", className)}/>
            {children}
        </div>
    } else {
        return <div className="form-group mb-2 ">
            <label className={"form-label"}>{label}</label>
            <select ref={ref as React.ForwardedRef<HTMLSelectElement>}
                    placeholder={placeholder}
                    defaultValue={""}
                    className={classNames("form-select bg-darkula border-secondary text-light", className)}
                    name={name}
                    onChange={onChange}
                    multiple={multiple}
                    onBlur={onBlur}>
                {options.map(o => (
                    <option placeholder={o.placeholder} key={o.value} disabled={o.disabled}
                            value={o.value}>{o.label}</option>
                ))}
            </select>
            {children}
        </div>

    }

}

export const FormControl = forwardRef(FormControlInner)
