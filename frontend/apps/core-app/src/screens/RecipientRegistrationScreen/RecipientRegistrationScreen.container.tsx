import * as React from 'react';
import { RouteProp, useRoute } from '@react-navigation/native';

import { RootParamList } from '../../navigators/types';

import { RecipientRegistrationScreenComponent } from './RecipientRegistrationScreen.component';

export const RecipientRegistrationScreenContainer: React.FC = () => {
  const route = useRoute<RouteProp<RootParamList, 'RecipientRegistration'>>();
  return <RecipientRegistrationScreenComponent route={route} />;
};
