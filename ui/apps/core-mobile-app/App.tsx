import React from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { TopBar } from '@core/core-ui';

export default function App() {
  return (
    <View style={styles.container}>
      <TopBar />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    alignItems: 'center',
    justifyContent: 'center'
  }
});
