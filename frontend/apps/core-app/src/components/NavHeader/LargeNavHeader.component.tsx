import * as React from 'react';
import { HStack, Image, Link, Text, VStack } from 'native-base';
import { StackHeaderProps } from '@react-navigation/stack';

export const LargeNavHeaderComponent: React.FC<StackHeaderProps> = ({
  route,
}) => {
  return (
    <VStack width="100%" pt={10} pb={10} maxWidth={1180} ml="auto" mr="auto">
      <HStack width="100%" alignItems="center" justifyContent="space-between">
        <Image
          width={65}
          height={65}
          source={require('../../../assets/png/nrc_logo.png')}
        />
        <HStack>
          <Link href="https://google.com" mr={8}>
            Beneficiary
          </Link>
          <Link href="https://google.com" mr={8}>
            Create beneficiary
          </Link>
          <Link href="https://google.com" mr={8}>
            Settings
          </Link>
          <Link href="https://google.com" mr={8}>
            Log out
          </Link>
        </HStack>
      </HStack>
      <Text variant="display" level="1" mt={12}>
        {route.name}
      </Text>
    </VStack>
  );
};
