import { WebBrowserAuthSessionResult, WebBrowserResultType } from '../types/types';
import Browser from '../types/browser';

type Resolve = (value: WebBrowserAuthSessionResult) => void;

export const listener = (event: MessageEvent, browser: Browser, resolve: Resolve) => {
  const { data, isTrusted, origin } = event;

  if (!isTrusted) {
    return;
  }
  if (origin !== window.location.origin) {
    return;
  }

  const handle = window.sessionStorage.getItem(browser.getHandle());

  if (data.sender === handle) {
    browser.dismissPopup();
    resolve({ type: WebBrowserResultType.SUCCESS, url: data.url });
  }
};

export const handler = (resolve: Resolve, browser: Browser) => {
  const localListener = (event: MessageEvent) => listener(event, browser, resolve);

  window.addEventListener('message', localListener, false);
  const interval = setInterval(() => {
    if (browser.getPopup()?.closed) {
      resolve({ type: WebBrowserResultType.DISMISS });
      clearInterval(interval);
      browser.dismissPopup();
    }
  }, 10);
  const popup = browser.getPopup();
  if (popup != null) {
    browser.getListenerMap().set(popup, { listener: localListener, interval });
  }
};
