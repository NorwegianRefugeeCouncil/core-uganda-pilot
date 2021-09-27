"use strict";
exports.__esModule = true;
var ajax_1 = require("rxjs/ajax");
var rxjs_1 = require("rxjs");
var operators_1 = require("rxjs/operators");
// let testSubject = new Subject();
var test = ajax_1.ajax({
    url: 'http://localhost:9000/apis/iam/v1/parties',
    headers: {
        'X-Authenticated-User-Subject': 'stephen.kabagambe@email.com'
    },
    method: 'GET'
}).pipe(operators_1.map(function (data) {
    console.log(data);
}, operators_1.catchError(function (error) {
    console.log('error: ', error);
    return rxjs_1.of(error);
})));
test.subscribe(function (data) {
    console.log(data);
});
//
// const myObserver = {
//   next: (x: number) => console.log('Observer got a next value: ' + x),
//   error: (err: Error) => console.error('Observer got an error: ' + err),
// };
//
// // Execute with the observer object
// test.subscribe(myObserver);
// test.subscribe((data)=> {
// })
