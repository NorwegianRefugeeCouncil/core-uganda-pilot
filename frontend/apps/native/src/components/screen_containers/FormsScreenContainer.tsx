import { FormDefinition } from 'core-api-client';
import React from 'react';

import { FormsScreenContainerProps } from '../../types/screens';
import client from '../../utils/clients';
import { FormsScreen } from '../screens/FormsScreen';

export const FormsScreenContainer = ({ navigation, route }: FormsScreenContainerProps) => {
  const [forms, setForms] = React.useState<FormDefinition[]>();
  const [isLoading, setIsLoading] = React.useState(true);
  React.useEffect(() => {
    client()
      .listForms({})
      .then((data) => {
        setIsLoading(false);
        setForms(data.response?.items);
      });
  }, [client]);
  return <FormsScreen isLoading={isLoading} forms={forms} navigation={navigation} />;
};
