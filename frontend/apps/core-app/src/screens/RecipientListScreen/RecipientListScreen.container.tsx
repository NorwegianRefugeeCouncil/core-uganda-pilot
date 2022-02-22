import * as React from 'react';
import {
  NavigationProp,
  RouteProp,
  useNavigation,
  useRoute,
} from '@react-navigation/native';

import { RootParamList } from '../../navigators/types';

import { RecipientListScreenComponent } from './RecipientListScreen.component';

export const RecipientListScreenContainer: React.FC = () => {
  const route = useRoute<RouteProp<RootParamList, 'RecipientList'>>();
  const navigation = useNavigation<NavigationProp<RootParamList>>();
  return <RecipientListScreenComponent route={route} navigation={navigation} />;
};
