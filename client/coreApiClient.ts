import { of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { Party } from './types/party';
import { Attribute } from './types/attributes';
import { Case } from './types/cases';
import { ajax } from 'rxjs/ajax';
import { XMLHttpRequest } from 'xhr2';

global.XMLHttpRequest = global.XMLHttpRequest ? global.XMLHttpRequest : XMLHttpRequest;

let test = ajax(
  {
    url: 'http://localhost:9000/apis/cms/v1/cases',
    // url: 'http://localhost:9000/apis/iam/v1/attributes',
    // url: 'http://localhost:9000/apis/iam/v1/parties',
    headers: {
      'X-Authenticated-User-Subject': 'stephen.kabagambe@email.com'
    },
    method: 'GET',
    async: true,
    timeout: 0,
    crossDomain: true,
    withCredentials: false,
    responseType: 'json'
  }
).pipe(map((response) => {

  // const newTypedThing = new Attribute(data.response.items[0]);
  // const newTypedThing = new Party(data.response.items[0]);
  // const newTypedThing = new Case(data.response.items[0]);
    console.log(response);
  }, catchError(error => {
    console.log('error: ', error);
    return of(error);
  }))
);

test.subscribe(() => {
});

