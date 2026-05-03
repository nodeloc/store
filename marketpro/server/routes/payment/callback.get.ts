/**
 * 支付回调透传路由
 * Nitro proxy 会跟随后端的 302 redirect，导致浏览器收不到重定向。
 * 这里手动转发请求，并将后端的 Location header 原样返回给浏览器。
 */
export default defineEventHandler(async (event) => {
  const backendUrl = process.env.BACKEND_URL || "http://backend:8080";
  const query = getQuery(event);

  // 拼接完整的后端回调 URL（带所有查询参数）
  const params = new URLSearchParams(query as Record<string, string>).toString();
  const targetUrl = `${backendUrl}/payment/callback?${params}`;

  // 请求后端，不跟随重定向
  const res = await fetch(targetUrl, { redirect: "manual" });

  // 后端返回 302 → 把 Location 透传给浏览器
  if (res.status === 301 || res.status === 302 || res.status === 303 || res.status === 307 || res.status === 308) {
    const location = res.headers.get("location");
    if (location) {
      return sendRedirect(event, location, res.status);
    }
  }

  // 其他情况（不太可能）：透传响应体
  const body = await res.text();
  setResponseStatus(event, res.status);
  return body;
});
