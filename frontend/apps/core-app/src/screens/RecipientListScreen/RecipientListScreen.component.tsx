import * as React from 'react';
import { Button, Text, VStack } from 'native-base';
import { RouteProp } from '@react-navigation/native';

import { RootParamList } from '../../navigation/types';

import * as Styles from './RecipientListScreen.styles';

type Props = {
  route: RouteProp<RootParamList, 'RecipientList'>;
  handleItemClick: (id: string) => void;
};

export const RecipientListScreenComponent: React.FC<Props> = ({
  route,
  handleItemClick,
}) => {
  return (
    <Styles.Container>
      <Text variant="display">{route.name}</Text>
      <VStack space={2} width="sm">
        <Button variant="major" onPress={() => handleItemClick('1')}>
          Recipient 1
        </Button>
        <Button variant="major" onPress={() => handleItemClick('2')}>
          Recipient 2
        </Button>
      </VStack>
    </Styles.Container>
  );
};
