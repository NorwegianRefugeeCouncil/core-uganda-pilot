import * as React from 'react';
import { Button } from 'native-base';

import * as Styles from './LoginScreen.styles';

type Props = {
  onLogin: () => any;
};

export const LoginScreenComponent: React.FC<Props> = ({ onLogin }) => {
  return (
    <Styles.Container>
      <Button onPress={onLogin} color="primary" variant="major">
        Login
      </Button>
    </Styles.Container>
  );
};
