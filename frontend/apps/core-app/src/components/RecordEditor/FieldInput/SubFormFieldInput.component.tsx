import * as React from 'react';
import { FormDefinition, FormType } from 'core-api-client';
import {
  FormProvider,
  useController,
  useForm,
  useFormContext,
} from 'react-hook-form';
import {
  Button,
  FormControl,
  HStack,
  Modal,
  ScrollView,
  Text,
  VStack,
} from 'native-base';

import { RecordEditor } from '..';
import { buildDefaultRecord } from '../../../screens/RecipientRegistrationScreen/buildDefaultRecord';
import { buildDefaultFormValues } from '../../../screens/RecipientRegistrationScreen/buildDefaultFormValues';

type Props = {
  form: FormDefinition;
  field: FormDefinition['fields'][number];
};

export const SubFormFieldInput: React.FC<Props> = ({ form, field }) => {
  const [open, setOpen] = React.useState(false);

  const { control } = useFormContext();

  const {
    field: { onChange, value },
    fieldState: { error, invalid },
  } = useController({
    name: `${form.id}.${field.id}`,
    control,
    rules: {}, // TODO Record validation
    defaultValue: [],
  });

  const subForm: FormDefinition = {
    id: field.id,
    name: field.name,
    formType: FormType.DefaultFormType,
    code: '',
    databaseId: form.databaseId,
    folderId: form.folderId,
    fields: field.fieldType.subForm?.fields ?? [],
  };

  const f = useForm({
    defaultValues: buildDefaultFormValues(
      [subForm],
      [buildDefaultRecord(subForm)],
    ),
  });

  const handleOpenModal = () => setOpen(true);

  const handleCloseModal = () => setOpen(false);

  const handleAdd = f.handleSubmit((data: any) => {
    onChange([...value, data[field.id]]);
    handleOpenModal();
  });

  return (
    <>
      <FormControl isInvalid={invalid}>
        <VStack width="100%" space={2}>
          <HStack justifyContent="space-between">
            <VStack>
              <FormControl.Label>
                <Text variant="heading" level={4}>
                  {field.name}
                </Text>
              </FormControl.Label>
              <FormControl.HelperText>
                {field.description}
              </FormControl.HelperText>
            </VStack>
            <Button
              onPress={handleOpenModal}
              colorScheme="secondary"
              variant="minor"
              w={50}
            >
              Add
            </Button>
          </HStack>
          {value.length === 0 ? (
            <Text variant="body" level="1">
              No entries
            </Text>
          ) : (
            <ScrollView>
              <HStack space={2} overflowX="scroll">
                {field.fieldType.subForm?.fields.map((f) => (
                  <VStack key={f.id} space={2}>
                    <Text variant="heading" level={5}>
                      {f.name}
                    </Text>
                    {value.map((v, i) => (
                      <Text key={`value-${f.id}-${i}`}>{v[f.id]}</Text>
                    ))}
                  </VStack>
                ))}
              </HStack>
            </ScrollView>
          )}
          <FormControl.ErrorMessage>{error}</FormControl.ErrorMessage>
        </VStack>
      </FormControl>

      <Modal isOpen={open} onClose={handleCloseModal}>
        <Modal.Content>
          <Modal.Header>{field.name}</Modal.Header>
          <Modal.Body>
            <FormProvider {...f}>
              <RecordEditor form={subForm} />
            </FormProvider>
          </Modal.Body>
          <Modal.Footer>
            <HStack space={4}>
              <Button
                onPress={handleCloseModal}
                colorScheme="secondary"
                variant="minor"
              >
                Cancel
              </Button>
              <Button
                onPress={handleAdd}
                colorScheme="secondary"
                variant="major"
              >
                Add
              </Button>
            </HStack>
          </Modal.Footer>
        </Modal.Content>
      </Modal>
    </>
  );
};
