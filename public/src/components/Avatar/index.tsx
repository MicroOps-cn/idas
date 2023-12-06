import type { AvatarProps as AntdAvatarProps } from 'antd';
import { Modal } from 'antd';
import { Upload } from 'antd';
import { Avatar as AntdAvatar } from 'antd';
import ImgCrop from 'antd-img-crop';
import type { RcFile, UploadFile } from 'antd/lib/upload';
import { isString } from 'lodash';
import { useEffect, useState } from 'react';

import { uploadFile as postFile } from '@/services/idas/files';
import { UploadOutlined } from '@ant-design/icons';
import { ProField } from '@ant-design/pro-components';
import type { ProFormFieldProps as ProFormFieldItemProps } from '@ant-design/pro-components';
import { createField } from '@ant-design/pro-form/es/BaseForm/createField';

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

interface ProAvatarUploadProps {
  onError?: (error: any) => void;
}

export const ProAvatarUpload: React.FC<ProFormFieldItemProps<ProAvatarUploadProps>> = ({
  fieldProps: { onError, ...fieldProps } = { onError: undefined },
  proFieldProps,
}: ProFormFieldItemProps<ProAvatarUploadProps>) => {
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

  return (
    <>
      <ProField
        mode="edit"
        fieldProps={fieldProps}
        renderFormItem={(text, { onChange, mode, value, ...props }) => {
          const avatarList = text ? [{ uid: value, name: 'img', url: getAvatarSrc(text) }] : [];
          return mode === 'edit' ? (
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
                {text ? null : <UploadOutlined />}
              </Upload>
            </ImgCrop>
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
export const AvatarUploadField =
  createField<ProFormFieldItemProps<ProAvatarUploadProps>>(ProAvatarUpload);
export default Avatar;
