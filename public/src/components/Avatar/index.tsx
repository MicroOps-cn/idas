import type { AvatarProps as AntdAvatarProps } from 'antd';
import { Divider, List, Modal, Popover, Skeleton } from 'antd';
import { Upload } from 'antd';
import { Avatar as AntdAvatar } from 'antd';
import ImgCrop from 'antd-img-crop';
import type { RcFile, UploadFile } from 'antd/lib/upload';
import { isString } from 'lodash';
import { useEffect, useState } from 'react';
import InfiniteScroll from 'react-infinite-scroll-component';

import { uploadFile as postFile } from '@/services/idas/files';
import type { RequestError } from '@/utils/request';
import { UploadOutlined } from '@ant-design/icons';
import { ProField } from '@ant-design/pro-components';
import type {
  ProFormFieldProps as ProFormFieldItemProps,
  RequestData,
} from '@ant-design/pro-components';
import { createField } from '@ant-design/pro-form/es/BaseForm/createField';
import { ClassNames } from '@emotion/react';

export type AvatarProps = {
  src?: string;
} & Omit<AntdAvatarProps, 'src'>;

declare const apiPath: string;

export const getAvatarSrc = (src?: string) => {
  if (src && src.match(/^[-_a-zA-Z0-9]+$/)) {
    if (apiPath.endsWith('/')) {
      return apiPath + `api/v1/files/${src}`;
    }
    return apiPath + `/api/v1/files/${src}`;
  }
  return src;
};

const Avatar: React.FC<AvatarProps> = ({ src, ...props }) => {
  return <AntdAvatar src={getAvatarSrc(src)} {...props} />;
};

interface AvatarUploadProps {
  onError?: (error: any) => void;
  onChange?: (url?: string) => void;
  value?: string;
}
const getBase64 = (file: RcFile): Promise<string> =>
  new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result as string);
    reader.onerror = (error) => reject(error);
  });

const handleUploadFile = async (
  filename: string,
  fileObj: RcFile | string | Blob,
  onError?: (err: any) => void,
): Promise<string> => {
  try {
    const formData = new FormData();
    formData.append(filename, fileObj);
    const resp = await postFile({ data: formData, requestType: 'form' });
    if (resp.data) {
      return resp.data[filename];
    }
  } catch (error) {
    onError?.(error);
  }
  return '';
};

export const AvatarUpload: React.FC<AvatarUploadProps> = ({
  onError,
  onChange,
  value,
  ...props
}: AvatarUploadProps) => {
  const [avatar, setAvatar] = useState<UploadFile>();

  const [previewOpen, setPreviewOpen] = useState(false);
  const [previewImage, setPreviewImage] = useState('');
  const handleCancel = () => setPreviewOpen(false);

  const handlePreview = async (file: UploadFile) => {
    if (!file.url && !file.preview) {
      file.preview = await getBase64(file.originFileObj as RcFile);
    }
    setPreviewImage(file.url || (file.preview as string));
    setPreviewOpen(true);
  };
  useEffect(() => {
    if (value) {
      setAvatar({
        uid: value,
        name: 'img',
        url: getAvatarSrc(value),
      });
    } else {
      setAvatar(undefined);
    }
  }, [value]);
  return (
    <>
      <ImgCrop
        beforeCrop={async (file: RcFile): Promise<boolean> => {
          if (file.type === 'image/svg+xml') {
            const fileId = await handleUploadFile(file.name, file, onError);
            if (fileId) {
              setAvatar({
                uid: fileId,
                name: file.name,
                url: getAvatarSrc(fileId),
              });
              onChange?.(fileId);
            }
            return false;
          }
          return true;
        }}
      >
        <Upload
          // showUploadList={false}
          accept="image/png, image/jpeg, image/svg+xml"
          fileList={avatar ? [avatar] : []}
          maxCount={1}
          listType="picture-card"
          onPreview={handlePreview}
          onRemove={() => {
            setAvatar(undefined);
            onChange?.();
          }}
          customRequest={async (options) => {
            try {
              let { filename } = options;
              const { file } = options;
              if (!isString(file) && file) {
                filename = (file as RcFile).name ?? filename;
              }
              filename = filename ?? new Date().getTime().toString();
              filename = filename.substring(0, filename.lastIndexOf('.')) + '.png';
              const fileId = await handleUploadFile(filename, options.file, onError);
              if (fileId) {
                setAvatar({
                  uid: fileId,
                  name: filename,
                  url: getAvatarSrc(fileId),
                });
                onChange?.(fileId);
              }
            } catch (error) {
              onError?.(error);
            }
          }}
          {...props}
        >
          {avatar ? undefined : <UploadOutlined />}
        </Upload>
      </ImgCrop>
      <Modal open={previewOpen} footer={null} onCancel={handleCancel}>
        <img alt="example" style={{ width: '100%' }} src={previewImage} />
      </Modal>
    </>
  );
};

interface ProAvatarUploadProps extends ProFormFieldItemProps {
  onError?: (error: any) => void;
  optionsRequest?: (
    params: {
      pageSize?: number;
      current?: number;
      keywords?: string;
    },
    props: any,
  ) => Promise<RequestData<{ id: string }>>;
}

export const ProAvatarUpload: React.FC<ProFormFieldItemProps<ProAvatarUploadProps>> = ({
  fieldProps: { onError, ...fieldProps } = { onError: undefined },
  proFieldProps,
  optionsRequest,
}: ProAvatarUploadProps) => {
  const [avatar, setAvatar] = useState<UploadFile>();

  const [previewOpen, setPreviewOpen] = useState(false);
  const [previewImage, setPreviewImage] = useState('');
  const handleCancel = () => setPreviewOpen(false);

  const handlePreview = async (file: UploadFile) => {
    if (!file.url && !file.preview) {
      file.preview = await getBase64(file.originFileObj as RcFile);
    }
    setPreviewImage(file.url || (file.preview as string));
    setPreviewOpen(true);
  };
  const [iconList, setIconList] = useState<{ id: string }[]>([]);

  const [iconListHasMore, setIconListHasMore] = useState<boolean>(true);
  const [iconListPageNumber, setIconListPageNumber] = useState<number>(0);
  const [loadingIconList, setLoadingIconList] = useState<boolean>(false);
  const [popupVisible, setPopupVisible] = useState<boolean>(false);
  const loadMoreIconList = async (params?: API.getAppIconsParams) => {
    if (optionsRequest) {
      try {
        setLoadingIconList(true);
        const resp = await optionsRequest(
          {
            current: iconListPageNumber + 1,
            pageSize: 40,
            ...params,
          },
          {},
        );
        if (resp && resp.data) {
          const { data: newData, current, pageSize, total } = resp;
          setIconList((oldData) => {
            if (current === 1) {
              return newData;
            }
            return [...oldData, ...newData];
          });
          setIconListPageNumber(current);
          if (total && total < current * pageSize) {
            setIconListHasMore(false);
          } else {
            setIconListHasMore(true);
          }
        } else {
          setIconListHasMore(false);
        }
      } catch (error) {
        if (!(error as RequestError).handled) {
          console.error(`failed to get user list: ${error}`);
        }
      } finally {
        setLoadingIconList(false);
      }
    }
  };
  // const setPopupVisible = (v: boolean) => {
  //   if (!v || optionsRequest) {
  //     if (v) {
  //       loadMoreIconList();
  //     } else {
  //       setIconListHasMore(true);
  //       setIconListPageNumber(0);
  //     }
  //   }
  //   _setPopupVisible(v);
  // };
  return (
    <>
      <ProField
        mode="edit"
        fieldProps={fieldProps}
        renderFormItem={(text, { onChange, mode, value, ...props }) => {
          const avatarList = text ? [{ uid: value, name: 'img', url: getAvatarSrc(text) }] : [];
          return mode === 'edit' ? (
            <>
              <div style={{ height: 112, width: 112 }}>
                <ImgCrop
                  beforeCrop={async (file: RcFile): Promise<boolean> => {
                    if (file.type === 'image/svg+xml') {
                      const fileId = await handleUploadFile(file.name, file, onError);
                      if (fileId) {
                        setAvatar({
                          uid: fileId,
                          name: file.name,
                          url: getAvatarSrc(fileId),
                        });
                        onChange?.(fileId);
                      }
                      return false;
                    }
                    return true;
                  }}
                >
                  <Upload
                    accept="image/png, image/jpeg, image/svg+xml"
                    fileList={avatarList}
                    maxCount={1}
                    listType="picture-card"
                    onDrop={() => {
                      setAvatar(undefined);
                      onChange?.();
                    }}
                    style={{}}
                    onPreview={handlePreview}
                    onRemove={() => {
                      setAvatar(undefined);
                      onChange?.();
                    }}
                    customRequest={async (options) => {
                      try {
                        let { filename } = options;
                        const { file } = options;
                        if (!isString(file) && file) {
                          filename = (file as RcFile).name ?? filename;
                        }
                        filename = filename ?? new Date().getTime().toString();
                        filename = filename.substring(0, filename.lastIndexOf('.')) + '.png';
                        const fileId = await handleUploadFile(filename, file, onError);
                        if (fileId) {
                          setAvatar({
                            uid: fileId,
                            name: filename,
                            url: getAvatarSrc(fileId),
                          });
                          onChange?.(fileId);
                        }
                      } catch (error) {
                        onError?.(error);
                      }
                    }}
                    {...props}
                  >
                    {text ? null : (
                      <Popover
                        onOpenChange={(v) => {
                          if (!v || (!text && optionsRequest)) {
                            setPopupVisible(v);
                            if (v) {
                              loadMoreIconList();
                            } else {
                              setIconListHasMore(true);
                              setIconListPageNumber(0);
                            }
                          }
                        }}
                        open={popupVisible}
                        content={
                          <div style={{ width: 360, height: 200 }}>
                            <div
                              id="iconsScrollableDiv"
                              style={{
                                height: '100%',
                                overflow: 'auto',
                              }}
                            >
                              <InfiniteScroll
                                dataLength={iconList.length}
                                next={loadMoreIconList}
                                hasMore={iconListHasMore}
                                loader={<Skeleton avatar paragraph={{ rows: 1 }} active />}
                                endMessage={<Divider plain>End</Divider>}
                                scrollableTarget="iconsScrollableDiv"
                              >
                                <List<{ id: string }>
                                  grid={{ gutter: 16, column: 8 }}
                                  dataSource={iconList}
                                  style={{ margin: '0 8px' }}
                                  loading={loadingIconList}
                                  renderItem={({ id }) => {
                                    return (
                                      <ClassNames>
                                        {({ css }) => (
                                          <div
                                            className={css`
                                              :hover {
                                                background: rgba(0, 0, 0, 0.12);
                                              }
                                              padding: 5px;
                                            `}
                                            onClick={(e) => {
                                              e.stopPropagation();
                                              setAvatar({
                                                uid: id,
                                                name: id,
                                                url: getAvatarSrc(id),
                                              });
                                              onChange?.(id);
                                              setPopupVisible(false);
                                            }}
                                          >
                                            <Avatar src={id} />
                                          </div>
                                        )}
                                      </ClassNames>
                                    );
                                  }}
                                />
                              </InfiniteScroll>
                            </div>
                          </div>
                        }
                        placement="bottom"
                        trigger="hover"
                      >
                        <UploadOutlined
                          style={{ width: 112, height: 112, placeContent: 'center' }}
                        />
                      </Popover>
                    )}
                  </Upload>
                </ImgCrop>
              </div>
            </>
          ) : (
            <Avatar src={avatar?.url ?? text} />
          );
        }}
        {...proFieldProps}
      />
      <Modal open={previewOpen} footer={null} onCancel={handleCancel}>
        <img alt="example" style={{ width: '100%' }} src={previewImage} />
      </Modal>
    </>
  );
};
export const AvatarUploadField = createField<ProAvatarUploadProps>(ProAvatarUpload);
export default Avatar;
