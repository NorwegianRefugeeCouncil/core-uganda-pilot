import {View} from "react-native";
import React from "react";
import {DataTable, Text} from "react-native-paper";
import {darkTheme} from "../../constants/theme";
import {InputProps} from "./FormControl";
// import {Picker} from "@react-native-picker/picker";

const ReferenceInput: React.FC<InputProps> = (
    {
        fieldDefinition,
        style,
        value,
        onChange,
        onBlur,
        error,
        invalid,
        isTouched,
        isDirty,
        isMultiple,
        isQuantity
    }) => {

    // console.log(isDirty, isTouched, error)
    const [selectedValue, setSelectedValue] = React.useState(value);

    return (
        <View style={style}>
            {fieldDefinition.name && <Text theme={darkTheme}>{fieldDefinition.name}</Text>}
            {fieldDefinition.description &&(
                <Text theme={darkTheme} style={{fontSize: 10}}>
                    {fieldDefinition.description}
                </Text>
            )}
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
            {isTouched && isDirty && error && (
                <Text>
                    {error.message == '' ? 'invalid' : error.message}
                </Text>
            )}
        </View>

    );
};

export default ReferenceInput;
