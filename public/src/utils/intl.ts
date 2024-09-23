import type { PrimitiveType } from 'intl-messageformat';
import type { IntlShape, MessageDescriptor, IntlFormatters, FormatDateOptions, FormatNumberOptions, FormatPluralOptions, FormatRelativeTimeOptions } from 'react-intl';

import locales from '@/locales/zh-CN';
import { DisplayNamesOptions } from '@formatjs/intl-displaynames/lib';
import { IntlListFormatOptions } from '@formatjs/intl-listformat';
import { FormattableUnit } from '@formatjs/intl-relativetimeformat';
import { ReactNode } from 'react';

export class IntlContext implements IntlFormatters {
  private id: string;
  private parent: IntlContext | IntlShape;
  locale: string;
  constructor(id: string, parent: IntlContext | IntlShape) {
    this.id = id;
    this.parent = parent;
    this.locale = parent.locale;
  }
  formatDate(value: string | number | Date | undefined, opts?: FormatDateOptions | undefined): string {
    return this.parent.formatDate(value, opts);
  }
  formatTime(value: string | number | Date | undefined, opts?: FormatDateOptions | undefined): string {
    return this.parent.formatTime(value, opts);
  }
  formatDateToParts(value: string | number | Date | undefined, opts?: FormatDateOptions | undefined): Intl.DateTimeFormatPart[] {
    return this.parent.formatDateToParts(value, opts);
  }
  formatTimeToParts(value: string | number | Date | undefined, opts?: FormatDateOptions | undefined): Intl.DateTimeFormatPart[] {
    return this.parent.formatTimeToParts(value, opts);
  }
  formatRelativeTime(value: number, unit?: FormattableUnit | undefined, opts?: FormatRelativeTimeOptions | undefined): string {
    return this.parent.formatRelativeTime(value, unit, opts);
  }
  formatNumber(value: number | bigint, opts?: FormatNumberOptions | undefined): string {
    return this.parent.formatNumber(value, opts);
  }
  formatNumberToParts(value: number | bigint, opts?: FormatNumberOptions | undefined): Intl.NumberFormatPart[] {
    return this.parent.formatNumberToParts(value, opts);
  }
  formatPlural(value: number, opts?: FormatPluralOptions | undefined): Intl.LDMLPluralRule {
    return this.parent.formatPlural(value, opts);
  }
  formatList(values: string[], opts?: IntlListFormatOptions | undefined): string;
  formatList(values: ReactNode[], opts?: IntlListFormatOptions | undefined): React.ReactNode;
  formatList(values: ReactNode[] | string[], opts?: IntlListFormatOptions | undefined): React.ReactNode {
    return this.parent.formatList(values, opts);
  }
  formatDisplayName(value: string | number | object, opts?: DisplayNamesOptions | undefined): string | undefined {
    return this.parent.formatDisplayName(value, opts);
  }
  formatHTMLMessage(descriptor: MessageDescriptor, values?: Record<string, PrimitiveType>): React.ReactNode {
    const id = this.getId(descriptor.id);
    if (!(this.parent instanceof IntlContext) && id) {
      if (!(locales as Record<string, string>)[id])
        console.log(`'${id}': '${descriptor.defaultMessage}',`);
    }
    return this.parent.formatHTMLMessage({ ...descriptor, id: id }, values);
  }
  formatMessage(descriptor: MessageDescriptor, values?: Record<string, PrimitiveType>): string {
    const id = this.getId(descriptor.id);
    if (!(this.parent instanceof IntlContext) && id) {
      if (!(locales as Record<string, string>)[id])
        console.log(`'${id}': '${descriptor.defaultMessage}',`);
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
