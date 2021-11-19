"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.resolveDiscoveryAsync = void 0;
const axios_1 = __importDefault(require("axios"));
function resolveDiscoveryAsync(issuerOrDiscovery) {
    return __awaiter(this, void 0, void 0, function* () {
        let issuer;
        if (typeof issuerOrDiscovery === "string") {
            issuer = issuerOrDiscovery;
        }
        else {
            issuer = issuerOrDiscovery.issuer;
        }
        const metadataEndpoint = `${issuer}/.well-known/openid-configuration`;
        return axios_1.default.get(metadataEndpoint)
            .then(value => value.data)
            .catch(err => {
            console.log("failed to get discovery document", err);
            return null;
        });
    });
}
exports.resolveDiscoveryAsync = resolveDiscoveryAsync;
