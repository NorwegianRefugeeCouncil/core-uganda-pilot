import { Button } from 'core-design-system';
import React from 'react';
import { Text, View } from 'react-native';
import { Title } from 'react-native-paper';

import { layout } from '../../styles';
import { DesignSystemScreenProps } from '../../types/screens';

const DesignSystemDemoScreen = ({}: DesignSystemScreenProps) => {
    return (
        <View style={layout.body}>
            <Title>Design System Demo</Title>

            <Button onPress={() => console.log('integrated design system')}>
                <Text>button</Text>
            </Button>
        </View>
    );
};

export default DesignSystemDemoScreen;
