"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.mockGenerator = void 0;
const tslib_1 = require("tslib");
const mockjs_1 = tslib_1.__importDefault(require("mockjs"));
const fs_1 = tslib_1.__importDefault(require("fs"));
const util_1 = require("./util");
const path_1 = require("path");
const index_1 = tslib_1.__importDefault(require("./openAPIParserMock/index"));
const log_1 = tslib_1.__importDefault(require("./log"));
const tiny_pinyin_1 = tslib_1.__importDefault(require("tiny-pinyin"));
mockjs_1.default.Random.extend({
    country() {
        const data = [
            '阿根廷',
            '澳大利亚',
            '巴西',
            '加拿大',
            '中国',
            '法国',
            '德国',
            '印度',
            '印度尼西亚',
            '意大利',
            '日本',
            '韩国',
            '墨西哥',
            '俄罗斯',
            '沙特阿拉伯',
            '南非',
            '土耳其',
            '英国',
            '美国',
        ];
        const id = (Math.random() * data.length).toFixed();
        return data[id];
    },
    phone() {
        const phonepreFix = ['111', '112', '114']; // 自己写前缀哈
        return this.pick(phonepreFix) + mockjs_1.default.mock(/\d{8}/); // Number()
    },
    status() {
        const status = ['success', 'error', 'default', 'processing', 'warning'];
        return status[(Math.random() * 4).toFixed(0)];
    },
    authority() {
        const status = ['admin', 'user', 'guest'];
        return status[(Math.random() * status.length).toFixed(0)];
    },
    avatar() {
        const avatar = [
            'https://gw.alipayobjects.com/zos/rmsportal/KDpgvguMpGfqaHPjicRK.svg',
            'https://gw.alipayobjects.com/zos/rmsportal/udxAbMEhpwthVVcjLXik.png',
            'https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png',
            'https://gw.alipayobjects.com/zos/rmsportal/ThXAXghbEsBCCSDihZxY.png',
            'https://gw.alipayobjects.com/zos/rmsportal/OKJXDXrmkNshAMvwtvhu.png',
            'https://avatars0.githubusercontent.com/u/507615?s=40&v=4',
            'https://avatars1.githubusercontent.com/u/8186664?s=40&v=4',
        ];
        const id = (Math.random() * avatar.length).toFixed();
        return avatar[id];
    },
    group() {
        const data = ['体验技术部', '创新科技组', '前端 6 组', '区块链平台部', '服务技术部'];
        const id = (Math.random() * data.length).toFixed();
        return data[id];
    },
    label() {
        const label = [
            '很有想法的',
            '小清新',
            '傻白甜',
            '阳光少年',
            '大咖',
            '健身达人',
            '程序员',
            '算法工程师',
            '川妹子',
            '名望程序员',
            '大长腿',
            '海纳百川',
            '专注设计',
            '爱好广泛',
            'IT 互联网',
        ];
        const id = (Math.random() * label.length).toFixed();
        return label[id];
    },
    href() {
        const href = [
            'https://preview.pro.ant.design/dashboard/analysis',
            'https://ant.design',
            'https://procomponents.ant.design/',
            'https://umijs.org/',
            'https://github.com/umijs/dumi',
        ];
        const id = (Math.random() * href.length).toFixed();
        return href[id];
    },
});
const genMockData = (example) => {
    if (!example) {
        return {};
    }
    if (typeof example === 'string') {
        return mockjs_1.default.mock(example);
    }
    if (Array.isArray(example)) {
        return mockjs_1.default.mock(example);
    }
    return Object.keys(example)
        .map((name) => {
        return {
            [name]: mockjs_1.default.mock(example[name]),
        };
    })
        .reduce((pre, next) => {
        return Object.assign(Object.assign({}, pre), next);
    }, {});
};
const genByTemp = ({ method, path, parameters, status, data, }) => {
    if (!['get', 'put', 'post', 'delete', 'patch'].includes(method.toLocaleLowerCase())) {
        return '';
    }
    let securityPath = path;
    parameters === null || parameters === void 0 ? void 0 : parameters.forEach(item => {
        if (item.in === "path") {
            securityPath = securityPath.replace(`{${item.name}}`, `:${item.name}`);
        }
    });
    return `'${method.toUpperCase()} ${securityPath}': (req: Request, res: Response) => {
    res.status(${status}).send(${data});
  }`;
};
const genMockFiles = (mockFunction) => {
    return util_1.prettierFile(` 
// @ts-ignore
import { Request, Response } from 'express';

export default {
${mockFunction.join('\n,')}
    }`)[0];
};
const mockGenerator = ({ openAPI, mockFolder }) => tslib_1.__awaiter(void 0, void 0, void 0, function* () {
    const openAPParse = new index_1.default(openAPI);
    const docs = openAPParse.parser();
    const pathList = Object.keys(docs.paths);
    const { paths } = docs;
    const mockActionsObj = {};
    pathList.forEach((path) => {
        const pathConfig = paths[path];
        Object.keys(pathConfig).forEach((method) => {
            var _a, _b, _c, _d;
            const methodConfig = pathConfig[method];
            if (methodConfig) {
                let conte = (_b = (methodConfig.operationId || ((_a = methodConfig === null || methodConfig === void 0 ? void 0 : methodConfig.tags) === null || _a === void 0 ? void 0 : _a.join('/')) ||
                    path.replace('/', '').split('/')[1])) === null || _b === void 0 ? void 0 : _b.replace(/[^\w^\s^\u4e00-\u9fa5]/gi, '');
                if (/[\u3220-\uFA29]/.test(conte)) {
                    conte = tiny_pinyin_1.default.convertToPinyin(conte, '', true);
                }
                if (!conte) {
                    return;
                }
                const data = genMockData((_d = (_c = methodConfig.responses) === null || _c === void 0 ? void 0 : _c['200']) === null || _d === void 0 ? void 0 : _d.example);
                if (!mockActionsObj[conte]) {
                    mockActionsObj[conte] = [];
                }
                const tempFile = genByTemp({
                    method,
                    path,
                    parameters: methodConfig.parameters,
                    status: '200',
                    data: JSON.stringify(data),
                });
                if (tempFile) {
                    mockActionsObj[conte].push(tempFile);
                }
            }
        });
    });
    Object.keys(mockActionsObj).forEach((file) => {
        if (!file || file === 'undefined') {
            return;
        }
        if (file.includes('/')) {
            const dirName = path_1.dirname(path_1.join(mockFolder, `${file}.mock.ts`));
            if (!fs_1.default.existsSync(dirName)) {
                fs_1.default.mkdirSync(dirName);
            }
        }
        util_1.writeFile(mockFolder, `${file}.mock.ts`, genMockFiles(mockActionsObj[file]));
    });
    log_1.default('✅ 生成 mock 文件成功');
});
exports.mockGenerator = mockGenerator;
