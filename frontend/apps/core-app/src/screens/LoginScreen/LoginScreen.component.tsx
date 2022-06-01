import * as React from 'react';
import { Box, Button, Text } from 'native-base';
import { Logo } from 'core-design-system';

type Props = {
  onLogin: () => any;
  isLoading: boolean;
};

export const LoginScreenComponent: React.FC<Props> = ({
  onLogin,
  isLoading,
}) => {
  return (
    <Box height="100%" flexDirection="row">
      <Box
        pl="130px"
        pt="40px"
        justifyContent="flex-start"
        alignItems="flex-start"
        width="700px"
        maxWidth="700px"
        height="100%"
        bg="red"
        zIndex={1}
      >
        <Logo size="65px" pb="55px" />
        <Text variant="heading" level="1" mb="40px">
          Login
        </Text>
        <Text variant="body" level="1" mb="35px">
          Access to your personal account
        </Text>
        <Button
          onPress={onLogin}
          color="primary"
          variant="major"
          w="440px"
          isLoading={isLoading}
          testID="login-button"
        >
          Login
        </Button>
      </Box>

      <Box flexGrow="1" bg="neutral.200" />
    </Box>
  );
};
