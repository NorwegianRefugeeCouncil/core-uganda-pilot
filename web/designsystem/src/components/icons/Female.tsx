import React from "react";
import {Path} from "react-native-svg"
import {IconVariants} from "../../types/icons";

export default (variant: IconVariants) => {
    return (
        <Path
            fill={variant}
            d="M24.7,20.9c1.3-1.3,2-2.9,2-4.7c0-1.8-0.7-3.5-2-4.7c-1.3-1.3-2.9-1.9-4.7-1.9s-3.5,0.7-4.7,1.9c-1.3,1.3-2,2.9-2,4.7
	c0,1.5,0.5,2.9,1.4,4.1c0.2,0.2,0.4,0.4,0.6,0.6c0.2,0.2,0.4,0.4,0.6,0.6c0.9,0.7,2,1.2,3.2,1.3v2.5h-2.7c-0.1,0-0.2,0.1-0.2,0.2
	v1.4c0,0.1,0.1,0.2,0.2,0.2h2.7v3.3c0,0.1,0.1,0.2,0.2,0.2h1.4c0.1,0,0.2-0.1,0.2-0.2V27h2.7c0.1,0,0.2-0.1,0.2-0.2v-1.4
	c0-0.1-0.1-0.2-0.2-0.2h-2.7v-2.5C22.3,22.6,23.7,21.9,24.7,20.9z M20,21c-1.3,0-2.5-0.5-3.4-1.4c-0.9-0.9-1.4-2.1-1.4-3.4
	c0-1.3,0.5-2.5,1.4-3.4c0.9-0.9,2.1-1.4,3.4-1.4s2.5,0.5,3.4,1.4c0.9,0.9,1.4,2.1,1.4,3.4c0,1.3-0.5,2.5-1.4,3.4
	C22.5,20.5,21.3,21,20,21z"
        />
    )
}
