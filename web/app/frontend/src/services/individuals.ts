import {Subject} from 'rxjs';
import {Individual} from "core-js-api-client/lib/types/models";
import iamClient from "../utils/clients";

const subject = new Subject();

export function getIndividual(id: string) {
    subject.pipe(iamClient.Parties().Get()).subscribe(console.log);
    return subject.next(id);
}

export function createIndividual(individual: Individual) {
    subject.pipe(iamClient.Individuals().Create()).subscribe(console.log);
    return subject.next(individual);
}
