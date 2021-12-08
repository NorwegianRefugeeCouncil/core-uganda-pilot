import {findSubForm, hasSubFormWithId, selectFormOrSubFormById, selectRootForm} from "../form";
import {RootState} from "../../app/store";

describe("form reducer", () => {
    describe("findSubForm", () => {
        it('should find the first level sub form', function () {
            const found = findSubForm("subformId", [{
                fieldType: {subForm: {fields: []}},
                id: 'subformId',
                code: 'code',
                name: 'name',
                description: 'description',
                required: false,
                key: false,
                options: []
            }])
            expect(found?.id).toEqual("subformId")
        });
        it('should find a nested sub form', function () {
            const found = findSubForm("nested", [{
                fieldType: {
                    subForm: {
                        fields: [{
                            fieldType: {
                                subForm: {
                                    fields: []
                                }
                            },
                            id: 'nested',
                            code: 'code',
                            name: 'name',
                            description: 'description',
                            required: false,
                            key: false,
                            options: []
                        }]
                    }
                },
                id: 'mainSubform',
                code: 'code',
                name: 'name',
                description: 'description',
                required: false,
                key: false,
                options: []
            }])
            expect(found?.id).toEqual("nested")
        });
        it('should return undefined if not found', function () {
            const found = findSubForm("bla", [{
                fieldType: {subForm: {fields: []}},
                id: 'foo',
                code: 'code',
                name: 'name',
                description: 'description',
                required: false,
                key: false,
                options: []
            }])
            expect(found?.id).toBeUndefined()
        });
    })
    describe("hasSubFormWithId", () => {
        it('should return true if has subform with id', function () {
            const found = hasSubFormWithId("subformId", [{
                fieldType: {subForm: {fields: []}},
                id: 'subformId',
                code: 'code',
                name: 'name',
                description: 'description',
                required: false,
                key: false,
                options: []
            }])
            expect(found).toBeTruthy()
        });
        it('should return false if has subform with no such id', function () {
            const found = hasSubFormWithId("bla", [{
                fieldType: {subForm: {fields: []}},
                id: 'subformId',
                code: 'code',
                name: 'name',
                description: 'description',
                required: false,
                key: false,
                options: []
            }])
            expect(found).toBeFalsy()
        });
        it('should return true with a nested subform', function () {
            const found = hasSubFormWithId("nested", [{
                fieldType: {
                    subForm: {
                        fields: [{
                            fieldType: {
                                subForm: {
                                    fields: []
                                }
                            },
                            id: 'nested',
                            code: 'code',
                            name: 'name',
                            description: 'description',
                            required: false,
                            key: false,
                            options: []
                        }]
                    }
                },
                id: 'mainSubform',
                code: 'code',
                name: 'name',
                description: 'description',
                required: false,
                key: false,
                options: []
            }])
            expect(found).toBeTruthy()
        });
    })
    describe("selectRootForm", () => {
        it('should return the root form id if given the root form id', function () {
            let state = {
                forms: {
                    ids: ["form"],
                    entities: {
                        "form": {id: "form"}
                    }
                }
            } as unknown;
            const found = selectRootForm(state as RootState, "form")
            expect(found?.id).toEqual("form")
        });
        it('should return the root form id if given the child form id', function () {
            let state = {
                forms: {
                    ids: ["form"],
                    entities: {
                        "form": {
                            id: "form", fields: [{
                                fieldType: {subForm: {}},
                                id: "subform"
                            }]
                        }
                    }
                }
            } as unknown;
            const found = selectRootForm(state as RootState, "subform")
            expect(found?.id).toEqual("form")
        });
    })
    describe("selectFormOrSubFormById", () => {
        it('should return the root form if given a root form id', function () {
            let state = {
                forms: {
                    ids: ["form"],
                    entities: {
                        "form": {id: "form", fields: [{fieldType: {subForm: {id: "subform"}}}]}
                    }
                }
            } as unknown;
            const found = selectFormOrSubFormById(state as RootState, "form")
            expect(found?.id).toEqual("form")
        });
        it('should return the root form if given a child form id', function () {
            let state = {
                forms: {
                    ids: ["form"],
                    entities: {
                        "form": {
                            id: "form", fields: [{
                                fieldType: {subForm: {}},
                                id: "subform"
                            }]
                        }
                    }
                }
            } as unknown;
            const found = selectFormOrSubFormById(state as RootState, "subform")
            expect(found?.id).toEqual("subform")
        });
    })
})
