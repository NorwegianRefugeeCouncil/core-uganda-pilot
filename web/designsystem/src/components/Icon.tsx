import React from 'react';
import {IconName, IconProps} from "../types/icons";
import Attachment from '../assets/svg/standard/Attachment'
import Beneficiary from "../assets/svg/standard/Beneficiary";
import {NativeBaseProvider} from "native-base";

const Icon = ({name}: IconProps) => {
    let iconComponent;
    switch (name) {
        case IconName.ATTACHMENT:
            iconComponent = Attachment;
            break;
        case IconName.BENEFICIARY:
            iconComponent = Beneficiary;
            break;
        default:
            iconComponent = Beneficiary;
    }

    console.log('IConComponent', typeof iconComponent, iconComponent)

    return (
        <NativeBaseProvider children={iconComponent()}/>
    );

    // if (name == IconName.ATTACHMENT) {
    //     return Attachment();
    // }
    // return Beneficiary();
}

export default Icon
