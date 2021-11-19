"use strict";
exports.__esModule = true;
var react_native_1 = require("react-native");
var react_1 = require("react");
var react_native_paper_1 = require("react-native-paper");
var theme_1 = require("../../constants/theme");
// import {Picker} from "@react-native-picker/picker";
var ReferenceInput = function (_a) {
    var fieldDefinition = _a.fieldDefinition, style = _a.style, value = _a.value, onChange = _a.onChange, onBlur = _a.onBlur, error = _a.error, invalid = _a.invalid, isTouched = _a.isTouched, isDirty = _a.isDirty, isMultiple = _a.isMultiple, isQuantity = _a.isQuantity;
    var _b = react_1["default"].useState(value), selectedValue = _b[0], setSelectedValue = _b[1];
    return (<react_native_1.View style={style}>
            {fieldDefinition.name && <react_native_paper_1.Text theme={theme_1.darkTheme}>{fieldDefinition.name}</react_native_paper_1.Text>}
            {fieldDefinition.description && (<react_native_paper_1.Text theme={theme_1.darkTheme} style={{ fontSize: 10 }}>
                    {fieldDefinition.description}
                </react_native_paper_1.Text>)}
            {/*<DataTable>*/}
            {/*    <DataTable.Header>*/}
            {/*        <DataTable.Title>Dessert</DataTable.Title>*/}
            {/*        <DataTable.Title numeric>Calories</DataTable.Title>*/}
            {/*        <DataTable.Title numeric>Fat</DataTable.Title>*/}
            {/*    </DataTable.Header>*/}

            {/*    <DataTable.Row>*/}
            {/*        <DataTable.Cell>Frozen yogurt</DataTable.Cell>*/}
            {/*        <DataTable.Cell numeric>159</DataTable.Cell>*/}
            {/*        <DataTable.Cell numeric>6.0</DataTable.Cell>*/}
            {/*    </DataTable.Row>*/}

            {/*    <DataTable.Row>*/}
            {/*        <DataTable.Cell>Ice cream sandwich</DataTable.Cell>*/}
            {/*        <DataTable.Cell numeric>237</DataTable.Cell>*/}
            {/*        <DataTable.Cell numeric>8.0</DataTable.Cell>*/}
            {/*    </DataTable.Row>*/}

            {/*    <DataTable.Pagination*/}
            {/*        page={page}*/}
            {/*        numberOfPages={3}*/}
            {/*        onPageChange={(page) => setPage(page)}*/}
            {/*        label="1-2 of 6"*/}
            {/*        optionsPerPage={optionsPerPage}*/}
            {/*        itemsPerPage={itemsPerPage}*/}
            {/*        setItemsPerPage={setItemsPerPage}*/}
            {/*        showFastPagination*/}
            {/*        optionsLabel={'Rows per page'}*/}
            {/*    />*/}
            {/*</DataTable>*/}
            {isTouched && isDirty && error && (<react_native_paper_1.Text>
                    {error.message == '' ? 'invalid' : error.message}
                </react_native_paper_1.Text>)}
        </react_native_1.View>);
};
exports["default"] = ReferenceInput;
