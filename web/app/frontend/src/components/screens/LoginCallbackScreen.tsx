import React from 'react';
import {Title} from 'react-native-paper';
import {layout} from '../../styles';
import {FlatList, Text, TouchableOpacity, View} from 'react-native';
import useApiClient from "../../utils/clients";
import routes from "../../constants/routes";
import {FormDefinition} from "core-js-api-client/lib/types/types";

const LoginCallbackScreen: React.FC<any> = ({navigation}) => {

    return (
        <View style={layout.body}>
            <Text>LOGIN TEXT</Text>
        </View>
    );
};

export default LoginCallbackScreen;
