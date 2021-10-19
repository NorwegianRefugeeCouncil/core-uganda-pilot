import {IAMClient} from "core-js-api-client";
import {Subject} from 'rxjs';
import host from "../constants/host";
import {Individual} from "core-js-api-client/lib/types/models";

const subject = new Subject();

const iamClient = new IAMClient('http', host, {'X-Authenticated-User-Subject': ['test@user.email']});

export function getIndividual(id: string) {
    subject.pipe(iamClient.Parties().Get()).subscribe(console.log);
    return subject.next(id);
}

export function createIndividual(person: Individual) {
    subject.pipe(iamClient.Parties().Create()).subscribe(console.log);
    return subject.next(person);
}
