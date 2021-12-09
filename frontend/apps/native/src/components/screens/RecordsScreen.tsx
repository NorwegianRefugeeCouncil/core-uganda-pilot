import React from 'react';
import { FAB, Title } from 'react-native-paper';
import { FlatList, Text, TouchableOpacity, View } from 'react-native';
import uuidv4 from 'uuid';
import { StackNavigationProp } from '@react-navigation/stack';

import { layout } from '../../styles';
import routes from '../../constants/routes';
import { RecordsStoreProps } from '../../reducers/recordsReducers';
import { RecordsScreenContainerProps, StackParamList } from '../../types/screens';
import testIds from '../../constants/testIds';

interface RecordsScreenProps {
  isLoading: boolean;
  navigation: StackNavigationProp<StackParamList, 'records'>;
  state: RecordsStoreProps;
  route: RecordsScreenContainerProps['route'];
}

const RecordsScreen = ({ isLoading, navigation, route, state }: RecordsScreenProps) => {
  const { formId } = route.params;

  return (
    <View style={[layout.container, layout.body]}>
      <Title>{routes.records.title}</Title>
      {isLoading ? (
        <View>
          <Text>Loading...</Text>
        </View>
      ) : (
        <View>
          <FlatList
            style={{ width: '100%' }}
            data={state.formsById[formId].records}
            keyExtractor={(_, index) => index.toString()}
            renderItem={({ item }) => (
              <TouchableOpacity
                testID={testIds.remoteRecord}
                accessibilityLabel="view record"
                onPress={() =>
                  navigation.navigate(routes.viewRecord.name, {
                    recordId: item.id,
                    formId,
                  })
                }
              >
                <View style={{ flexDirection: 'row', flex: 1 }}>
                  <View style={{ justifyContent: 'center', paddingRight: 12 }}>
                    <Text>{item.id}</Text>
                  </View>
                </View>
              </TouchableOpacity>
            )}
          />
          <FlatList
            style={{ width: '100%' }}
            data={state.formsById[formId].localRecords}
            keyExtractor={(_, index) => index.toString()}
            renderItem={({ item }) => (
              <TouchableOpacity
                testID={testIds.localRecord}
                accessibilityLabel="add record"
                onPress={() =>
                  navigation.navigate(routes.addRecord.name, {
                    recordId: item,
                    formId,
                  })
                }
              >
                <View style={{ flexDirection: 'row', flex: 1 }}>
                  <View style={{ justifyContent: 'center', paddingRight: 12 }}>
                    <Text>{item}</Text>
                  </View>
                </View>
              </TouchableOpacity>
            )}
          />
        </View>
      )}

      <FAB
        style={layout.fab}
        icon="plus"
        color="white"
        onPress={() => navigation.navigate(routes.addRecord.name, { formId, recordId: uuidv4() })}
      />
    </View>
  );
};

export default RecordsScreen;
