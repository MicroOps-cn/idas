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
const { basePath = '/', publicPath = '/' } = commandArgs;

export default (api: IApi) => {
  api.addHTMLLinks(() => [{ rel: 'icon', href: path.join(basePath, publicPath, 'logo.svg') }]);
};
