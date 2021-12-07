export {default as Button} from './components/Button';

// this breaks the web/app/frontend.
// it seems like this import is importing a native library that jest cannot run
// TODO: investigate
// export {default as Icon} from './components/Icon';
