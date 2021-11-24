import {storiesOf} from '@storybook/react-native';
import React from 'react';
import {View} from 'react-native';
import {Icon} from 'core-design-system'
import {IconName} from "core-design-system/lib/types/icons";

storiesOf('Icon', module)
    // .addDecorator((getStory) =>
    //     <CenterView>{getStory()}</CenterView>
    // )
    .add('Icon', () => {
        console.log('ICON story', typeof Icon, Icon)
        return (
            <View>
                <Icon name={IconName.ATTACHMENT}/>
                <Icon name={IconName.BENEFICIARY}/>
                <Icon name={IconName.CALL}/>
                <Icon name={IconName.CALENDAR}/>
                <Icon name={IconName.CASE}/>
                <Icon name={IconName.CHAT}/>
                <Icon name={IconName.DELETE}/>
                <Icon name={IconName.EDIT}/>
            </View>
        )
    });
