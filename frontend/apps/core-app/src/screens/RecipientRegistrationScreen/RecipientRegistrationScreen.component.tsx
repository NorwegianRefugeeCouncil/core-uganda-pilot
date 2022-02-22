import * as React from 'react';
import { Text } from 'native-base';
import { RouteProp, useRoute } from '@react-navigation/native';

import { RootParamList } from '../../navigators/types';

import * as Styles from './RecipientRegistrationScreen.styles';

export const RecipientRegistrationScreenComponent: React.FC = () => {
  const route =
    useRoute<RouteProp<RootParamList, 'RecipientRegistration'>>();

  return (
    <Styles.Container>
      <Text variant="display">{route.name}</Text>
    </Styles.Container>
  );
};
