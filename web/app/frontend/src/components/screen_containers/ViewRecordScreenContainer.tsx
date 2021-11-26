import React from 'react';
import {ViewRecordScreenContainerProps} from '../../types/screens';
import {FormDefinition} from 'core-js-api-client/lib/types/types';
import {useApiClient} from '../../utils/useApiClient';
import {useForm} from 'react-hook-form';
import {ViewRecordScreen} from '../screens/ViewRecordScreen';

export const ViewRecordScreenContainer = ({
                                              route,
                                              state,
                                          }: ViewRecordScreenContainerProps) => {
    const {formId, recordId} = route.params;

    const [isLoading, setIsLoading] = React.useState(true);
    const [form, setForm] = React.useState<FormDefinition>();

    const client = useApiClient();
    const {control, reset} = useForm();

    React.useEffect(() => {
        client.getForm({id: formId}).then(data => {
            setForm(data.response);
        });
    }, [formId]);

    React.useEffect(() => {
        if (form) {
            reset(state.formsById[formId].recordsById[recordId].values);
            setIsLoading(false);
        }
    }, [form]);

    return <ViewRecordScreen isLoading={isLoading} control={control}/>;
};
