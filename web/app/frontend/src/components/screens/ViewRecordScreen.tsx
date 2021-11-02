import React from 'react';
import {FAB, Snackbar, Switch, Title} from 'react-native-paper';
import {common, layout} from '../../styles';
import routes from '../../constants/routes';
import {Button, FlatList, Platform, ScrollView, Text, TouchableOpacity, View} from 'react-native';
import useApiClient from "../../utils/clients";
import {FormDefinition} from "core-js-api-client/lib/types/types";
import {useForm} from "react-hook-form";
import {getEncryptionKey} from "../../utils/getEncryptionKey";
import * as SecureStore from "expo-secure-store";
import CryptoJS from "react-native-crypto-js";
import AsyncStorage from "@react-native-async-storage/async-storage";
import _ from "lodash";
import * as Network from "expo-network";
import FormControl from "../form/FormControl";

const ViewRecordScreen: React.FC<any> = ({navigation, route}) => {
    const {formId, id} = route.params;

    const [record, setRecord] = React.useState<any[]>();
    const [isLoading, setIsLoading] = React.useState(true);
    const [form, setForm] = React.useState<FormDefinition>();

    const client = useApiClient();
    const {control} = useForm();

    React.useEffect(() => {
        client.getForm({id: formId})
            .then((data) => {
                // setRecord(data.response?.fields)
                setForm(data.response)
                setIsLoading(false)
            })
        // client.get
    }, []);


    return (
        <View style={[layout.container, layout.body, common.darkBackground]}>
            <ScrollView>
                {!isLoading && (
                    <View>
                        {form?.fields.map((field) => {
                            return (
                                <FormControl
                                    key={field.code}
                                    fieldDefinition={field}
                                    style={{width: '100%'}}
                                    control={control}
                                    name={field.id}
                                />
                            )
                        })}
                    </View>
                )}
            </ScrollView>
        </View>
    );
};

export default ViewRecordScreen;
