import React from 'react';
import {Title} from 'react-native-paper';
import {layout} from '../../styles';
import {Text, View} from 'react-native';
import routes from "../../constants/routes";
import {Button} from "core-design-system";

const DesignSystemDemoScreen = () => {

    return (
        <View style={layout.body}>
            <Title>{routes.designSystem.title}</Title>

            <Button
                onPress={() => console.log('integrated design system')}
            >
                <Text>
                    button
                </Text>
            </Button>
        </View>
    );
};

export default DesignSystemDemoScreen;
