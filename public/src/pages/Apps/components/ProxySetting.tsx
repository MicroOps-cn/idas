import type { FormInstance, FormItemProps, InputRef } from 'antd';
import { Button, Form, Input, InputNumber, Select, Table } from 'antd';
import type { StoreValue } from 'antd/lib/form/interface';
import type { DefaultOptionType } from 'antd/lib/select';
import type { ColumnsType } from 'antd/lib/table';
import type { ColumnTitle } from 'antd/lib/table/interface';
import { arrayMoveImmutable } from 'array-move';
import { isArray, isString } from 'lodash';
import type { BaseSelectRef } from 'rc-select';
import type { Dispatch, SetStateAction } from 'react';
import React, { useContext, useRef, useEffect, useState, useMemo } from 'react';
import type { SortableContainerProps, SortEnd } from 'react-sortable-hoc';
import { SortableContainer, SortableElement, SortableHandle } from 'react-sortable-hoc';

import { IntlContext } from '@/utils/intl';
import { newId } from '@/utils/uuid';
import { DeleteOutlined, MenuOutlined, PlusOutlined } from '@ant-design/icons';
import { ProFormCheckbox, ProFormField, ProFormGroup, ProFormText } from '@ant-design/pro-form';

import styles from '../style.less';

interface ProxyConfigProps {
  parentIntl: IntlContext;
  dataSource: API.AppProxyInfo & { id?: string };
  setDataSource: Dispatch<SetStateAction<API.AppProxyInfo>>;
}

const DragHandle = SortableHandle(() => <MenuOutlined style={{ cursor: 'grab', color: '#999' }} />);

const EditableContext = React.createContext<FormInstance<any> | null>(null);

const SortableItem = SortableElement((props: React.HTMLAttributes<HTMLTableRowElement>) => {
  const [form] = Form.useForm();
  return (
    <Form form={form} component={false}>
      <EditableContext.Provider value={form}>
        <tr className="editable-row" {...props} />
      </EditableContext.Provider>
    </Form>
  );
});

const SortableBody = SortableContainer((props: React.HTMLAttributes<HTMLTableSectionElement>) => (
  <tbody {...props} />
));

const methodOptions: DefaultOptionType[] = [
  'GET',
  'POST',
  'PUT',
  'PATCH',
  'DELETE',
  'OPTIONS',
  'HEAD',
].map((value) => ({
  value,
  label: value,
}));

const isNullNode = (node: React.ReactNode): boolean => {
  if (isArray(node)) {
    for (const elem of node) {
      if (!isNullNode(elem)) return false;
    }
    return true;
  } else if (isString(node)) {
    return !Boolean(node.trim());
  }
  return !Boolean(node);
};

interface EditableCellProps<Item = Record<string, any>> {
  title?: ColumnTitle<Item>;
  editable?: boolean;
  children: React.ReactNode;
  dataIndex: string | number;
  record: Item;
  inputType?: 'select' | 'text' | 'number';
  handleSave: (record: Item) => void;
  options?: DefaultOptionType[];
  formProps: FormItemProps;
}
type IRef = InputRef & BaseSelectRef & HTMLInputElement;

const EditableCell: React.FC<EditableCellProps> = ({
  title,
  editable,
  children,
  dataIndex,
  record,
  handleSave,
  options,
  inputType,
  formProps,
  ...restProps
}) => {
  const [editing, setEditing] = useState(false);
  const inputRef = useRef<IRef>(null);
  const form = useContext(EditableContext)!;
  useEffect(() => {
    if (editing) {
      inputRef.current!.focus();
    }
  }, [editing]);

  const toggleEdit = () => {
    setEditing(!editing);
    form.setFieldsValue({ [dataIndex]: record[dataIndex] });
  };

  const save = async () => {
    try {
      const values = await form.validateFields();
      toggleEdit();
      handleSave({ ...record, ...values });
    } catch (errInfo) {
      console.error('Save failed:', errInfo);
    }
  };
  const inputNode = (() => {
    switch (inputType) {
      case 'number':
        return <InputNumber ref={inputRef} size="small" onPressEnter={save} onBlur={save} />;
      case 'select':
        return <Select ref={inputRef} size="small" onBlur={save} options={options} />;
      default:
        return <Input ref={inputRef} size="small" onPressEnter={save} onBlur={save} />;
    }
  })();
  let childNode = children;
  if (editable) {
    childNode = editing ? (
      <Form.Item
        className={styles.ProxyUrlInput}
        style={{ margin: 0 }}
        name={dataIndex}
        rules={[{ required: true, message: '' }]}
        {...formProps}
      >
        {inputNode}
      </Form.Item>
    ) : (
      <div className="editable-cell-value-wrap" style={{ paddingRight: 24 }} onClick={toggleEdit}>
        {isNullNode(children) ? <>&nbsp;</> : children}
      </div>
    );
  }

  return <td {...restProps}>{childNode}</td>;
};

type ColumnType<T> = ColumnsType<T>[number] & {
  editable?: boolean;
  inputType?: 'select' | 'text' | 'number';
  onCell?: (record: API.AppProxyUrl) => Exclude<EditableCellProps, { children: React.ReactNode }>;
  dataIndex?: string;
  options?: DefaultOptionType[];
  formProps?: FormItemProps;
};

const ProxySetting: React.FC<ProxyConfigProps> = ({ parentIntl, dataSource, setDataSource }) => {
  const [tableHeight, setTableHeight] = useState<number>(240);
  const {
    id,
    domain,
    urls,
    upstream: initUpstream,
    insecureSkipVerify,
    transparentServerName,
    jwtProvider: initJwtProvider,
    jwtSecret,
    jwtCookieName,
    hstsOffload,
  } = dataSource;
  const [urlMaxIndex, setUrlMaxIndex] = useState<number>(0);
  const [upstream, setUpstream] = useState<string>(initUpstream);
  const [jwtProvider, setJwtProvider] = useState<boolean>(initJwtProvider);
  const [errorMessage, setErrorMessage] = useState<React.ReactNode>();
  const intl = useMemo(() => {
    return new IntlContext('proxy', parentIntl);
  }, [parentIntl]);
  useEffect(() => {
    const resetTableHeight = () => {
      if (window.innerHeight > 890) {
        setTableHeight(window.innerHeight - 690);
      } else {
        setTableHeight(200);
      }
    };
    resetTableHeight();
    window.addEventListener('resize', resetTableHeight);
  }, []);

  useEffect(() => {
    let errMsg: React.ReactNode = '';
    for (const url of urls) {
      if (!url.name || !url.name.trim()) {
        errMsg = intl.t('name.required', 'name cannot be empty!');
      }
      if (!url.method || !url.method.trim()) {
        errMsg = intl.t('method.required', 'method cannot be empty!');
      }
      if (!url.url || !url.url.trim()) {
        errMsg = intl.t('url.required', 'URL cannot be empty!');
      }
    }
    setErrorMessage(errMsg);
  }, [intl, urls]);

  const onSortEnd = ({ oldIndex, newIndex }: SortEnd) => {
    if (oldIndex !== newIndex) {
      const newUrls = arrayMoveImmutable(urls.slice(), oldIndex, newIndex).filter(
        (el: API.AppProxyUrl) => !!el,
      );
      setDataSource({ ...dataSource, urls: newUrls });
    }
  };

  const DraggableContainer: React.FC<SortableContainerProps> = (props: SortableContainerProps) => (
    <SortableBody
      useDragHandle
      disableAutoscroll
      helperClass="row-dragging"
      onSortEnd={onSortEnd}
      {...props}
    />
  );

  const DraggableBodyRow: React.FC<any> = ({ className, style, ...restProps }) => {
    // function findIndex base on Table rowKey props and should always be a right array index
    const index = urls.findIndex((x) => x.id === restProps['data-row-key']);
    return <SortableItem index={index} {...restProps} />;
  };

  const setUrls = (record: API.AppProxyUrl) => {
    setDataSource((ori) => {
      const newUrls =
        ori?.urls.map((item) => {
          if (item.id === record.id) {
            return record;
          }
          return item;
        }) ?? [];
      return { ...ori, urls: newUrls };
    });
  };

  const columns: ColumnType<API.AppProxyUrl>[] = [
    {
      dataIndex: 'sort',
      width: 30,
      className: 'drag-visible',
      render: () => <DragHandle />,
    },
    {
      width: 150,
      title: intl.t('name.title', 'Name'),
      dataIndex: 'name',
      inputType: 'text',
      editable: true,
    },
    {
      title: intl.t('method.title', 'Method'),
      dataIndex: 'method',
      inputType: 'select',
      width: 100,
      options: methodOptions,
      editable: true,
    },
    {
      title: intl.t('url.title', 'URL'),
      dataIndex: 'url',
      inputType: 'text',
      editable: true,
    },
    {
      title: intl.t('url.upstream', 'Upstream'),
      dataIndex: 'upstream',
      inputType: 'text',
      editable: true,
      formProps: {
        required: false,
        rules: [],
      },
      render: (value: any) => {
        return value ?? upstream;
      },
    },
    {
      width: 30,
      render: (_, record) => [
        <a
          key="delete"
          onClick={async () => {
            setDataSource((ori) => ({
              ...ori,
              urls: ori.urls.filter((item) => item.id !== record.id),
            }));
          }}
        >
          <DeleteOutlined />
        </a>,
      ],
    },
  ];
  const editableColumns: ColumnType<API.AppProxyUrl>[] = columns.map(
    ({ formProps, ...col }: ColumnType<API.AppProxyUrl>) => {
      if (!col.editable) {
        return col;
      }
      return {
        ...col,
        onCell: (record: API.AppProxyUrl) => ({
          record,
          editable: col.editable,
          dataIndex: col.dataIndex,
          inputType: col.inputType,
          options: col.options,
          title: col.title,
          handleSave: setUrls,
          formProps: formProps,
        }),
      } as ColumnType<API.AppProxyUrl>;
    },
  );

  return (
    <>
      <ProFormText
        name={['proxy', 'domain']}
        label={intl.t('domain.label', 'Domain')}
        colProps={{ span: 12 }}
        initialValue={domain}
        rules={[
          {
            required: true,
            message: intl.t('domain.required', 'domain cannot be empty!'),
          },
        ]}
      />
      <ProFormText
        name={['proxy', 'upstream']}
        label={intl.t('upstream.label', 'Upstream')}
        colProps={{ span: 12 }}
        initialValue={upstream}
        fieldProps={{
          onChange: (e) => {
            setUpstream(e.target.value);
          },
        }}
        tooltip={
          <>
            {intl.t('upstream.example', 'Example')}:<br />
            <li>abc.com</li>
            <li>http://abc.com:80</li>
            <li>http://1.2.3.4</li>
            <li>https://abc.com</li>
          </>
        }
        rules={[
          {
            required: true,
            message: intl.t('upstream.required', 'upstream cannot be empty!'),
          },
        ]}
      />
      <ProFormCheckbox
        name={['proxy', 'insecureSkipVerify']}
        label={intl.t('insecureSkipVerify.label', 'Skip TLS Verify')}
        colProps={{ span: 8 }}
        initialValue={insecureSkipVerify}
        tooltip={intl.t(
          'insecureSkipVerify.describe',
          'When requesting back-end servers, the certificate verification is ignored (insecure).',
        )}
      />
      <ProFormCheckbox
        name={['proxy', 'transparentServerName']}
        label={intl.t('transparentServerName.label', 'Transparent Server Name')}
        colProps={{ span: 8 }}
        initialValue={transparentServerName}
        tooltip={intl.t(
          'transparentServerName.describe',
          'When requesting the backend, the domain name requested by the client will be transparently transmitted.',
        )}
      />
      <ProFormCheckbox
        name={['proxy', 'hstsOffload']}
        label={intl.t('hstsOffload.label', 'HSTS Offload')}
        colProps={{ span: 8 }}
        initialValue={hstsOffload}
        tooltip={intl.t('hstsOffload.describe', 'Delete HSTS field in response header')}
      />
      <ProFormGroup colProps={{ span: 24 }}>
        <ProFormCheckbox
          name={['proxy', 'jwtProvider']}
          label={intl.t('jwtProvider.label', 'JWT provider')}
          colProps={{ span: 4 }}
          initialValue={jwtProvider}
          fieldProps={{
            onChange: (e) => {
              setJwtProvider(e.target.checked);
            },
          }}
          tooltip={intl.t('jwtProvider.describe', 'As a JWT provider, issue tokens to clients.')}
        />
        <ProFormText
          name={['proxy', 'jwtCookieName']}
          label={intl.t('jwtCookieName.label', 'JWT Cookie Name')}
          colProps={{ span: 4 }}
          initialValue={jwtCookieName}
          disabled={!jwtProvider}
          rules={[
            {
              validator: async (_, value: string) => {
                if (jwtProvider && !value) {
                  return Promise.reject();
                }
                return Promise.resolve();
              },
            },
          ]}
          tooltip={intl.t('jwtCookieName.describe', 'Store the name of the JWT Cookie.')}
        />
        <ProFormText
          name={['proxy', 'jwtSecret']}
          label={intl.t('jwtSecret.label', 'JWT secret')}
          colProps={{ span: 12 }}
          initialValue={jwtSecret}
          fieldProps={{ maxLength: 200 }}
          placeholder={id ? '*******************************' : undefined}
          rules={[
            {
              pattern:
                /^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{4}|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)$/,
              message: intl.t('proxy.jwtSecret-message', 'Please enter a valid Base64 encoding.'),
            },
            {
              validator: async (rule, value: StoreValue) => {
                if (jwtProvider && !value && !id) {
                  return Promise.reject();
                }
                return Promise.resolve();
              },
              message: intl.t('proxy.jwtSecret-message', 'Please enter a valid Base64 encoding.'),
            },
          ]}
          disabled={!jwtProvider}
          tooltip={intl.t('jwtSecret.describe', 'Use this key to issue tokens.(Base64 encoding)')}
        />
      </ProFormGroup>
      <ProFormField colProps={{ span: 24 }} label={intl.t('url.label', 'URL')}>
        <Table
          pagination={false}
          dataSource={urls}
          columns={editableColumns}
          scroll={{ y: tableHeight }}
          rowKey="id"
          // rowClassName={() => 'editable-row'}
          size="small"
          components={{
            body: {
              wrapper: DraggableContainer,
              row: DraggableBodyRow,
              cell: EditableCell,
            },
          }}
        />

        <Button
          onClick={() => {
            setDataSource({
              ...dataSource,
              urls: [...urls, { name: '', id: newId(), method: '*', url: '' }],
            });
            setUrlMaxIndex(urlMaxIndex + 1);
          }}
          type="dashed"
          block
        >
          <PlusOutlined />
        </Button>
        <div style={{ color: 'red' }}>{errorMessage}</div>
      </ProFormField>
    </>
  );
};
export default ProxySetting;
