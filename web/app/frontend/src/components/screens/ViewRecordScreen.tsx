import React from 'react';
import {common, layout} from '../../styles';
import {ScrollView, View} from 'react-native';
import {useApiClient} from "../../utils/useApiClient";
import {FormDefinition} from "core-js-api-client/lib/types/types";
import {useForm} from "react-hook-form";
import FormControl from "../form/FormControl";

const ViewRecordScreen: React.FC<any> = ({route, state}) => {
    const {formId, recordId} = route.params;

    const [isLoading, setIsLoading] = React.useState(true);
    const [form, setForm] = React.useState<FormDefinition>();

    const client = useApiClient();
    const {control, reset} = useForm();

    React.useEffect(() => {
        client.getForm({id: formId})
            .then((data) => {
                setForm(data.response)
            })
    }, [formId]);

    React.useEffect(() => {
        if (form) {
            reset(state.formsById[formId].recordsById[recordId].values)
            setIsLoading(false)
        }
    }, [form])

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
