import * as React from 'react';
import { Link } from 'react-router-dom';

export const TopNav: React.FC = props => {
  return (
    <nav className='navbar navbar-expand-lg navbar-light bg-light shadow-sm'>
      <div className='container-fluid'>
        <a className='navbar-brand' href='#'>Navbar</a>
        <button className='navbar-toggler' type='button' data-bs-toggle='collapse'
                data-bs-target='#navbarSupportedContent' aria-controls='navbarSupportedContent' aria-expanded='false'
                aria-label='Toggle navigation'>
          <span className='navbar-toggler-icon' />
        </button>
        <div className='collapse navbar-collapse' id='navbarSupportedContent'>
          <ul className='navbar-nav me-auto mb-2 mb-lg-0'>
            <li className='nav-item'>
              <Link to='/' className='nav-link active'>
                Home
              </Link>
            </li>
            <li className='nav-item'>
              <Link to='/formdefinitions' className='nav-link'>
                Form Definitions
              </Link>
            </li>
            <li className='nav-item'>
              <Link to='/store' className='nav-link'>
                Store
              </Link>
            </li>
          </ul>
        </div>
      </div>
    </nav>
  );
};
