import { Direction } from '@ui-helpers/types';

export function menuPositionFromDir(direction: Direction) {
  switch (direction) {
    case 'down':
      return {
        top: '100%',
        left: '0',
      };
    case 'up':
      return {
        top: 'auto',
        bottom: '100%',
      };
    case 'right':
    case 'end':
      return {
        top: '0',
        right: 'auto',
        left: '100%',
      };
    case 'left':
    case 'start':
      return {
        top: '0',
        right: '100%',
        left: 'auto',
      };
    default:
      return {};
  }
}
