/**
 * @see https://umijs.org/zh-CN/plugins/plugin-access
 * */

export default function access(initialState: { currentUser?: API.UserInfo | undefined }) {
  const { currentUser } = initialState || {};
  return {
    canAnonymous: true,
    canUser: currentUser?.id,
    forbidden: false,
    canAdmin: currentUser && currentUser.role === 'admin',
    canEditor: currentUser && (currentUser.role === 'editor' || currentUser.role === 'admin'),
    canViewer:
      currentUser &&
      (currentUser.role === 'viewer' ||
        currentUser.role === 'editor' ||
        currentUser.role === 'admin'),
  };
}
