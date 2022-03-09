import * as React from 'react';
import { Button, Text, VStack } from 'native-base';
import { RouteProp } from '@react-navigation/native';

import { RootParamList } from '../../navigation/types';

import * as Styles from './RecipientListScreen.styles';

type Props = {
  route: RouteProp<RootParamList, 'RecipientList'>;
  handleItemClick: (id: string) => void;
};

const IDS = ['a76cea53-5a0e-48f3-8917-70a9e74a1b32', 'fake-id']; // TODO remove when actual list available

export const RecipientListScreenComponent: React.FC<Props> = ({
  route,
  handleItemClick,
}) => {
  return (
    <Styles.Container>
      <Text variant="display">{route.name}</Text>
      <VStack space={2} width="sm">
        {IDS.map((id) => (
          <Button variant="major" onPress={() => handleItemClick(id)} key={id}>
            Recipient {id}
          </Button>
        ))}
      </VStack>
    </Styles.Container>
  );
};
