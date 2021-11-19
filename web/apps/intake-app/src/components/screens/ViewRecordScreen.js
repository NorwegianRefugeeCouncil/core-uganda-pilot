"use strict";
exports.__esModule = true;
var react_1 = require("react");
var styles_1 = require("../../styles");
var react_native_1 = require("react-native");
var clients_1 = require("../../utils/clients");
var react_hook_form_1 = require("react-hook-form");
var FormControl_1 = require("../form/FormControl");
var ViewRecordScreen = function (_a) {
    var route = _a.route, state = _a.state;
    var _b = route.params, formId = _b.formId, recordId = _b.recordId;
    var _c = react_1["default"].useState(true), isLoading = _c[0], setIsLoading = _c[1];
    var _d = react_1["default"].useState(), form = _d[0], setForm = _d[1];
    var client = (0, clients_1["default"])();
    var _e = (0, react_hook_form_1.useForm)(), control = _e.control, reset = _e.reset;
    react_1["default"].useEffect(function () {
        client.getForm({ id: formId })
            .then(function (data) {
            setForm(data.response);
        });
    }, [formId]);
    react_1["default"].useEffect(function () {
        if (form) {
            reset(state.formsById[formId].recordsById[recordId].values);
            setIsLoading(false);
        }
    }, [form]);
    return (<react_native_1.View style={[styles_1.layout.container, styles_1.layout.body, styles_1.common.darkBackground]}>
            <react_native_1.ScrollView>
                {!isLoading && (<react_native_1.View>
                        {form === null || form === void 0 ? void 0 : form.fields.map(function (field) {
                return (<FormControl_1["default"] key={field.code} fieldDefinition={field} style={{ width: '100%' }} control={control} name={field.id}/>);
            })}
                    </react_native_1.View>)}
            </react_native_1.ScrollView>
        </react_native_1.View>);
};
exports["default"] = ViewRecordScreen;
