export declare const getAbsolutePath: (filePath: string) => string;
export declare const mkdir: (dir: string) => void;
export declare const prettierFile: (content: string) => [string, boolean];
export declare const writeFile: (folderPath: string, fileName: string, content: string) => boolean;
export declare const getTagName: (name: string) => string;
/**
 * 根据当前的数据源类型，对请求回来的 apiInfo 进行格式化
 * 如果是 op 数据源，对 tags 以及 path 中的 tags 进行处理
 * - before: 前缀（产品集.产品码） + 操作对象（必填）+ 子操作对象（可选）+ 动作（必填）
 * - after: 操作对象（必填）+ 子操作对象（可选） ==> 驼峰
 */
export declare const formatApiInfo: (apiInfo: Record<string, any>) => any;
declare type serviceParam = {
    title: string;
    type: string;
    description: string;
    default: string;
    [key: string]: any;
};
declare type serviceParams = Record<string, serviceParam>;
/**
 * 一方化场景下，由于 onex 会对请求的响应做处理
 *  1. 将 Response & Request 中的参数字段会变更为小驼峰写法
 *  onex 相关代码 ： http://gitlab.alipay-inc.com/one-console/sdk/blob/master/src/request.ts#L110
 *  2. 另外要注意：
 *  op 返回的数据，请求参数的类型格式 需要做额外的处理
 *  - (name) key.n, (type) string  ==> key: string []
 *  - (name) key.m,  (type) string ===>  key: string []
 *  - (name) key.key1 , (type) string ==> key: {key1:string}
 *  - (name) key.n.key1 ,(type) string => key:{ key1 :string}[]
 *  - (name) key.n.key1.m,(type) string ==> key:{key1: string[]}[]
 */
export declare function formatParamsForYFH(params: serviceParams, paramsObject?: serviceParams): serviceParams;
export declare const stripDot: (str: string) => string;
export {};
