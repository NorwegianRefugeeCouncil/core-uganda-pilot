import { StyleSheet } from 'react-native';

export default StyleSheet.create({
  container: {
    flex: 1,
    flexGrow: 1,
  },
  body: {
    flexGrow: 1,
    alignItems: 'flex-start',
    justifyContent: 'flex-start',
    padding: 18
  },
  fab: {
    position: 'absolute',
    bottom: 28,
    right: 18
  }
});
