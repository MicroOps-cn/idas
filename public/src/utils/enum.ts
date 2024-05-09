import type { PrimitiveType } from 'intl-messageformat';
import type { MessageDescriptor } from 'react-intl';

export type LabelValue = { label: string; value: string | number; key: string };

interface IntlShape {
  formatMessage: (descriptor: MessageDescriptor, values?: Record<string, PrimitiveType>) => string;
}

export const enumToOptions = (
  enumObj: any,
  intl: IntlShape,
  intlIdPrefix: string,
  filter?: (item: string, element: any) => boolean,
): LabelValue[] => {
  const opts: LabelValue[] = [];
  for (const key in enumObj) {
    if (Object.prototype.propertyIsEnumerable.call(enumObj, key) && isNaN(Number(key))) {
      const element = enumObj[key];
      if (filter && !filter(key, element)) {
        continue;
      }
      opts.push({
        label: intl.formatMessage({
          id: `${intlIdPrefix}.${key}`,
          defaultMessage: key,
        }),
        key,
        value: element,
      });
    }
  }
  return opts;
};

type EnumMap = Map<string, string | number>;

export const enumToMap = (
  enumObj: any,
  intl: IntlShape,
  intlIdPrefix: string,
  filter?: (item: string, element: any) => boolean,
): EnumMap => {
  const opts: EnumMap = new Map();
  for (const key in enumObj) {
    if (Object.prototype.hasOwnProperty.call(enumObj, key) && isNaN(Number(key))) {
      const element = enumObj[key];
      if (filter && !filter(key, element)) {
        continue;
      }
      opts.set(
        element,
        intl.formatMessage({
          id: `${intlIdPrefix}.${key}`,
          defaultMessage: key,
        }),
      );
    }
  }
  return opts;
};

type Status = 'Success' | 'Error' | 'Processing' | 'Warning' | 'Default';

type IValueEnum = Record<
  string,
  | React.ReactNode
  | {
    text: React.ReactNode;
    status: Status;
  }
>;

export const enumToStatusEnum = (
  enumObj: any,
  intl: IntlShape,
  intlIdPrefix: string,
  valueStatsuMap: Record<string, Status>,
): IValueEnum => {
  const opts: IValueEnum = {};
  for (const key in enumObj) {
    if (Object.prototype.hasOwnProperty.call(enumObj, key) && !isNaN(Number(key))) {
      const element = enumObj[key];
      opts[key] = {
        text: intl.formatMessage({
          id: `${intlIdPrefix}.${element}`,
          defaultMessage: element,
        }),
        status: valueStatsuMap[element] ?? 'Default',
      };
    }
  }
  return opts;
};
