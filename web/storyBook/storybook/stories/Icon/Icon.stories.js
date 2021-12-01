import {storiesOf} from '@storybook/react-native';
import React from 'react';
import {Icon} from 'core-design-system'
import {IconName, IconVariants} from "core-design-system/lib/types/icons";
import CenterView from "../CenterView";
import {select} from "@storybook/addon-knobs";

storiesOf('Icon', module)
    .addDecorator((getStory) =>
        <CenterView>
            {getStory()}
        </CenterView>
    )
    .add('Icon', () => {
        const IconNameList = Object.entries(IconName)
            .map(([_, value]) => (value));
        const IconVariantList = Object.entries(IconVariants)
            .map(([_, value]) => (value));

        return (
            <Icon
                name={select('icon name', IconNameList, IconName.ATTACHMENT)}
                variant={select('variant', IconVariantList, IconVariants.DARK)}
            />
        )
    });
