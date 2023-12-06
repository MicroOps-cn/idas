"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.normalizeArray = exports.inferSchema = exports.isFunc = exports.objectify = exports.get = exports.isObject = void 0;
function isObject(obj) {
    return !!obj && typeof obj === 'object';
}
exports.isObject = isObject;
function objectify(thing) {
    if (!isObject(thing))
        return {};
    return thing;
}
exports.objectify = objectify;
function get(entity, path) {
    let current = entity;
    for (let i = 0; i < path.length; i += 1) {
        if (current === null || current === undefined) {
            return undefined;
        }
        current = current[path[i]];
    }
    return current;
}
exports.get = get;
function normalizeArray(arr) {
    if (Array.isArray(arr))
        return arr;
    return [arr];
}
exports.normalizeArray = normalizeArray;
function isFunc(thing) {
    return typeof thing === 'function';
}
exports.isFunc = isFunc;
function inferSchema(thing) {
    if (thing.schema) {
        return thing.schema;
    }
    if (thing.properties) {
        return Object.assign(Object.assign({}, thing), { type: 'object' });
    }
    return thing;
}
exports.inferSchema = inferSchema;
