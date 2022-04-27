import * as React from 'react';
import { ScrollView as SV } from 'react-native';
import { Button, HStack, ScrollView, VStack } from 'native-base';
import { useForm, FormProvider, FieldValue } from 'react-hook-form';
import { FormWithRecord } from 'core-api-client';
import { Recipient } from 'core-api-client/src/types/client/Recipient';

import * as ReactHookFormTransformer from '../../utils/ReactHookFormTransformer';
import { RecipientEditor } from '../../components/Recipient/RecipientEditor';
import { RecipientViewer } from '../../components/Recipient/RecipientViewer';

type Props = {
  mode: 'edit' | 'review';
  data: FormWithRecord<Recipient>[];
  onSubmit: (data: FormWithRecord<Recipient>[]) => void;
  onCancel: () => void;
  error: string | null;
  loading: boolean;
};

export const RecipientRegistrationScreenComponent: React.FC<Props> = ({
  mode,
  data,
  onSubmit,
  onCancel,
  error,
  loading,
}) => {
  const scrollRef = React.useRef<SV | null>(null);

  React.useEffect(() => {
    scrollRef.current?.scrollTo({
      y: 0,
      animated: true,
    });
  }, [mode]);

  const f = useForm({
    mode: 'all',
    defaultValues: ReactHookFormTransformer.toReactHookForm(data),
  });

  React.useEffect(() => {
    f.reset(ReactHookFormTransformer.toReactHookForm(data));
  }, [JSON.stringify(data)]);

  const handleSubmit =
    mode === 'edit'
      ? f.handleSubmit((submittedData: FieldValue<any>) => {
          onSubmit(
            ReactHookFormTransformer.fromReactHookForm(data, submittedData),
          );
        })
      : () => onSubmit(data);

  if (
    loading ||
    error ||
    data.length === 0 ||
    Object.keys(f.getValues()).length === 0
  )
    return null;

  return (
    <FormProvider {...f}>
      <ScrollView ref={scrollRef} width="100%" maxWidth="1180px" marginX="auto">
        <VStack space={4}>
          {mode === 'edit' && <RecipientEditor data={data} />}
          {mode === 'review' && <RecipientViewer data={data} />}
          <HStack space={4} justifyContent="flex-end">
            <Button
              key={`cancel-${mode}`}
              testID="recipient-registration-cancel-button"
              onPress={onCancel}
              colorScheme="secondary"
              variant="minor"
            >
              {mode === 'edit' ? 'Cancel' : 'Back'}
            </Button>
            <Button
              key={`save-${mode}`}
              testID="recipient-registration-submit-button"
              onPress={handleSubmit}
              colorScheme="primary"
              variant="major"
            >
              {mode === 'edit' ? 'Review' : 'Save'}
            </Button>
          </HStack>
        </VStack>
      </ScrollView>
    </FormProvider>
  );
};
