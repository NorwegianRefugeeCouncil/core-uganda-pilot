import React, { FC, ReactNode } from 'react';
import { Text, Box, Pressable } from 'native-base';

type Props = {
  header: string;
  content: ReactNode;
};

export const AccordionComponent: FC<Props> = ({ header, content }) => {
  const [isActive, setIsActive] = React.useState(false);

  return (
    <Box>
      <Pressable bg="neutral.200" p="2" onPress={() => setIsActive(!isActive)}>
        <Text variant="heading">{header}</Text>
      </Pressable>
      {isActive && (
        <Box bg="neutral.100" p="2">
          {content}
        </Box>
      )}
    </Box>
  );
};
