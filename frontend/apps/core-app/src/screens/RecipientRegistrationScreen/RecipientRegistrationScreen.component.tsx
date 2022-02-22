import * as React from 'react';
import { Text } from 'native-base';
import { NativeStackScreenProps } from 'react-native-screens/native-stack';

import { RootStackParamList } from '../../navigators/types';

import * as Styles from './RecipientRegistrationScreen.styles';

type Props = NativeStackScreenProps<
  RootStackParamList,
  'RecipientRegistration'
>;

export const RecipientRegistrationScreenComponent: React.FC<Props> = ({
  route,
}) => {
  return (
    <Styles.Container>
      <Text variant="display">{route.name}</Text>
    </Styles.Container>
  );
};
