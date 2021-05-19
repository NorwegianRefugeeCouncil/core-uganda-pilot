import * as React from 'react';

export function useAccordion({ id, onPointerDown = null }) {
  const [activeKey, setActiveKey] = React.useState(id);

  const handlePointerDown = (id: string) => {
    if (onPointerDown != null) onPointerDown(id);
    console.log(id);
    setActiveKey(id);
  };

  return { activeKey, handlePointerDown };
}
