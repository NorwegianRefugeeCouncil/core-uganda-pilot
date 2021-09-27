import { of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { request } from 'universal-rxjs-ajax';
import { Party } from './types/party';
import { Attribute } from './types/attributes';
import { Case } from './types/cases';

let test = request(
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
).pipe(map((data) => {

  // const newTypedThing = new Attribute(data.response.items[0]);
  // const newTypedThing = new Party(data.response.items[0]);
  const newTypedThing = new Case(data.response.items[0]);
    console.log(newTypedThing);
  }, catchError(error => {
    console.log('error: ', error);
    return of(error);
  }))
);

test.subscribe(() => {
});

