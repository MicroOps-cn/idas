import { defaults } from 'lodash';
import minimist from 'minimist';
import path from 'path';
import type { IApi } from 'umi';

interface CommandArgs {
  basePath: string;
  apiPath: string;
  publicPath?: string;
}

const commandArgs: CommandArgs = defaults(minimist(process.argv.slice(2)), {
  basePath: '/',
  apiPath: '/',
});
const { apiPath, basePath, publicPath = basePath } = commandArgs;

const { REACT_APP_ENV = 'dev' } = process.env;



export default (api: IApi) => {
  api.addHTMLLinks(() => [{ rel: 'icon', href: path.join(basePath, publicPath, 'logo.svg') }]);
  api.onStart(() => {
    console.log(`basePath: ${basePath}, apiPath: ${apiPath}, publicPath: ${publicPath}`);
    console.log(`env: ${REACT_APP_ENV}`)
  })
};
