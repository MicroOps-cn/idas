import type { ButtonProps } from 'antd';
import { message } from 'antd';
import { Button } from 'antd';
import { isBoolean } from 'lodash';
import { useState } from 'react';

interface MultiStatusButtonProps extends ButtonProps {
  success?: React.ReactElement;
  failed?: React.ReactElement;
  onClick: (event: React.MouseEvent<HTMLElement, MouseEvent>) => Promise<void | boolean>;
}

export const MultiStatusButton: React.FC<MultiStatusButtonProps> = ({
  success,
  failed,
  onClick,
  ...props
}) => {
  const [loading, setLoading] = useState<boolean>(false);
  const [status, setStatus] = useState<0 | 1 | 2>(0);
  return (
    <Button
      {...props}
      loading={loading}
      onClick={async (e) => {
        if (onClick) {
          try {
            setLoading(true);
            const ok = await onClick(e);
            if (isBoolean(ok) && !ok) {
              setStatus(2);
            } else {
              setStatus(1);
            }
          } catch (error) {
            message.error(`${error}`, 3);
            setStatus(2);
          } finally {
            setLoading(false);
          }
        }
      }}
    >
      {status === 1 && success ? success : status === 2 && failed ? failed : props.children}
    </Button>
  );
};
export default MultiStatusButton;
