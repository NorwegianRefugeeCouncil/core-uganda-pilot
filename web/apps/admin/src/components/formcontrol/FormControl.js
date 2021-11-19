"use strict";
exports.__esModule = true;
exports.FormControl = void 0;
var react_1 = require("react");
var classnames_1 = require("classnames");
exports.FormControl = (0, react_1.forwardRef)(function (props, ref) {
    var label = props.label, name = props.name, className = props.className, children = props.children, onChange = props.onChange, onBlur = props.onBlur, options = props.options, placeholder = props.placeholder, multiple = props.multiple, sensitive = props.sensitive, value = props.value, readOnly = props.readOnly, allowCopy = props.allowCopy;
    var _a = (0, react_1.useState)(false), reveal = _a[0], setReveal = _a[1];
    function copyTextToClipboard(text) {
        if ('clipboard' in navigator) {
            return navigator.clipboard.writeText(text);
        }
        else {
            return document.execCommand('copy', true, text);
        }
    }
    if (!options) {
        return <div className={"form-group pb-3 pt-2 border-secondary"}>
            <label className={"form-label fw-bold"}>{label}</label>
            <div className={"input-group"}>
                <input placeholder={placeholder} name={name} ref={ref} type={sensitive && !reveal ? "password" : "text"} onChange={onChange} onBlur={onBlur} value={value} readOnly={readOnly} className={(0, classnames_1["default"])("form-control form-control-darkula", className)}/>

                {allowCopy && <button type={"button"} onClick={function () {
                    copyTextToClipboard(value);
                }} className={"btn btn-outline-secondary"} title={"Copy value"}>
                    <i className={"bi bi-clipboard"}/>
                </button>}

                {sensitive && (<button onClick={function () { return setReveal(!reveal); }} type={"button"} className={"btn btn-outline-secondary"} title={reveal ? "Hide" : "Show"}>
                        {reveal && <i className={"bi bi-eye-slash"}/>}
                        {!reveal && <i className={"bi bi-eye"}/>}
                    </button>)}

            </div>
            {children}
        </div>;
    }
    else {
        return <div className="form-group pb-3 pt-2">
            <label className={"form-label fw-bold"}>{label}</label>
            <select ref={ref} placeholder={placeholder} defaultValue={""} className={(0, classnames_1["default"])("form-select bg-darkula border-secondary text-light", className)} name={name} onChange={onChange} multiple={multiple} onBlur={onBlur}>
                {options.map(function (o) { return (<option placeholder={o.placeholder} key={o.value} disabled={o.disabled} value={o.value}>{o.label}</option>); })}
            </select>
            {children}
        </div>;
    }
});
