import type { OpenAPIObject, OperationObject, SchemaObject } from 'openapi3-ts';
export declare type GenerateServiceProps = {
    requestLibPath?: string;
    requestImportStatement?: string;
    /**
     * api 的前缀
     */
    apiPrefix?: string | ((params: {
        path: string;
        method: string;
        namespace: string;
        functionName: string;
        autoExclude?: boolean;
    }) => string);
    /**
     * 生成的文件夹的路径
     */
    serversPath?: string;
    /**
     * Swagger 2.0 或 OpenAPI 3.0 的地址
     */
    schemaPath?: string;
    /**
     * 项目名称
     */
    projectName?: string;
    hook?: {
        /** change open api data after constructor */
        afterOpenApiDataInited?: (openAPIData: OpenAPIObject) => OpenAPIObject;
        /** 自定义函数名称 */
        customFunctionName?: (data: OperationObject) => string;
        /** 自定义类型名称 */
        customTypeName?: (data: OperationObject) => string;
        /** 自定义类名 */
        customClassName?: (tagName: string) => string;
        /**
         * 自定义获取type hook
         * 返回非字符串将使用默认方法获取type
         * @example set number to string
         * function customType(schemaObject,namespace){
         *  if(schemaObject.type==='number' && !schemaObject.format){
         *    return 'BigDecimalString';
         *  }
         * }
         */
        customType?: (schemaObject: SchemaObject | undefined, namespace: string, originGetType: (schemaObject: SchemaObject | undefined, namespace: string) => string) => string;
        /**
         * 自定义生成文件名，可返回多个，表示生成多个文件
         * 返回为空，则使用默认的获取方法获取
         * @example  使用operationId生成文件名
         * function customFileNames(operationObject,apiPath){
         *   const operationId=operationObject.operationId;
         *   if (!operationId) {
         *      console.warn('[Warning] no operationId', apiPath);
         *      return;
         *    }
         *    const res = operationId.split('_');
         *    if (res.length > 1) {
         *      res.shift();
         *      if (res.length > 2) {
         *        console.warn('[Warning]  operationId has more than 2 part', apiPath);
         *      }
         *      return [res.join('_')];
         *    } else {
         *      const controllerName = (res || [])[0];
         *      if (controllerName) {
         *        return [controllerName];
         *      }
         *      return;
         *    }
         * }
         */
        customFileNames?: (operationObject: OperationObject, apiPath: string, _apiMethod: string) => string[];
    };
    namespace?: string;
    /**
     * 默认为false，true时使用null代替可选
     */
    nullable?: boolean;
    mockFolder?: string;
    /**
     * 模板文件的文件路径
     */
    templatesFolder?: string;
    /**
     * 枚举样式
     */
    enumStyle?: 'string-literal' | 'enum';
    /**
     * response中数据字段
     * example: ['result', 'res']
     */
    dataFields?: string[];
};
export declare const getSchema: (schemaPath: string) => Promise<any>;
export declare const generateService: ({ requestLibPath, schemaPath, mockFolder, nullable, ...rest }: GenerateServiceProps) => Promise<void>;
