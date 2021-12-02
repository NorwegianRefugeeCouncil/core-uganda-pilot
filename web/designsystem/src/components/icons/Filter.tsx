import React from "react";
import {Path} from "react-native-svg"
import {IconVariants} from "../../types/icons";

export default (variant: IconVariants) => {
    return (
        <Path
            fill={variant}
            d="M18,26h4v-2h-4V26z M11,14v2h18v-2H11z M14,21h12v-2H14V21z"
        />
    )
}