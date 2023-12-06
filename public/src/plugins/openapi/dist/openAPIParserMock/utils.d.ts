declare function isObject(obj: any): boolean;
declare function objectify(thing: any): any;
declare function get(entity: any, path: (string | number)[]): any;
declare function normalizeArray(arr: any): any[];
declare function isFunc(thing: any): boolean;
declare function inferSchema(thing: any): any;
export { isObject, get, objectify, isFunc, inferSchema, normalizeArray };
