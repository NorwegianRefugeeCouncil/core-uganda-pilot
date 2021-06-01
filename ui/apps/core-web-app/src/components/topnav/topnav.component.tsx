import * as React from 'react';
import { Link } from 'react-router-dom';
import { SVGProps } from 'react';


type IconType = 'home' | 'envelope' | 'person' | 'gear'


const renderIcon = (iconType: IconType) => {

  let path: SVGProps<SVGPathElement>;
  let className = 'text-white flex-shrink-1 bi';

  if (iconType === 'home') {
    path = <path
      d='M6.5 14.5v-3.505c0-.245.25-.495.5-.495h2c.25 0 .5.25.5.5v3.5a.5.5 0 0 0 .5.5h4a.5.5 0 0 0 .5-.5v-7a.5.5 0 0 0-.146-.354L13 5.793V2.5a.5.5 0 0 0-.5-.5h-1a.5.5 0 0 0-.5.5v1.293L8.354 1.146a.5.5 0 0 0-.708 0l-6 6A.5.5 0 0 0 1.5 7.5v7a.5.5 0 0 0 .5.5h4a.5.5 0 0 0 .5-.5z' />;
    className += ' bi-house-door-fill';
  }
  if (iconType === 'envelope') {
    path = <path
      d='M0 4a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2V4zm2-1a1 1 0 0 0-1 1v.217l7 4.2 7-4.2V4a1 1 0 0 0-1-1H2zm13 2.383-4.758 2.855L15 11.114v-5.73zm-.034 6.878L9.271 8.82 8 9.583 6.728 8.82l-5.694 3.44A1 1 0 0 0 2 13h12a1 1 0 0 0 .966-.739zM1 11.114l4.758-2.876L1 5.383v5.73z' />;
    className += ' bi-bi-envelope';
  }
  if (iconType === 'person') {
    path = <path
      d='M8 8a3 3 0 1 0 0-6 3 3 0 0 0 0 6zm2-3a2 2 0 1 1-4 0 2 2 0 0 1 4 0zm4 8c0 1-1 1-1 1H3s-1 0-1-1 1-4 6-4 6 3 6 4zm-1-.004c-.001-.246-.154-.986-.832-1.664C11.516 10.68 10.289 10 8 10c-2.29 0-3.516.68-4.168 1.332-.678.678-.83 1.418-.832 1.664h10z' />;
    className += ' bi-person';
  }
  if (iconType === 'gear') {
    path = <path
      d='M8.932.727c-.243-.97-1.62-.97-1.864 0l-.071.286a.96.96 0 0 1-1.622.434l-.205-.211c-.695-.719-1.888-.03-1.613.931l.08.284a.96.96 0 0 1-1.186 1.187l-.284-.081c-.96-.275-1.65.918-.931 1.613l.211.205a.96.96 0 0 1-.434 1.622l-.286.071c-.97.243-.97 1.62 0 1.864l.286.071a.96.96 0 0 1 .434 1.622l-.211.205c-.719.695-.03 1.888.931 1.613l.284-.08a.96.96 0 0 1 1.187 1.187l-.081.283c-.275.96.918 1.65 1.613.931l.205-.211a.96.96 0 0 1 1.622.434l.071.286c.243.97 1.62.97 1.864 0l.071-.286a.96.96 0 0 1 1.622-.434l.205.211c.695.719 1.888.03 1.613-.931l-.08-.284a.96.96 0 0 1 1.187-1.187l.283.081c.96.275 1.65-.918.931-1.613l-.211-.205a.96.96 0 0 1 .434-1.622l.286-.071c.97-.243.97-1.62 0-1.864l-.286-.071a.96.96 0 0 1-.434-1.622l.211-.205c.719-.695.03-1.888-.931-1.613l-.284.08a.96.96 0 0 1-1.187-1.186l.081-.284c.275-.96-.918-1.65-1.613-.931l-.205.211a.96.96 0 0 1-1.622-.434L8.932.727zM8 12.997a4.998 4.998 0 1 1 0-9.995 4.998 4.998 0 0 1 0 9.996z' />;
    className += ' bi-gear-wide';
  }

  return <svg xmlns='http://www.w3.org/2000/svg'
              fill='currentColor'
              className={className}
              viewBox='0 0 16 16'
              style={{
                height: '1.5rem'
              }}
  >
    {path}
  </svg>;
};

type NavItemProps = {
  iconType: IconType
  label: string,
  width: string
  to: string
}


export const NavItem: React.FC<NavItemProps> = (props) => {
  const { iconType, label, width, to } = props;
  return <li className={'nav-item flex-shrink-1'} style={{
    width: width,
    maxWidth: width
  }}>
    <Link to={to} className={'d-flex flex-column align-items-center px-2'}>
      {renderIcon(iconType)}
      <span className={'fw-bold pb-1 small text-white'}>{label}</span>
    </Link>
  </li>;
};

export const TopNav: React.FC = props => {

  return (
    <ul
      style={{
        boxShadow: '0 -0.1rem 0.5rem rgba(0, 0, 0, 0.1)',
        width: '100vw',
        maxWidth: '100vw',
        height: '3.5rem'
      }}
      className={'nav nav-pills fixed-bottom bg-primary text-white justify-content-evenly flex-nowrap px-2 pt-2'}
    >
      <NavItem iconType='home' label='Home' width={'3rem'} to='/' />
      <NavItem iconType='person' label='Beneficiaries' width={'6rem'} to='/beneficiaries' />
      <NavItem iconType='envelope' label='Messages' width={'4.5rem'} to='/messages' />
      <NavItem iconType='home' label='Home' width={'3rem'} to='/' />
      <NavItem iconType='gear' label='Settings' width={'4rem'} to='/settings' />
    </ul>
  );
};
