import * as React from 'react';
import { Text } from 'native-base';
import { RouteProp } from '@react-navigation/native';

import { RootParamList } from '../../navigation/types';

import * as Styles from './RecipientRegistrationScreen.styles';

type Props = {
  route: RouteProp<RootParamList, 'RecipientRegistration'>;
};

export const RecipientRegistrationScreenComponent: React.FC<Props> = ({
  route,
}) => {
  return (
    <Styles.Container>
      <Text variant="display">{route.name}</Text>
    </Styles.Container>
  );
};
