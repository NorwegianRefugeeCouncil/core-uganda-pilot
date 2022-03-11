import * as React from 'react';
import { Button, HStack } from 'native-base';
import { useForm, FormProvider } from 'react-hook-form';
import { FormDefinition, Record } from 'core-api-client';
import { Accordion } from 'core-design-system';

import { RecordEditor } from '../../components/RecordEditor';

import { buildDefaultFormValues } from './buildDefaultFormValues';

type Props = {
  forms: FormDefinition[];
  records: Record[];
  onSubmit: () => void;
  onCancel: () => void;
};

export const RecipientRegistrationScreenComponent: React.FC<Props> = ({
  forms,
  records,
  onSubmit,
  onCancel,
}) => {
  const f = useForm({ defaultValues: buildDefaultFormValues(forms, records) });

  return (
    <FormProvider {...f}>
      <Styles.Container>
        {forms.map((form, i) => (
          <Accordion key={form.id} header={form.name}>
            <RecordEditor form={form} record={records[i]} onChange={() => {}} />
          </Accordion>
        ))}
        <HStack>
          <Button onPress={onCancel} />
          <Button onPress={onSubmit} />
        </HStack>
      </Styles.Container>
    </FormProvider>
  );
};
