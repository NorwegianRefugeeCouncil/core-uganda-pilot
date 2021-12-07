import React from "react";
import {Path} from "react-native-svg"
import {IconVariants} from "../../types/icons";

export default (variant: IconVariants) => {
    return (
        <Path
            fill={variant}
            d="M19.2,20.8v8.5h1.5v-8.5h8.5v-1.5h-8.5v-8.5h-1.5v8.5h-8.5v1.5H19.2z"
        />
    )
}
