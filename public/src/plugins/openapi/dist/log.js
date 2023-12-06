"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const tslib_1 = require("tslib");
const chalk_1 = tslib_1.__importDefault(require("chalk"));
// eslint-disable-next-line no-console
const Log = (...rest) => console.log(`${chalk_1.default.blue('[openAPI]')}: ${rest.join('\n')}`);
exports.default = Log;
