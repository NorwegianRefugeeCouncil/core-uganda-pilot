import React from 'react';
import {Title} from 'react-native-paper';
import {layout} from '../../styles';
import {View} from 'react-native';
import routes from "../../constants/routes";
import {Button} from "core-design-system";

const DesignSystemDemoScreen = () => {

    return (
        <View style={layout.body}>
            <Title>{routes.designSystem.title}</Title>

            <Button
                onPress={() => console.log('integrated design system')}
                text={'demo button'}
            />
        </View>
    );
};

export default DesignSystemDemoScreen;
