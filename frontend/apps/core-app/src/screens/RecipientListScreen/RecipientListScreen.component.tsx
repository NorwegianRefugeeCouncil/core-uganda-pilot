import * as React from 'react';
import { Button, Text, VStack } from 'native-base';
import {
  NavigationProp,
  RouteProp,
  useNavigation,
  useRoute,
} from '@react-navigation/native';

import { RootParamList } from '../../navigators/types';

import * as Styles from './RecipientListScreen.styles';

export const RecipientListScreenComponent: React.FC = () => {
  const route = useRoute<RouteProp<RootParamList, 'RecipientList'>>();
  const navigation = useNavigation<NavigationProp<RootParamList>>();

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
