import type { AxiosResponse } from 'axios';
import axios from 'axios';
import fs from 'fs';
import path from 'path';
import type { IApi } from 'umi';

import en_US from '../../src/locales/en-US';
import ja_JP from '../../src/locales/ja-JP';
import pt_BR from '../../src/locales/pt-BR';
import zh_CN from '../../src/locales/zh-CN';
import zh_TW from '../../src/locales/zh-TW';

// import bn_BD from '../../src/locales/bn-BD';
// import fa_IR from '../../src/locales/fa-IR';
// import id_ID from '../../src/locales/id-ID';

const allLocales = new Map<string, { locale: any; code: string; fileLocales: Record<string, any> }>(
  [
    // ["bn-BD", { locale: bn_BD, fileLocales: {}, code: 'ben' }],
    // ["id-ID", { locale: id_ID, fileLocales: {}, code: 'id' }],
    // ["fa-IR", { locale: fa_IR, fileLocales: {}, code: 'per' }], //波斯语
    ['en-US', { locale: en_US, fileLocales: {}, code: 'en' }],
    ['ja-JP', { locale: ja_JP, fileLocales: {}, code: 'jp' }],
    ['pt-BR', { locale: pt_BR, fileLocales: {}, code: 'pt' }], //葡萄牙语
    ['zh-CN', { locale: zh_CN, fileLocales: {}, code: 'zh' }],
    ['zh-TW', { locale: zh_TW, fileLocales: {}, code: 'cht' }],
  ],
);

const srcLocaleNmae = 'zh-CN';

type TranslateResult = {
  from?: string;
  to?: string;
  trans_result?: [{ src: string; dst: string }];
  error_code?: string;
  error_msg: string;
};

// console.log(zhLang)
export default (api: IApi) => {
  const appid = process.env.BAIDU_TRANSLATE_APPID
  const appKey = process.env.BAIDU_TRANSLATE_APPKEY
  const crypto = require('crypto');
  const cryptoMD5 = (content: string): string => {
    const md5 = crypto.createHash('md5');
    md5.update(content);
    return md5.digest('hex');
  };
  const translate = (
    msg: string,
    options: { from: string; to: string },
  ): Promise<AxiosResponse<TranslateResult>> => {
    const salt = new Date().getTime();
    const str1 = appid + msg + salt + appKey;
    const sign = cryptoMD5(str1);
    const { from, to } = options;
    return axios.get('http://fanyi-api.baidu.com/api/trans/vip/translate', {
      params: {
        q: msg,
        appid: appid,
        salt: salt,
        from: from,
        to: to,
        sign: sign,
      },
    });
  };
  const sleep = (ms: number) => {
    return new Promise((resolve) => setTimeout(resolve, ms));
  };
  api.registerCommand({
    name: 'translate',
    fn: async () => {
      const srcLocaleCode = allLocales.get(srcLocaleNmae)?.code ?? 'auto';
      fs.readdirSync(`src/locales/`, { withFileTypes: true }).forEach((file) => {
        if (file.isDirectory()) {
          const locales = allLocales.get(file.name);
          if (locales) {
            fs.readdirSync(`src/locales/${file.name}/`).forEach((name) => {
              if (path.extname(name) === '.ts') {
                try {
                  const locale = require(`../../src/locales/${file.name}/${name}`);
                  locales.fileLocales = { ...locales.fileLocales, [name]: locale.default };
                } catch (error: any) {
                  if (
                    error.name === 'SyntaxError' &&
                    error.message === "Unexpected token 'export'"
                  ) {
                    throw Error(
                      `src/locales/${file.name}/${name} can't be import to src/locales/${file.name}`,
                    );
                  } else {
                    throw error;
                  }
                }
              }
            });
          }
        }
      });

      const srcLocales = new Map<string, { filename: string; msg: any }>();
      allLocales.forEach(({ fileLocales: fileLocals }, localeName) => {
        if (localeName === srcLocaleNmae) {
          for (const filename in fileLocals) {
            if (Object.prototype.hasOwnProperty.call(fileLocals, filename)) {
              const fileLocale = fileLocals[filename];
              // srcLocales.set(key,fileLocale)
              for (const key in fileLocale) {
                if (Object.prototype.hasOwnProperty.call(fileLocale, key)) {
                  srcLocales.set(key, { filename, msg: fileLocale[key] });
                }
              }
            }
          }
        }
      });

      for (const localeName of allLocales.keys()) {
        const { locale, code, fileLocales = {} } = allLocales.get(localeName) ?? {};
        if (localeName === srcLocaleNmae || !code) {
          continue;
        }
        for (const msgKey of srcLocales.keys()) {
          const { filename, msg } = srcLocales.get(msgKey) ?? {};
          if (!filename) {
            continue;
          }
          const fileLocale = fileLocales[filename] ?? {};
          if (!locale[msgKey]) {
            let dstMsg: string | undefined;
            // console.log(localeName, filename, msgKey, msg)
            while (true) {
              const res = await translate(msg, { from: srcLocaleCode, to: code });
              // console.log(res.data.trans_result);
              dstMsg = res.data.trans_result?.map(({ dst }) => dst).join('\n');
              if (!dstMsg) {
                switch (res.data.error_code) {
                  case '58001':
                    throw Error(`invalid dest code: ${code} (${JSON.stringify(res.data)})`);
                  case '52001':
                    console.log(`${res.data.error_msg}, wait 2s...`);
                    await sleep(2000);
                    continue;
                  case '54003':
                    console.log('API triggered rate limit, wait 10s...');
                    await sleep(10000);
                    continue;
                  default:
                    throw Error(
                      `tranlate failed: ${msgKey}(${msg}) ==> ${JSON.stringify(res.data)}`,
                    );
                }
              }
              break;
            }

            console.log(localeName, filename, msgKey, dstMsg);
            fileLocale[msgKey] = dstMsg;
            fileLocale._changed = true;
            fileLocales[filename] = fileLocale;
            await sleep(1000);
          }
        }
        for (const filename in fileLocales) {
          if (Object.prototype.hasOwnProperty.call(fileLocales, filename)) {
            const fileLocale = fileLocales[filename];
            if (fileLocale._changed) {
              delete fileLocale._changed;
              console.log(`save file: ${filename}`);
              const module_name = filename.substring(0, filename.lastIndexOf('.'));
              if (!fs.existsSync(`src/locales/${localeName}/${filename}`)) {
                let buf = fs.readFileSync(`src/locales/${localeName}.ts`, {
                  encoding: 'utf8',
                  flag: 'r',
                });
                buf = `import ${module_name} from './${localeName}/${module_name}';\n${buf}`;
                const lines = buf.split('\n');
                for (let index = 0; index < lines.length; index++) {
                  const line = lines[index];
                  if (line.match(/^\s*\.\.\./)) {
                    lines[index] = `${line}\n  ...${module_name},`;
                    break;
                  }
                }
                buf = lines.join('\n');
                fs.writeFileSync(`src/locales/${localeName}.ts`, buf);
              }
              let cnt = `export default ${JSON.stringify(fileLocale, null, 2)}`;
              if (cnt.slice(cnt.length - 2) === '\n}') {
                cnt = cnt.slice(0, cnt.length - 2) + ',\n}\n';
              }
              fs.writeFileSync(`src/locales/${localeName}/${filename}`, cnt);
            }
          }
        }
      }
    },
  });
};
