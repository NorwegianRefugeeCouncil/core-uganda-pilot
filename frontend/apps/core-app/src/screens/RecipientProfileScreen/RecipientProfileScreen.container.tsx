import * as React from 'react';
import { RouteProp, useRoute } from '@react-navigation/native';

import { RootParamList } from '../../navigators/types';

import { RecipientProfileScreenComponent } from './RecipientProfileScreen.component';

export const RecipientProfileScreenContainer: React.FC = () => {
  const route = useRoute<RouteProp<RootParamList, 'RecipientProfile'>>();
  return <RecipientProfileScreenComponent route={route} />;
};
