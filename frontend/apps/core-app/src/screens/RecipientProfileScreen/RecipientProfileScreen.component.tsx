import * as React from 'react';
import { Text } from 'native-base';
import { RouteProp, useRoute } from '@react-navigation/native';

import { RootStackParamList } from '../../navigators/types';

import * as Styles from './RecipientProfileScreen.styles';

export const RecipientProfileScreenComponent: React.FC = () => {
  const route = useRoute<RouteProp<RootStackParamList, 'RecipientProfile'>>();
  return (
    <Styles.Container>
      <Text variant="display">
        {route.name}: {route.params.id}
      </Text>
    </Styles.Container>
  );
};
