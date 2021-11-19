"use strict";
exports.__esModule = true;
var react_1 = require("react");
var react_native_paper_1 = require("react-native-paper");
var react_native_1 = require("react-native");
var react_native_paper_dates_1 = require("react-native-paper-dates");
var theme_1 = require("../../constants/theme");
var DateInput = function (_a) {
    var fieldDefinition = _a.fieldDefinition, style = _a.style, value = _a.value, onChange = _a.onChange, onBlur = _a.onBlur, error = _a.error, invalid = _a.invalid, isTouched = _a.isTouched, isDirty = _a.isDirty, isMultiple = _a.isMultiple, isQuantity = _a.isQuantity;
    // export function DateModal(props: { date: Date; setDate: React.Dispatch<React.SetStateAction<Date>>; }) {
    // const [date, setDate] = props.date != null && props.setDate != null ? [props.date, props.setDate] : useState(new Date(Date.now()));
    // const [date, setDate] = useState(new Date(Date.now()));
    var _b = (0, react_1.useState)(false), isOpen = _b[0], setIsOpen = _b[1];
    var show = function () { return setIsOpen(true); };
    var hide = function () { return setIsOpen(false); };
    var datePicker;
    if (react_native_1.Platform.OS === "web") {
        datePicker = (<input type="date" value={toDateString(value)} onChange={function (event) { return onChange(toDate(event.target.value)); }}/>);
    }
    else {
        datePicker = (<>
                {fieldDefinition.name && <react_native_paper_1.Text theme={theme_1.darkTheme}>{fieldDefinition.name}</react_native_paper_1.Text>}
                {fieldDefinition.description && (<react_native_paper_1.Text theme={theme_1.darkTheme} style={{ fontSize: 10 }}>
                        {fieldDefinition.description}
                    </react_native_paper_1.Text>)}
                <react_native_paper_1.Button onPress={show}>
                    {toDateString(value)}
                </react_native_paper_1.Button>
                <react_native_paper_dates_1.DatePickerModal locale="en" mode={"single"} onConfirm={function (p) {
                hide();
                onChange(p.date);
            }} onDismiss={hide} visible={isOpen}/>
            </>);
    }
    return (<>
            {datePicker}
        </>);
};
function toDateString(date) {
    console.log('TO DATE STRING', date);
    if (!date) {
        return 'yyyy-mm-dd';
    }
    var y = date.getFullYear();
    var zeroPad = function (n) { return n < 9 ? '0' + (n + 1) : n + 1; };
    var m = zeroPad(date.getMonth());
    var d = date.getDate();
    return y + "-" + m + "-" + d;
}
function toDate(yyyymmdd) {
    var _a = yyyymmdd.split("-").map(function (s) { return s[0] === "0" ? +s[1] : +s; }), yyyy = _a[0], mm = _a[1], dd = _a[2];
    return new Date(yyyy, mm - 1, dd);
}
exports["default"] = DateInput;
