import * as React from 'react';
import { Box, Text } from 'native-base';
import { FormDefinition } from 'core-api-client';

import { formsClient } from '../clients/formsClient';
import { useAPICall } from '../hooks/useAPICall';

type ProviderProps = {
  children: React.ReactNode;
};

const RecipientFormsContext = React.createContext<FormDefinition[]>([]);

export const useRecipientForms = (): FormDefinition[] =>
  React.useContext(RecipientFormsContext);

export const RecipientFormsProvider: React.FC<ProviderProps> = ({
  children,
}) => {
  const [_, state] = useAPICall(
    formsClient.Recipient.getRecipientForms,
    [],
    true,
  );

  if (state.error) {
    return (
      <Box>
        <Text>Error loading recipient forms: {state.error}</Text>
      </Box>
    );
  }

  return (
    <RecipientFormsContext.Provider value={state.data || []}>
      {children}
    </RecipientFormsContext.Provider>
  );
};
