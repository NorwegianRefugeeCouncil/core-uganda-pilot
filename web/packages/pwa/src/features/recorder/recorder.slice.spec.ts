import {selectPostRecords} from "./recorder.slice";

describe("selectPostRecords", () => {
    it('should return the ordered records', function () {

        const got = selectPostRecords({
            recorder: {
                ids: ["record2", "record1"],
                entities: {
                    "record1": {
                        recordId: "record1",
                        formId: "form1",
                        parentRecordId: undefined,
                        values: {
                            "a": "b"
                        }
                    },
                    "record2": {
                        recordId: "record2",
                        formId: "form2",
                        parentRecordId: "record1",
                        values: {
                            "a": "b"
                        }
                    }
                },
            }
        } as any)

        expect(got).toHaveLength(2)
        expect(got[0].id).toEqual("record1")
        expect(got[1].id).toEqual("record1")

    });
})
