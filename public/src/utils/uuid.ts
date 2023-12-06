import { Worker } from 'snowflake-uuid';

class IdGenerator {
  private generator: Worker;
  private static instance: IdGenerator;

  private constructor() {
    this.generator = new Worker(0, 1, {
      workerIdBits: 5,
      datacenterIdBits: 5,
      sequenceBits: 12,
    });
  }
  static generate() {
    return this.getInstance().generator.nextId().toString();
  }
  private static getInstance() {
    // 静态方法中可以使用 this.静态属性，或者Dep.静态属性。
    if (this.instance) {
      return this.instance;
    }
    this.instance = new IdGenerator();
    return this.instance;
  }
}

export const newId = (): string => {
  return IdGenerator.generate();
};
