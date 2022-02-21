import * as React from 'react';
import { Text } from 'native-base';

import * as Styles from './RecipientListScreen.styles';

type Props = {};

export const RecipientListScreenComponent: React.FC<Props> = () => {
  return (
    <Styles.Container>
      <Text>RecipientList</Text>
    </Styles.Container>
  );
};
