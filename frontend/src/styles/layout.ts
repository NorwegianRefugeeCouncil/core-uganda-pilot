import { StyleSheet } from 'react-native';
import theme from '../constants/theme';

export default StyleSheet.create({
  container: {
    flex: 1,
    flexGrow: 1,
  },
  navigation: {
    paddingHorizontal: 16,
    paddingVertical: 10,
    justifyContent: 'space-between',
    alignContent: 'center',
    alignItems: 'center'
  },
  body: {
    flexGrow: 1,
    alignItems: 'flex-start',
    justifyContent: 'flex-start',
    padding: 18,
  }
});
