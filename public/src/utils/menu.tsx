import Avatar from '@/components/Avatar';
import { getPages } from '@/services/idas/pages';
import type { MenuDataItem } from '@ant-design/pro-components';

export class MenuRender {
  private dynamicMenu: MenuDataItem[] | null;
  private constructor() {
    this.dynamicMenu = null;
  }
  private static instance: MenuRender | null = null;
  static getInstance(): MenuRender {
    if (MenuRender.instance === null) {
      MenuRender.instance = new MenuRender();
    }
    return MenuRender.instance;
  }
  public render(call: (menu: MenuDataItem[]) => void, force: boolean = false) {
    if (this.dynamicMenu !== null && !force) {
      return;
    }
    this.dynamicMenu = [];
    getPages({ pageSize: 1000, status: 'enabled' })
      .then(({ data }) => {
        if (data) {
          this.dynamicMenu = data.map((item) => ({
            key: `/page/${item.id}`,
            path: `/page/${item.id}`,
            locale: false,
            name: item.name,
            icon: item.icon ? <Avatar className="anticon" size={14} src={item.icon} /> : undefined,
            exact: true,
          }));
          call(this.dynamicMenu);
        }
      })
      .catch((err) => {
        console.error(`failed to load page menu: ${err}`);
      });
  }
}
