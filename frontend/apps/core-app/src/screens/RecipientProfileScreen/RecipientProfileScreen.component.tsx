import * as React from 'react';
import { Text } from 'native-base';

import * as Styles from './RecipientProfileScreen.styles';

type Props = {};

export const RecipientProfileScreenComponent: React.FC<Props> = () => {
  return (
    <Styles.Container>
      <Text>RecipientList</Text>
    </Styles.Container>
  );
};
