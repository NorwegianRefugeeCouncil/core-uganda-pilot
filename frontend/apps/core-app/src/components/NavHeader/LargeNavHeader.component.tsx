import * as React from 'react';
import { HStack, Image, Text, VStack, Button } from 'native-base';
import { StackHeaderProps } from '@react-navigation/stack';
import { getFocusedRouteNameFromRoute } from '@react-navigation/native';

import { Link } from '../Link';
import { routes } from '../../constants/routes';

import { logout } from './logout';

const getHeaderTitle = (route: StackHeaderProps['route']): string => {
  const routeName =
    getFocusedRouteNameFromRoute(route) ?? routes.recipientsRoot.name;

  return routes[routeName as keyof typeof routes].title;
};

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
          <Link
            to={{
              screen: routes.recipientsRoot.name,
              params: { screen: routes.recipientsList.name },
            }}
            mr={8}
          >
            Beneficiaries
          </Link>
          <Link
            to={{
              screen: routes.recipientsRoot.name,
              params: { screen: routes.recipientsRegistration.name },
            }}
            mr={8}
          >
            Register beneficiary
          </Link>
          <Link
            to={{
              screen: routes.recipientsRoot.name,
              params: { screen: routes.recipientsList.name },
            }}
            mr={8}
          >
            Settings
          </Link>
          <Link to={{ screen: '' }} onPress={logout} mr={8}>
            Log out
          </Link>
        </HStack>
      </HStack>
      <Text variant="display" level="1" mt={12}>
        {getHeaderTitle(route)}
      </Text>
    </VStack>
  );
};
