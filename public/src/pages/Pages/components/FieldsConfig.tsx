import { isArray } from 'lodash';
import React, { useState } from 'react';

import { PageFieldType } from '@/services/idas/enums';
import { enumToMap } from '@/utils/enum';
import type { IntlContext } from '@/utils/intl';
import type { ProColumns } from '@ant-design/pro-components';
import { ProField } from '@ant-design/pro-components';
import { EditableProTable } from '@ant-design/pro-components';
import { createField } from '@ant-design/pro-form/es/BaseForm';
import type { ProFormFieldItemProps } from '@ant-design/pro-form/lib/typing';

interface FieldsConfigProps extends ProFormFieldItemProps {
  parentIntl: IntlContext;
}

export interface FieldConfig extends API.FieldConfig {
  id: string;
  index?: number;
}

export default createField<FieldsConfigProps>(
  ({ parentIntl: intl, fieldProps }: FieldsConfigProps) => {
    const [editableKeys, setEditableRowKeys] = useState<React.Key[]>([]);
    return (
      <ProField
        mode="edit"
        fieldProps={fieldProps}
        renderFormItem={(value: FieldConfig[] | '', { onChange }) => {
          const columns: ProColumns<FieldConfig>[] = [
            {
              title: intl.t('field.name.title', 'Field Name'),
              dataIndex: 'name',
              formItemProps: {
                rules: [
                  {
                    required: true,
                    message: intl.t('field.name.require', 'This item is required!'),
                  },
                ],
              },
              width: '25%',
            },
            {
              title: intl.t('field.displayName.title', 'Display Name'),
              dataIndex: 'displayName',
              width: '25%',
            },
            {
              title: intl.t('field.fieldType.title', 'Type'),
              dataIndex: 'valueType',
              valueType: 'select',
              valueEnum: enumToMap(PageFieldType, intl, 'fieldType'),
              width: '20%',
            },
            {
              title: intl.t('field.option.title', 'Option'),
              valueType: 'option',
              width: '20%',
              render: (text, record, _, action) => [
                <a
                  key="editable"
                  onClick={() => {
                    action?.startEditable?.(record.id);
                  }}
                >
                  {intl.t('field.edit.button', 'Edit')}
                </a>,
                <a
                  key="delete"
                  onClick={() => {
                    onChange?.(
                      (isArray(value) ? value : []).filter((item) => item.name !== record.name),
                    );
                  }}
                >
                  {intl.t('field.delete.button', 'Delete')}
                </a>,
              ],
            },
          ];
          return (
            <EditableProTable<FieldConfig>
              rowKey="id"
              headerTitle={intl.t('fields.title', 'Fields configuration')}
              maxLength={100}
              recordCreatorProps={{
                position: 'bottom',
                record: (idx, data) => {
                  const newData = {
                    name: 'New Field',
                    id: `${new Date().getTime()}${(Math.random() * 3).toFixed(0).toString()}`,
                    valueType: PageFieldType.text,
                  };
                  if (data.find((item) => item.name == newData.name)) {
                    newData.name = `New Field ${idx}`;
                  }
                  for (let index = 0; index < data.length; index++) {
                    const element = data[index];
                    if (element.id == newData.id) {
                      console.error(new Error('Unknown error'));
                    }
                  }
                  return newData;
                },
                newRecordType: 'dataSource',
              }}
              loading={false}
              toolBarRender={() => []}
              columns={columns}
              expandedRowRender={(record) => {
                return <>{record}</>;
              }}
              value={(isArray(value) ? value : []).map((item) => ({
                ...item,
                id: item.id ?? `id-${item.name}`,
              }))}
              onChange={onChange}
              controlled
              editable={{
                type: 'multiple',
                editableKeys,
                // onValuesChange: (record, recordList) => {
                //   setDataSource(recordList);
                // },
                onChange: setEditableRowKeys,
              }}
            />
          );
        }}
      />
    );
  },
);
