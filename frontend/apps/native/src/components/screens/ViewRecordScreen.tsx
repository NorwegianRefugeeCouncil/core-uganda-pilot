import React from 'react';
import { ScrollView, Text, View } from 'react-native';
import { FormDefinition } from 'core-api-client';
import { Control } from 'react-hook-form';

import { common, layout } from '../../styles';
import FormControl from '../form/FormControl';

export interface ViewRecordScreenProps {
  isLoading: boolean;
  form?: FormDefinition;
  control: Control;
}

const ViewRecordScreen = ({ isLoading, form, control }: ViewRecordScreenProps) => {
  return (
    <View style={[layout.container, layout.body, common.darkBackground]}>
      <ScrollView>
        {isLoading ? (
          <View>
            <Text>Loading...</Text>
          </View>
        ) : (
          <View>
            {form?.fields.map((field) => {
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

export default ViewRecordScreen;
