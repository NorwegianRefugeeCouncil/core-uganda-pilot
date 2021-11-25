import React from "react";
import {Path} from "react-native-svg"
import {IconVariants} from "../../types/icons";

export default (variant: IconVariants) => {
    return (
        <Path
            fill={variant}
            d="M29.8,11.8H10.2c-0.4,0-0.8,0.3-0.8,0.8v15c0,0.4,0.3,0.8,0.8,0.8h19.5c0.4,0,0.8-0.3,0.8-0.8v-15
	C30.5,12.1,30.2,11.8,29.8,11.8z M28.8,14.3v12.2H11.2V14.3l-0.6-0.5l0.9-1.2l1,0.8h15.1l1-0.8l0.9,1.2L28.8,14.3z M27.5,13.4
	L20,19.3l-7.5-5.9l-1-0.8l-0.9,1.2l0.6,0.5l8,6.2c0.2,0.2,0.5,0.3,0.8,0.3s0.6-0.1,0.8-0.3l8-6.2l0.6-0.5l-0.9-1.2L27.5,13.4z"
        />
    )
}
