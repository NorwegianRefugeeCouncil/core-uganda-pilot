import React from 'react';
import { act, fireEvent, waitFor } from '@testing-library/react-native';
import { FormDefinition, FormType, FormWithRecord } from 'core-api-client';
import { Recipient } from 'core-api-client/src/types/client/Recipient';

import { render } from '../../../testUtils/render';
import { RecipientListScreenContainer } from '../RecipientListScreen.container';
import * as hooks from '../../../hooks/useAPICall';
import { formsClient } from '../../../clients/formsClient';
import configuration from '../../../config';
import { routes } from '../../../constants/routes';
import {
  makeField,
  makeForm,
  makeFormWithRecord,
} from '../../../testUtils/mockData';

const mockNavigate = jest.fn();

jest.mock('@react-navigation/native', () => {
  const actualNav = jest.requireActual('@react-navigation/native');
  return {
    ...actualNav,
    useNavigation: () => ({
      navigate: mockNavigate,
    }),
  };
});

const field0 = makeField(0, false, false, { text: {} });
const field1 = makeField(0, false, false, { text: {} });
const forms = [
  makeForm(0, FormType.RecipientFormType, [field0]),
  makeForm(1, FormType.RecipientFormType, [field1]),
];
const data = [makeFormWithRecord(0), makeFormWithRecord(1)];

describe('RecipientListScreenContainer', () => {
  const useAPICallSpy = jest.spyOn(hooks, 'useAPICall');

  afterEach(() => {
    useAPICallSpy.mockReset();
  });

  describe('should match the snapshot', () => {
    it('data', () => {
      useAPICallSpy.mockImplementation((func, _, __) =>
        func === formsClient.Recipient.getRecipientForms
          ? [
              () => Promise.resolve(),
              { data: forms, loading: false, error: null },
            ]
          : [() => Promise.resolve(), { data, loading: false, error: null }],
      );
      const { toJSON } = render(<RecipientListScreenContainer />);
      expect(toJSON()).toMatchSnapshot();
    });

    it('loading', () => {
      useAPICallSpy.mockImplementation((func, _, __) =>
        func === formsClient.Form.getAncestors
          ? [
              () => Promise.resolve(),
              { data: null, loading: true, error: null },
            ]
          : [
              () => Promise.resolve(),
              { data: null, loading: false, error: null },
            ],
      );
      const { toJSON } = render(<RecipientListScreenContainer />);
      expect(toJSON()).toMatchSnapshot();
    });

    it('error', () => {
      useAPICallSpy.mockImplementation((func, _, __) =>
        func === formsClient.Form.getAncestors
          ? [
              () => Promise.resolve(),
              { data: null, loading: false, error: 'formError' },
            ]
          : [
              () => Promise.resolve(),
              { data: null, loading: false, error: null },
            ],
      );
      const { toJSON } = render(<RecipientListScreenContainer />);
      expect(toJSON()).toMatchSnapshot();
    });
  });

  it('should call the api', () => {
    useAPICallSpy.mockImplementation((func, _, __) =>
      func === formsClient.Form.getAncestors
        ? [
            () => Promise.resolve(),
            { data: forms, loading: false, error: null },
          ]
        : [() => Promise.resolve(), { data, loading: false, error: null }],
    );
    render(<RecipientListScreenContainer />);

    expect(useAPICallSpy).toHaveBeenCalledTimes(6);
    expect(useAPICallSpy).toHaveBeenCalledWith(
      formsClient.Form.getAncestors,
      [configuration.recipient.registrationForm.formId],
      true,
    );
    expect(useAPICallSpy).toHaveBeenCalledWith(
      formsClient.Recipient.list,
      [
        {
          formId: configuration.recipient.registrationForm.formId,
          databaseId: configuration.recipient.registrationForm.databaseId,
        },
      ],
      true,
    );
  });

  it('should call handleItemClick', async () => {
    useAPICallSpy.mockImplementation((func, _, __) =>
      func === formsClient.Form.getAncestors
        ? [
            () => Promise.resolve(),
            { data: forms, loading: false, error: null },
          ]
        : [() => Promise.resolve(), { data, loading: false, error: null }],
    );
    const { findAllByTestId } = render(<RecipientListScreenContainer />);
    expect(useAPICallSpy).toHaveBeenCalledTimes(6);

    const testRows = await findAllByTestId('recipient-list-table-row-0');

    act(() => {
      fireEvent.press(testRows[0]);
    });

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(routes.recipientsProfile.name, {
        id: '0',
      });
    });
  });
});
