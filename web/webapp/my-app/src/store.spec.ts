import {fetchDatabases, newStore, reducer, setupFetchDBEffects} from "./store";
import {DatabaseList} from "./client";
import {of, OperatorFunction, throwError} from "rxjs";


describe("store", () => {

    it('should create the store', function () {
        const store = newStore(reducer)
        expect(store).not.toBeFalsy()
    });

    it('should fetch the databases', function (done) {
        const store = newStore(reducer)
        const client = {
            getDatabases(): OperatorFunction<null, DatabaseList> {
                return () => of({items: [{name: "database"}]})
            }
        }
        setupFetchDBEffects(store, client)
        store.state$.subscribe(s => {
            console.log(s)
            if (s.databases.length > 0 && !s.pending && !s.error) {
                done()
            }
        })
        store.dispatch(fetchDatabases())
    });

    it('should apply error from fetching the databases', function (done) {
        const store = newStore(reducer)
        const client = {
            getDatabases(): OperatorFunction<null, DatabaseList> {
                return () => throwError(() => "error")
            }
        }
        setupFetchDBEffects(store, client)
        store.state$.subscribe(s => {
            if (s.error && !s.pending) {
                done()
            }
        })
        store.dispatch(fetchDatabases())
    });

})
