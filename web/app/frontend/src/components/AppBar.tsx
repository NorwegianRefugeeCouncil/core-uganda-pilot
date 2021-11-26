import React from 'react';
import { Box, HStack, Icon, IconButton, StatusBar, Text } from 'native-base';
import { AntDesign, MaterialIcons } from '@expo/vector-icons';
import {StackHeaderProps, StackNavigationOptions, StackNavigationProp} from "@react-navigation/stack";
import {ParamListBase} from "@react-navigation/native";

export type AppBarProps = {
    children?: React.ReactNode
    navigation: StackNavigationProp<ParamListBase, string>
    options: StackNavigationOptions
    back?: StackHeaderProps['back']
}

export const AppBar = ({ navigation, back, children }: AppBarProps) => {
    return (
        <>
            <StatusBar backgroundColor="#3700B3" barStyle="light-content" />

            <Box safeAreaTop backgroundColor="#6200ee" />

            <HStack
                bg="#6200ee"
                px="1"
                py="3"
                justifyContent="space-between"
                alignItems="center"
            >
                <HStack space="4" alignItems="center">
                    {back && (
                        <IconButton
                            _icon={{ as: AntDesign, name: 'back' }}
                            size="sm"
                            color="white"
                            onPress={navigation.goBack}
                            accessibilityLabel={'back'}
                        />
                    )}
                </HStack>
                <HStack space="2">
                    {children}
                </HStack>
            </HStack>
        </>
    );
};
