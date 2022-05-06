import React from 'react';
import { act, fireEvent, waitFor } from '@testing-library/react-native';

import { render } from '../../../../testUtils/render';
import { RecipientListTableFilterContainer } from '../RecipientListTableFilter.container';

describe('GlobalTableFilterContainer', () => {
  it('should match the snapshot', () => {
    const { toJSON } = render(
      <RecipientListTableFilterContainer
        filter="filter"
        setFilter={jest.fn()}
      />,
    );
    expect(toJSON()).toMatchSnapshot();
  });

  it('should call setFilter on change', async () => {
    const setFilterMock = jest.fn();
    const { getByTestId } = render(
      <RecipientListTableFilterContainer
        filter="filter"
        setFilter={setFilterMock}
      />,
    );
    const input = getByTestId('recipient-list-table-filter');
    expect(input.props.value).toEqual('filter');

    act(() => {
      fireEvent.changeText(input, 'new filter text');
    });
    await waitFor(() => expect(input.props.value).toEqual('new filter text'));
    expect(setFilterMock).toHaveBeenCalledWith('new filter text');
  });
});
