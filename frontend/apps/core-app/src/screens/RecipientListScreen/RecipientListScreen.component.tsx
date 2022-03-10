import * as React from 'react';
import { Button, Text, VStack } from 'native-base';
import { RouteProp } from '@react-navigation/native';

import { RootParamList } from '../../navigation/types';

import * as Styles from './RecipientListScreen.styles';

type Props = {
  route: RouteProp<RootParamList, 'RecipientList'>;
  handleItemClick: (id: string) => void;
};

const IDS = [
  '8090092f-c983-4ff4-8599-214429218eb0',
  '30e683fa-2dd7-479f-98a3-9477d0079383',
  'fake-id',
]; // TODO remove when actual list available

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
