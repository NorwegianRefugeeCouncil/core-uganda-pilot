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
import * as ReactHookFormTransformer from '../../../utils/ReactHookFormTransformer';
import { formsClient } from '../../../clients/formsClient';

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
    defaultValues: ReactHookFormTransformer.toReactHookForm([
      {
        form: subForm,
        record: formsClient.Record.buildDefaultRecord(subForm),
      },
    ]),
  });

  const handleOpenModal = () => setOpen(true);

  const handleCloseModal = () => {
    setOpen(false);
    f.reset();
  };

  const handleAdd = f.handleSubmit((data: any) => {
    onChange([...value, data[field.id]]);
    handleCloseModal();
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
              testID="sub-form-field-input-open-modal-button"
              onPress={handleOpenModal}
              colorScheme="secondary"
              variant="minor"
              w={50}
            >
              Add
            </Button>
          </HStack>
          {value.length === 0 ? (
            <Text testID="sub-form-field-input-empty" variant="body" level="1">
              No entries
            </Text>
          ) : (
            <ScrollView>
              <HStack space={2}>
                {field.fieldType.subForm?.fields.map((subField) => (
                  <VStack key={subField.id} space={2}>
                    <Text variant="heading" level={5}>
                      {subField.name}
                    </Text>
                    {value.map((v, i) => (
                      <Text
                        key={`value-${subField.id}-${i}`}
                        testID="sub-form-field-input-value"
                      >
                        {v[subField.id]}
                      </Text>
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
          <Modal.Header testID="sub-form-field-input-modal-header">
            {field.name}
          </Modal.Header>
          <Modal.Body>
            <FormProvider {...f}>
              <RecordEditor form={subForm} />
            </FormProvider>
          </Modal.Body>
          <Modal.Footer>
            <HStack space={4}>
              <Button
                testID="sub-form-field-input-modal-cancel-button"
                onPress={handleCloseModal}
                colorScheme="secondary"
                variant="minor"
              >
                Cancel
              </Button>
              <Button
                testID="sub-form-field-input-modal-add-button"
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
