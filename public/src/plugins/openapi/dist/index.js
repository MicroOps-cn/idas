"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.generateService = exports.getSchema = void 0;
const tslib_1 = require("tslib");
/* eslint-disable global-require */
/* eslint-disable import/no-dynamic-require */
const http_1 = tslib_1.__importDefault(require("http"));
const https_1 = tslib_1.__importDefault(require("https"));
const node_fetch_1 = tslib_1.__importDefault(require("node-fetch"));
const swagger2openapi_1 = tslib_1.__importDefault(require("swagger2openapi"));
const log_1 = tslib_1.__importDefault(require("./log"));
const mockGenerator_1 = require("./mockGenerator");
const serviceGenerator_1 = require("./serviceGenerator");
const getImportStatement = (requestLibPath) => {
    if (requestLibPath && requestLibPath.startsWith('import')) {
        return requestLibPath;
    }
    if (requestLibPath) {
        return `import request from '${requestLibPath}'`;
    }
    return `import { request } from "umi"`;
};
const converterSwaggerToOpenApi = (swagger) => {
    if (!swagger.swagger) {
        return swagger;
    }
    return new Promise((resolve, reject) => {
        swagger2openapi_1.default.convertObj(swagger, {}, (err, options) => {
            log_1.default(['💺 将 Swagger 转化为 openAPI']);
            if (err) {
                reject(err);
                return;
            }
            resolve(options.openapi);
        });
    });
};
const getSchema = (schemaPath) => tslib_1.__awaiter(void 0, void 0, void 0, function* () {
    if (schemaPath.startsWith('http')) {
        const protocol = schemaPath.startsWith('https:') ? https_1.default : http_1.default;
        try {
            const agent = new protocol.Agent({
                rejectUnauthorized: false,
            });
            const json = yield node_fetch_1.default(schemaPath, { agent }).then((rest) => rest.json());
            return json;
        }
        catch (error) {
            // eslint-disable-next-line no-console
            console.log('fetch openapi error:', error);
        }
        return null;
    }
    const schema = require(schemaPath);
    return schema;
});
exports.getSchema = getSchema;
const getOpenAPIConfig = (schemaPath) => tslib_1.__awaiter(void 0, void 0, void 0, function* () {
    const schema = yield exports.getSchema(schemaPath);
    if (!schema) {
        return null;
    }
    const openAPI = yield converterSwaggerToOpenApi(schema);
    return openAPI;
});
// 从 appName 生成 service 数据
const generateService = (_a) => tslib_1.__awaiter(void 0, void 0, void 0, function* () {
    var { requestLibPath, schemaPath, mockFolder, nullable = false } = _a, rest = tslib_1.__rest(_a, ["requestLibPath", "schemaPath", "mockFolder", "nullable"]);
    const openAPI = yield getOpenAPIConfig(schemaPath);
    const requestImportStatement = getImportStatement(requestLibPath);
    const serviceGenerator = new serviceGenerator_1.ServiceGenerator(Object.assign({ namespace: 'API', requestImportStatement, enumStyle: 'string-literal', nullable }, rest), openAPI);
    serviceGenerator.genFile();
    if (mockFolder) {
        yield mockGenerator_1.mockGenerator({
            openAPI,
            mockFolder: mockFolder || './mocks/',
        });
    }
});
exports.generateService = generateService;
