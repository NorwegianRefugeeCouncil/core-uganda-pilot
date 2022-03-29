import React, { FC, ReactNode } from 'react';
import { Text, Box, Pressable } from 'native-base';

type Props = {
  header: string;
  children: ReactNode;
  defaultOpen?: boolean;
};

export const Accordion: FC<Props> = ({
  header,
  children,
  defaultOpen = false,
}) => {
  const [isExpanded, setIsExpanded] = React.useState(defaultOpen);

  const handleOnPress = () => setIsExpanded(!isExpanded);

  return (
    <Box>
      <Pressable bg="secondary.500" p="2" onPress={handleOnPress}>
        <Text color="white" fontSize="18px" lineHeight="21px">
          {header}
        </Text>
      </Pressable>
      {isExpanded && (
        <Box bg="secondary.100" p="2">
          {children}
        </Box>
      )}
    </Box>
  );
};
