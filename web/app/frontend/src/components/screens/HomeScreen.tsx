import React, {useState} from 'react';
import {Subheading, Text, Title} from 'react-native-paper';
import {layout} from '../../styles';
import {View} from 'react-native';
import {DateModal} from "../DateModal";

const HomeScreen = () => {
    const [date, setDate] = useState(new Date(Date.now()))

    return (
        <View style={layout.body}>
            <Title>Home</Title>
            <Subheading>Date picker example</Subheading>
            <DateModal date={date} setDate={setDate}/>
        </View>
    );
};

export default HomeScreen;
