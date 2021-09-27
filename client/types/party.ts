/* Do not change, this code is generated from Golang structs */

export class Party {
  id: string;
  partyTypeIds: string[];
  attributes: {[key: string]: string[]};

  constructor(source: any = {}) {
    if ('string' === typeof source) source = JSON.parse(source);
    this.id = source["id"];
    this.partyTypeIds = source["partyTypeIds"];
    this.attributes = source["attributes"];
  }
}