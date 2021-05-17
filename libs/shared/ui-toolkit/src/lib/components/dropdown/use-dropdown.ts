import * as React from 'react';

function useDropdown(onChange: (value: any) => void) {
  const toggleBtnRef = React.useRef(null);
  const menuRef = React.useRef(null);

  const [menuIsOpen, setMenuIsOpen] = React.useState(false);

  const toggleMenu = (event) => setMenuIsOpen(!menuIsOpen);

  // Attach listener to document to close menu if click outside
  const hideMenu = (event) => {
    const menu = menuRef.current;
    if (menu === undefined || menu.contains(event.target)) {
      return;
    }
    setMenuIsOpen(false);
  };

  React.useEffect(() => {
    if (menuIsOpen) {
      document.addEventListener('pointerdown', hideMenu);
    } else {
      document.removeEventListener('pointerdown', hideMenu);
    }
    return () => {
      document.removeEventListener('pointerdown', hideMenu);
    };
  }, [menuIsOpen]);

  const handleChange = (selectedValue: any) => {
    if (onChange != null) onChange(selectedValue);
    setMenuIsOpen(false);
  };

  return {
    menuRef,
    toggleBtnRef,
    menuIsOpen,
    toggleMenu,
    handleChange,
  };
}

export default useDropdown;
