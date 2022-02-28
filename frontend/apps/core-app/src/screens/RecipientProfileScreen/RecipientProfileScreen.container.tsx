import * as React from 'react';
import { RouteProp, useRoute } from '@react-navigation/native';

import { RootParamList } from '../../navigation/types';

import { RecipientProfileScreenComponent } from './RecipientProfileScreen.component';

export const RecipientProfileScreenContainer: React.FC = () => {
  const route = useRoute<RouteProp<RootParamList, 'RecipientProfile'>>();
  return <RecipientProfileScreenComponent route={route} />;
};
