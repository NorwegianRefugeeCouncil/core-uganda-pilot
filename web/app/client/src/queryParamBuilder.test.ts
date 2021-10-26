import {expect} from 'chai'
import {objectToQueryString} from './queryParamBuilder'

const testCases = {
    testCase1: {
        obj: {
            a: 1,
            b: 2,
            c: [3, 4, 5, null, 7],
            d: {
                e: 2,
                f: [2, 3, 4]
            },
            g: true
        },
        res: 'a=1&b=2&c=3&c=4&c=5&c=7&d[e]=2&d[f]=2&d[f]=3&d[f]=4&g=true'
    },
    testCase2: {
        obj: {
            a: 1,
            b: 2,
            c: [3, 4, 5, null, 7],
            d: {
                e: 2,
                f: [2, 3, 4]
            }
        },
        res: 'a=1&b=2&c=3&c=4&c=5&c=7&d[e]=2&d[f]=2&d[f]=3&d[f]=4'
    }
}

describe('Unit Tests: Query Param Builder', () => {
    it('Should correctly transform a complex JS object into Core query string format', () => {
        for (const [testCase, testCaseDefinition] of Object.entries(testCases)) {
            expect(objectToQueryString(testCaseDefinition.obj)).to.equal(testCaseDefinition.res)
        }
    })
})
