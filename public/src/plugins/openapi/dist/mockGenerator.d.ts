export declare type genMockDataServerConfig = {
    openAPI: any;
    mockFolder: string;
};
declare const mockGenerator: ({ openAPI, mockFolder }: genMockDataServerConfig) => Promise<void>;
export { mockGenerator };
