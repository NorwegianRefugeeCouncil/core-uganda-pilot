import * as React from 'react';
import { Button, Text, VStack } from 'native-base';
import { NativeStackScreenProps } from 'react-native-screens/native-stack';

import { RootStackParamList } from '../../navigators/types';

import * as Styles from './RecipientListScreen.styles';

type Props = NativeStackScreenProps<RootStackParamList, 'RecipientList'>;

export const RecipientListScreenComponent: React.FC<Props> = ({
  route,
  navigation,
}) => {
  return (
    <Styles.Container>
      <Text variant="display">{route.name}</Text>
      <VStack space={2} width="sm">
        <Button
          variant="major"
          onPress={() =>
            navigation.navigate('RecipientProfile', {
              id: '1',
            })
          }
        >
          Recipient 1
        </Button>
        <Button
          variant="major"
          onPress={() =>
            navigation.navigate('RecipientProfile', {
              id: '2',
            })
          }
        >
          Recipient 2
        </Button>
      </VStack>
    </Styles.Container>
  );
};
