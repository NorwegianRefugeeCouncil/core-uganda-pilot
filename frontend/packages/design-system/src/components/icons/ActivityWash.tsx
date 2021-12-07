import React from "react";
import {Path} from "react-native-svg"
import {IconVariants} from "../../types/icons";

export default (variant: IconVariants) => {
    return (
        <Path
            fill={variant}
            d="M27.1,11.4h1.1c0.4,0,0.7,0.3,0.7,0.7v6.7c0,0.4-0.3,0.7-0.7,0.7h-1.1c-0.4,0-0.7-0.3-0.7-0.7V17H16.2
	c-0.8,0-1.4,0.6-1.4,1.4v0.8h0.4c0.2,0,0.4,0.2,0.4,0.4v1.1c0,0.2-0.2,0.4-0.4,0.4h-3.7c-0.2,0-0.4-0.2-0.4-0.4v-1.1
	c0-0.2,0.2-0.4,0.4-0.4h0.4v-0.8c0-2.5,2-4.4,4.4-4.4H20v-1.9h-1.7c-0.3,0-0.6-0.3-0.6-0.6S18,11,18.3,11h5.2c0.3,0,0.6,0.3,0.6,0.6
	s-0.3,0.6-0.6,0.6h-1.7V14h4.5v-1.9C26.4,11.7,26.7,11.4,27.1,11.4z M11,26.7c0-1.1,1.7-3.4,2.1-4.1c0.1-0.1,0.2-0.1,0.2,0
	c0.5,0.6,2.1,3,2.1,4.1c0,1.2-1,2.3-2.2,2.3S11,28,11,26.7z"
        />
    )
}