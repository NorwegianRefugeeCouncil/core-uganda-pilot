import { FormDefinition } from 'core-js-api-client/lib/types/types';
import React from 'react';
import {Control, FieldValues} from 'react-hook-form';
import { ScrollView, View } from 'react-native';

import { common, layout } from '../../styles';
import { useApiClient } from '../../utils/useApiClient';
import FormControl from '../form/FormControl';

export type ViewRecordScreenProps = {
    isLoading: boolean;
    form?: FormDefinition;
    control: Control<FieldValues, Object>
}

export const ViewRecordScreen = ({isLoading, form, control}: ViewRecordScreenProps) => {

    return (
        <View style={[layout.container, layout.body, common.darkBackground]}>
            <ScrollView>
                {!isLoading && (
                    <View>
                        {form?.fields.map(field => {
                            return (
                                <FormControl
                                    key={field.code}
                                    fieldDefinition={field}
                                    style={{ width: '100%' }}
                                    control={control}
                                    name={field.id}
                                />
                            );
                        })}
                    </View>
                )}
            </ScrollView>
        </View>
    );
};
