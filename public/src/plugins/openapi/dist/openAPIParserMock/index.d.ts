declare class OpenAPIGeneratorMockJs {
    protected openAPI: any;
    constructor(openAPI: any);
    sampleFromSchema: (schema: any, propsName?: string[]) => any;
    parser: () => any;
}
export default OpenAPIGeneratorMockJs;
