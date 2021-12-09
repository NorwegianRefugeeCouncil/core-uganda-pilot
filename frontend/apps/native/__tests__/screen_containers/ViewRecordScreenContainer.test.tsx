import { render, waitFor } from '@testing-library/react-native';
import React from 'react';
import { FormDefinition } from 'core-js-api-client';

import { ViewRecordScreenContainer } from '../../src/components/screen_containers/ViewRecordScreenContainer';

const mockFormId = 'a form id';
const mockRecordId = 'a record id';

const mockRoute = {
  params: {
    formId: mockFormId,
    recordId: mockRecordId,
  },
};

const mockState: any = {
  formsById: {},
};

mockState.formsById[mockFormId] = {
  recordsById: {},
};

const mockValue = {
  fieldId: 'a field id',
  value: 'a value',
};
const mockValues = Array(10).fill(mockValue);

mockState.formsById[mockFormId].recordsById[mockRecordId] = { values: mockValues };

const mockForm: FormDefinition = {
  id: mockFormId,
  code: 'a code',
  databaseId: 'a database id',
  folderId: 'a folder id',
  name: 'a name',
  fields: [],
};

const mockFormData: { response: FormDefinition } = {
  response: mockForm,
};
const mockGetForm = jest.fn().mockResolvedValue(mockFormData);

jest.mock('../../src/utils/clients', () => () => ({
  getForm: mockGetForm,
}));

const mockControl = {};
const mockReset = jest.fn();
jest.mock('react-hook-form', () => {
  const rhf = jest.requireActual('react-hook-form');
  return {
    ...rhf,
    useForm: () => ({
      control: mockControl,
      reset: mockReset,
    }),
  };
});

const mockProps: any = {
  route: mockRoute,
  state: mockState,
};

describe(ViewRecordScreenContainer.name, () => {
  test('renders correctly', async () => {
    const { toJSON } = await waitFor(() => render(<ViewRecordScreenContainer {...mockProps} />));
    expect(toJSON()).toMatchSnapshot();
  });

  test('fetches a form', async () => {
    await waitFor(() => render(<ViewRecordScreenContainer {...mockProps} />));

    const expectedGetFormArg = { id: mockFormId };
    expect(mockGetForm).toHaveBeenCalledWith(expectedGetFormArg);
  });

  test("call react-hook-form's reset method with form data", async () => {
    await waitFor(() => render(<ViewRecordScreenContainer {...mockProps} />));
    await waitFor(() => expect(mockReset).toHaveBeenCalledWith(mockValues));
  });
});
