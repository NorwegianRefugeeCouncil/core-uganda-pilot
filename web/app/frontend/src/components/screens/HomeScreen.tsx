import React, {useState} from 'react';
import {Subheading, Title} from 'react-native-paper';
import {layout} from '../../styles';
import {Text, View} from 'react-native';
import {DateModal} from "../DateModal";
import {Button} from 'core-design-system'

const HomeScreen = () => {
    const [date, setDate] = useState(new Date(Date.now()))

    return (
        <View style={layout.body}>
            <Title>Home</Title>
            <Subheading>Date picker example</Subheading>
            <DateModal date={date} setDate={setDate}/>
            <Button onPress={() => console.log('integrated design system')}>
                <Text>
                    button
                </Text>
            </Button>
        </View>
    );
};

export default HomeScreen;
