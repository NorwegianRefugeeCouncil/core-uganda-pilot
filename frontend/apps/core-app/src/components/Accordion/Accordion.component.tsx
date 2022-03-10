import React, { FC, ReactNode } from 'react';
import { Text, Box, Pressable } from 'native-base';

type Props = {
  header: string;
  children: ReactNode;
};

export const AccordionComponent: FC<Props> = ({ header, children }) => {
  const [isExpanded, setIsExpanded] = React.useState(false);

  const handleOnPress = () => setIsExpanded(!isExpanded);

  return (
    <Box p="2">
      <Pressable
        bg="secondary.500"
        p="2"
        onPress={handleOnPress}
        display="flex"
        flexDirection="row"
        justifyContent="space-between"
      >
        <Text color="white" fontSize="18px" lineHeight="21px">
          {header}
        </Text>
        <Text color="white">{isExpanded ? 'open' : 'closed'}</Text>
      </Pressable>
      {isExpanded && (
        <Box bg="secondary.100" p="2">
          {children}
        </Box>
      )}
    </Box>
  );
};
