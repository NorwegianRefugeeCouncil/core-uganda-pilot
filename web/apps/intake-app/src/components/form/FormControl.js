"use strict";
exports.__esModule = true;
var react_native_1 = require("react-native");
var react_1 = require("react");
var TextInput_1 = require("./TextInput");
var Select_1 = require("./Select");
var react_hook_form_1 = require("react-hook-form");
var api_client_1 = require("@core/api-client");
var api_client_2 = require("@core/api-client");
var ReferenceInput_1 = require("./ReferenceInput");
var FormControl = function (_a) {
    var fieldDefinition = _a.fieldDefinition, style = _a.style, control = _a.control, name = _a.name, value = _a.value;
    return (
    // TODO: apply errors to all input types
    <react_native_1.View style={{ margin: 10 }}>
            <react_hook_form_1.Controller name={name} control={control} defaultValue={value} rules={{}} render={function (_a) {
            var _b = _a.field, onChange = _b.onChange, onBlur = _b.onBlur, value = _b.value, ref = _b.ref, fieldState = _a.fieldState, formState = _a.formState;
            var fieldKind = (0, api_client_2.getFieldKind)(fieldDefinition.fieldType);
            switch (fieldKind) {
                case api_client_1.FieldKind.Reference:
                    return (<ReferenceInput_1["default"] fieldDefinition={fieldDefinition} style={style} value={value} onBlur={onBlur} onChange={onChange}/>);
                case api_client_1.FieldKind.Quantity:
                    return (<TextInput_1["default"] fieldDefinition={fieldDefinition} style={style} value={value} onBlur={onBlur} onChange={onChange} isQuantity={true} {...fieldState}/>);
                case api_client_1.FieldKind.MultilineText:
                    return (<TextInput_1["default"] fieldDefinition={fieldDefinition} style={style} value={value} onBlur={onBlur} onChange={onChange} isMultiple={true} {...fieldState}/>);
                case api_client_1.FieldKind.SingleSelect:
                    return (<Select_1["default"] fieldDefinition={fieldDefinition} style={style} value={value} onBlur={onBlur} onChange={onChange}/>);
                default:
                    return (<TextInput_1["default"] fieldDefinition={fieldDefinition} style={style} value={value} onBlur={onBlur} onChange={onChange} {...fieldState}/>);
            }
        }}/>
        </react_native_1.View>);
};
exports["default"] = FormControl;
