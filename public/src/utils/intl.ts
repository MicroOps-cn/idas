import type { PrimitiveType } from 'intl-messageformat';
import type { IntlShape, MessageDescriptor } from 'react-intl';

import locales from '@/locales/zh-CN';

export class IntlContext {
  private id: string;
  private parent: IntlContext | IntlShape;
  constructor(id: string, parent: IntlContext | IntlShape) {
    this.id = id;
    this.parent = parent;
    this.locale = parent.locale;
  }
  locale: string;
  formatMessage(descriptor: MessageDescriptor, values?: Record<string, PrimitiveType>): string {
    const id = this.getId(descriptor.id);
    if (!(this.parent instanceof IntlContext) && id) {
      if (!(locales as Record<string, string>)[id])
        console.log(`"${id}": "${descriptor.defaultMessage}",`);
    }
    return this.parent.formatMessage({ ...descriptor, id: id }, values);
  }
  private getId(id?: string | number): string | undefined | number {
    if (!this.id) return id ?? undefined;
    return id ? `${this.id}.${id}` : undefined;
  }
  t(
    id: string | number,
    defaultMessage?: string,
    description?: string | object,
    values?: Record<string, PrimitiveType>,
  ): string {
    return this.formatMessage({ id, defaultMessage, description }, values);
  }
}
