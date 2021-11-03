import React from 'react';
import {common, layout} from '../../styles';
import {ScrollView, View} from 'react-native';
import useApiClient from "../../utils/clients";
import {FormDefinition, Record} from "core-js-api-client/lib/types/types";
import {useForm} from "react-hook-form";
import FormControl from "../form/FormControl";

const ViewRecordScreen: React.FC<any> = ({route, state, dispatch}) => {
    const {formId, id} = route.params;

    const [isLoading, setIsLoading] = React.useState(true);
    const [form, setForm] = React.useState<FormDefinition>();
    const [record, setRecord] = React.useState<Record>();

    const client = useApiClient();
    const {control} = useForm();

    React.useEffect(() => {
        client.getForm({id: formId})
            .then((data) => {
                setForm(data.response)
                setIsLoading(false)
                setRecord(state.formsById[formId].recordsById[id])
            })
    }, [formId]);


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
                                    value={record?.values[field.id]}
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
