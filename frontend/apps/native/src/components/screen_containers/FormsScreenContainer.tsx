import { FormDefinition } from 'core-api-client';
import React from 'react';

import { FormsScreenContainerProps } from '../../types/screens';
import useApiClient from '../../utils/clients';
import { FormsScreen } from '../screens/FormsScreen';

export const FormsScreenContainer = ({
  navigation,
  route,
}: FormsScreenContainerProps) => {
  const [forms, setForms] = React.useState<FormDefinition[]>();
  const [isLoading, setIsLoading] = React.useState(true);
  const apiClient = useApiClient();

  React.useEffect(() => {
    apiClient.Form.list({}).then((data) => {
      setIsLoading(false);
      setForms(data.response?.items);
    });
  }, []);
  return (
    <FormsScreen isLoading={isLoading} forms={forms} navigation={navigation} />
  );
};
