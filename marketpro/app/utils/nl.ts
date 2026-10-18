export interface NlShop {
  id: number;
  name: string;
  logo?: string;
  features?: string;
}

export interface NlCategory {
  id: number;
  name: string;
  icon?: string;
  image?: string;
  product_count?: number;
  is_active?: boolean;
}

export interface NlProduct {
  id: number;
  name: string;
  description?: string;
  price: number;
  orig_price?: number;
  image?: string;
  images?: string;
  stock_count: number;
  sales_count?: number;
  is_active?: boolean;
  delivery_type?: number;
  shop_id?: number;
  shop?: NlShop;
  category_id?: number;
  category?: NlCategory;
  promo_price?: number;
  promo_start?: string | null;
  promo_end?: string | null;
  is_promo_active?: boolean;
  effective_price?: number;
  avg_rating?: number;
  review_count?: number;
  created_at?: string;
}

export const prefersReducedMotion = () =>
  import.meta.client && !!window.matchMedia?.("(prefers-reduced-motion: reduce)").matches;

export function money(n: number | undefined | null) {
  return Number(n || 0).toFixed(2);
}

export function priceOf(p: NlProduct) {
  const now = Number(p.effective_price ?? p.price ?? 0);
  let was: number | undefined;
  if (p.is_promo_active && p.price > now) was = p.price;
  else if (p.orig_price && p.orig_price > now) was = p.orig_price;
  return { now, was, promo: !!p.is_promo_active };
}

export function stockOf(p: Pick<NlProduct, "sales_count" | "stock_count">) {
  const sold = p.sales_count || 0;
  const left = Math.max(0, p.stock_count || 0);
  const total = sold + left;
  const pct = total ? Math.min(100, Math.round((sold / total) * 100)) : 0;
  return { sold, left, pct, out: left <= 0, low: left > 0 && left <= 5 };
}

const NAME_ICONS: [RegExp, string][] = [
  [/vps|主机|服务器|云|cloud|server/i, "ph-hard-drives"],
  [/ai|gpt|claude|智能/i, "ph-sparkle"],
  [/游戏|game|steam|点卡/i, "ph-game-controller"],
  [/域名|domain/i, "ph-globe-hemisphere-west"],
  [/礼品|gift|充值/i, "ph-gift"],
  [/软件|授权|license|激活|code/i, "ph-code"],
  [/会员|账号|account/i, "ph-user-circle"],
  [/影视|视频|film|video/i, "ph-film-strip"],
  [/key|卡密|兑换/i, "ph-key"],
];

/** Category icon field may be an image URL, a Phosphor class, or a PascalCase Phosphor name. */
export function categoryIcon(cat?: Partial<NlCategory> | null): { img?: string; cls: string } {
  const fallback = "ph ph-package";
  if (!cat) return { cls: fallback };
  const icon = (cat.icon || "").trim();
  if (cat.image) return { img: cat.image, cls: fallback };
  if (/^(https?:)?\/\//.test(icon) || icon.startsWith("/")) return { img: icon, cls: fallback };
  if (icon.includes("ph-")) return { cls: icon.includes(" ") ? icon : `ph ${icon}` };
  if (/^[A-Z][A-Za-z]+$/.test(icon)) {
    const kebab = icon.replace(/([a-z0-9])([A-Z])/g, "$1-$2").toLowerCase();
    return { cls: `ph ph-${kebab}` };
  }
  const hit = NAME_ICONS.find(([re]) => re.test(cat.name || ""));
  return { cls: hit ? `ph ${hit[1]}` : fallback };
}

export function relTime(d?: string | null) {
  if (!d) return "";
  const s = Math.max(0, Math.floor((Date.now() - new Date(d).getTime()) / 1000));
  if (s < 60) return `${s} 秒前`;
  if (s < 3600) return `${Math.floor(s / 60)} 分钟前`;
  if (s < 86400) return `${Math.floor(s / 3600)} 小时前`;
  return `${Math.floor(s / 86400)} 天前`;
}

export function fmtDateTime(d?: string | null) {
  if (!d) return "";
  const t = new Date(d);
  const p = (n: number) => String(n).padStart(2, "0");
  return `${t.getFullYear()}-${p(t.getMonth() + 1)}-${p(t.getDate())} ${p(t.getHours())}:${p(t.getMinutes())}`;
}

export const ORDER_STATUS: Record<number, { label: string; tone: string; icon: string }> = {
  0: { label: "待支付", tone: "warning", icon: "ph ph-clock" },
  1: { label: "待发货", tone: "brand", icon: "ph ph-package" },
  4: { label: "已发货", tone: "success", icon: "ph ph-check" },
  2: { label: "已完成", tone: "success", icon: "ph ph-check-circle" },
  3: { label: "已取消", tone: "", icon: "ph ph-x" },
};

export function isNew(p: NlProduct) {
  if (!p.created_at) return false;
  return Date.now() - new Date(p.created_at).getTime() < 7 * 86400 * 1000;
}

export function productFlags(p: NlProduct) {
  const flags: { label: string; tone: string; icon?: string }[] = [];
  if (p.is_promo_active) flags.push({ label: "限时特惠", tone: "danger" });
  if ((p.delivery_type ?? 0) !== 1) flags.push({ label: "自动发货", tone: "volt", icon: "ph-fill ph-lightning" });
  if (flags.length < 2 && isNew(p)) flags.push({ label: "新品", tone: "brand" });
  if (flags.length < 2 && (p.sales_count || 0) >= 100) flags.push({ label: "热销", tone: "ink" });
  return flags.slice(0, 2);
}

export function productImages(p?: NlProduct | null) {
  if (!p) return [] as string[];
  let extras: string[] = [];
  try {
    extras = (JSON.parse(p.images || "[]") as string[]).filter(Boolean);
  } catch {
    extras = [];
  }
  if (p.image && !extras.includes(p.image)) extras = [p.image, ...extras];
  return extras;
}
