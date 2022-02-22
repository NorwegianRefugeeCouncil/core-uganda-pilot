import { NavigationProp, useNavigation } from '@react-navigation/native';
import React from 'react';

import { render } from '../../../testUtils/render';
import { RecipientListScreenComponent } from '../RecipientListScreen.component';
import { RootParamList } from '../../../navigators/types';

it('should match the snapshot', () => {
  const Wrapper: React.FC = () => {
    const navigation = useNavigation<NavigationProp<RootParamList>>();
    return (
      <RecipientListScreenComponent
        route={{ name: 'RecipientList', params: {}, key: 'key' }}
        navigation={navigation}
      />
    );
  };

  const { toJSON } = render(<Wrapper />);
  expect(toJSON()).toMatchSnapshot();
});
