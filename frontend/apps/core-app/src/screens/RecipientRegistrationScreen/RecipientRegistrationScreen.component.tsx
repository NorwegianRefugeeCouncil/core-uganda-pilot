import * as React from 'react';
import { Button, HStack, ScrollView, VStack } from 'native-base';
import { useForm, FormProvider } from 'react-hook-form';
import { FormDefinition, Record } from 'core-api-client';
import { Accordion } from 'core-design-system';

import { RecordEditor } from '../../components/RecordEditor';

import { buildDefaultFormValues } from './buildDefaultFormValues';

type Props = {
  forms: FormDefinition[];
  records: Record[];
  onSubmit: (data: any) => void;
  onCancel: () => void;
};

export const RecipientRegistrationScreenComponent: React.FC<Props> = ({
  forms,
  records,
  onSubmit,
  onCancel,
}) => {
  const f = useForm({ defaultValues: buildDefaultFormValues(forms, records) });

  React.useEffect(() => {
    f.reset();
  }, [JSON.stringify(forms)]);

  return (
    <FormProvider {...f}>
      <ScrollView width="100%" maxWidth="1180px" marginX="auto">
        <VStack space={4}>
          {forms.map((form) => (
            <Accordion key={form.id} header={form.name} defaultOpen>
              <RecordEditor form={form} />
            </Accordion>
          ))}
          <HStack space={4} justifyContent="flex-end">
            <Button onPress={onCancel} colorScheme="secondary" variant="minor">
              Cancel
            </Button>
            <Button
              onPress={f.handleSubmit(onSubmit)}
              colorScheme="primary"
              variant="major"
            >
              Review
            </Button>
          </HStack>
        </VStack>
      </ScrollView>
    </FormProvider>
  );
};
