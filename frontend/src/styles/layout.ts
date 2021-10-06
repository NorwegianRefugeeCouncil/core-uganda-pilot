import { StyleSheet } from 'react-native';
import theme from '../constants/theme';

export default StyleSheet.create({
  container: {
    flex: 1,
    // alignItems: 'center',
    // justifyContent: 'center',
    backgroundColor: theme.colors.green
  },
  navigation: {
    // flex: 1,
    display: 'flex',
    width: 230,
    alignItems: 'flex-start',
    justifyContent: 'flex-start',
    backgroundColor: theme.colors.red
  },
  body: {
    flex: 1,
    alignItems: 'flex-start',
    justifyContent: 'flex-start',
    backgroundColor: theme.colors.yellow
  }
});
