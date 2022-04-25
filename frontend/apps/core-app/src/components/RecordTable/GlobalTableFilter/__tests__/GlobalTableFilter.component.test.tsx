import React from 'react';
import { fireEvent, waitFor } from '@testing-library/react-native';

import { render } from '../../../../testUtils/render';
import { GlobalTableFilterComponent } from '../GlobalTableFilter.component';

describe('GlobalTableFilterComponent', () => {
  it('should match the snapshot', () => {
    const { toJSON } = render(
      <GlobalTableFilterComponent handleChange={jest.fn()} value="test" />,
    );
    expect(toJSON()).toMatchSnapshot();
  });

  it('should call onChange handler', async () => {
    const setGlobalFilterSpy = jest.fn();

    const { getByPlaceholderText } = render(
      <GlobalTableFilterComponent
        handleChange={setGlobalFilterSpy}
        value="test"
      />,
    );

    const input = getByPlaceholderText('Search');
    fireEvent.changeText(input, 'a new value');

    await waitFor(() => {
      expect(setGlobalFilterSpy).toHaveBeenCalledWith('a new value');
    });
  });
});
